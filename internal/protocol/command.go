// Copyright 2018-2021 Burak Sezer
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

package protocol

const (
	GetCmd                 = "dm.get"
	GetEntryCmd            = "dm.getentry"
	PutCmd                 = "dm.put"
	PutEntryCmd            = "dm.putentry"
	DelCmd                 = "dm.del"
	DelEntryCmd            = "dm.delentry"
	ExpireCmd              = "dm.expire"
	PExpireCmd             = "dm.pexpire"
	DestroyCmd             = "dm.destroy"
	QueryCmd               = "dm.query"
	IncrCmd                = "dm.incr"
	DecrCmd                = "dm.decr"
	GetPutCmd              = "dm.getput"
	LockCmd                = "dm.lock"
	UnlockCmd              = "dm.unlock"
	LockLeaseCmd           = "dm.locklease"
	PLockLeaseCmd          = "dm.plocklease"
	ScanCmd                = "dm.scan"
	PingCmd                = "ping"
	MoveFragmentCmd        = "olric.internal.movefragment"
	UpdateRoutingCmd       = "olric.internal.updaterouting"
	LengthOfPartCmd        = "olric.internal.lengthofpart"
	ClusterRoutingTableCmd = "cluster.routingtable"
)

const StatusOK = "OK"

type PubSubCommands struct {
	Publish        string
	Subscribe      string
	PSubscribe     string
	PubSubChannels string
	PubSubNumpat   string
	PubSubNumsub   string
}

var PubSub = &PubSubCommands{
	Publish:        "publish",
	Subscribe:      "subscribe",
	PSubscribe:     "psubscribe",
	PubSubChannels: "pubsub channels",
	PubSubNumpat:   "pubsub numpat",
	PubSubNumsub:   "pubsub numsub",
}
