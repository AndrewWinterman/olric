// Copyright 2018-2020 Burak Sezer
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package olric

import (
	"fmt"
	"github.com/buraksezer/olric/internal/cluster/partitions"
	"github.com/buraksezer/olric/internal/discovery"
	"github.com/buraksezer/olric/internal/protocol"
	"github.com/buraksezer/olric/pkg/storage"
	"github.com/vmihailenco/msgpack"
)

// TODO: Rename this
type dmapbox struct {
	PartID    uint64
	Kind      partitions.Kind
	Name      string
	Payload   []byte
	AccessLog map[uint64]int64
}

func (dm *dmap) Move(partID uint64, kind partitions.Kind, name string, owner discovery.Member) error {
	dm.Lock()
	defer dm.Unlock()

	payload, err := dm.storage.Export()
	if err != nil {
		return err
	}
	data := &dmapbox{
		PartID:  partID,
		Kind:    kind,
		Name:    name,
		Payload: payload,
	}
	// config structure will be regenerated by mergeDMap. Just pack the accessLog.
	if dm.config != nil && dm.config.accessLog != nil {
		data.AccessLog = dm.config.accessLog
	}
	value, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}

	req := protocol.NewSystemMessage(protocol.OpMoveDMap)
	req.SetValue(value)
	// TODO: Check errors etc
	_, err = dm.client.RequestTo2(owner.String(), req)
	return err
}

func (db *Olric) selectVersionForMerge(dm *dmap, hkey uint64, entry storage.Entry) (storage.Entry, error) {
	current, err := dm.storage.Get(hkey)
	if err == storage.ErrKeyNotFound {
		return entry, nil
	}
	if err != nil {
		return nil, err
	}
	versions := []*version{{entry: current}, {entry: entry}}
	versions = db.sortVersions(versions)
	return versions[0].entry, nil
}

func (db *Olric) mergeDMaps(part *partitions.Partition, data *dmapbox) error {
	dm, err := db.getOrCreateDMap(part, data.Name)
	if err != nil {
		return err
	}

	// Acquire dmap's lock. No one should work on it.
	dm.Lock()
	defer dm.Unlock()
	defer part.Map().Store(data.Name, dm)

	engine, err := dm.storage.Import(data.Payload)
	if err != nil {
		return err
	}

	// Merge accessLog.
	if dm.config != nil && dm.config.accessLog != nil {
		dm.config.Lock()
		for hkey, t := range data.AccessLog {
			if _, ok := dm.config.accessLog[hkey]; !ok {
				dm.config.accessLog[hkey] = t
			}
		}
		dm.config.Unlock()
	}

	if dm.storage.Stats().Length == 0 {
		// DMap has no keys. Set the imported storage instance.
		// The old one will be garbage collected.
		dm.storage = engine
		return nil
	}

	// DMap has some keys. Merge with the new one.
	var mergeErr error
	engine.Range(func(hkey uint64, entry storage.Entry) bool {
		winner, err := db.selectVersionForMerge(dm, hkey, entry)
		if err != nil {
			mergeErr = err
			return false
		}
		// TODO: Don't put the winner again if it comes from dm.storage
		mergeErr = dm.storage.Put(hkey, winner)
		if mergeErr == storage.ErrFragmented {
			db.wg.Add(1)
			go db.callCompactionOnStorage(dm)
			mergeErr = nil
		}
		if mergeErr != nil {
			return false
		}
		return true
	})
	return mergeErr
}

func (db *Olric) checkOwnership(part *partitions.Partition) bool {
	owners := part.Owners()
	for _, owner := range owners {
		if owner.CompareByID(db.rt.This()) {
			return true
		}
	}
	return false
}

func (db *Olric) moveDMapOperation(w, r protocol.EncodeDecoder) {
	err := db.isOperable()
	if err != nil {
		db.errorResponse(w, err)
		return
	}

	req := r.(*protocol.SystemMessage)
	box := &dmapbox{}
	err = msgpack.Unmarshal(req.Value(), box)
	if err != nil {
		db.log.V(2).Printf("[ERROR] Failed to unmarshal dmap: %v", err)
		db.errorResponse(w, err)
		return
	}

	var part *partitions.Partition
	if box.Kind == partitions.PRIMARY {
		part = db.primary.PartitionById(box.PartID)
	} else {
		part = db.backup.PartitionById(box.PartID)
	}

	// Check ownership before merging. This is useful to prevent data corruption in network partitioning case.
	if !db.checkOwnership(part) {
		db.log.V(2).Printf("[ERROR] Received DMap: %s on PartID: %d (kind: %s) doesn't belong to this node (%s)",
			box.Name, box.PartID, box.Kind, db.rt.This())
		err := fmt.Errorf("partID: %d (kind: %s) doesn't belong to %s: %w", box.PartID, box.Kind, db.rt.This(), ErrInvalidArgument)
		db.errorResponse(w, err)
		return
	}

	db.log.V(2).Printf("[INFO] Received DMap (kind: %s): %s on PartID: %d", box.Kind, box.Name, box.PartID)

	err = db.mergeDMaps(part, box)
	if err != nil {
		db.log.V(2).Printf("[ERROR] Failed to merge DMap: %v", err)
		db.errorResponse(w, err)
		return
	}
	w.SetStatus(protocol.StatusOK)
}