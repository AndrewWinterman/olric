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

package routing_table

import (
	"errors"
	"time"
)

var (
	ErrServerGone  = errors.New("server is gone")
	ErrClusterJoin = errors.New("cannot join the cluster")
)

// bootstrapCoordinator prepares the very first routing table and bootstraps the coordinator node.
func (r *RoutingTable) bootstrapCoordinator() error {
	r.Lock()
	defer r.Unlock()

	r.fillRoutingTable()
	_, err := r.updateRoutingTableOnCluster()
	if err != nil {
		return err
	}
	// The coordinator bootstraps itself.
	r.markBootstrapped()
	r.log.V(2).Printf("[INFO] The cluster coordinator has been bootstrapped")
	return nil
}

func (r *RoutingTable) attemptToJoin() error {
	attempts := 0
	for attempts < r.config.MaxJoinAttempts {
		select {
		case <-r.ctx.Done():
			// The node is gone.
			return ErrServerGone
		default:
		}

		attempts++
		n, err := r.discovery.Join()
		if err == nil {
			r.log.V(2).Printf("[INFO] Join completed. Synced with %d initial nodes", n)
			return nil
		}

		r.log.V(2).Printf("[ERROR] Join attempt returned error: %s", err)
		if r.IsBootstrapped() {
			r.log.V(2).Printf("[INFO] Bootstrapped by the cluster coordinator")
			return nil
		}

		r.log.V(2).Printf("[INFO] Awaits for %s to join again (%d/%d)",
			r.config.JoinRetryInterval, attempts, r.config.MaxJoinAttempts)
		<-time.After(r.config.JoinRetryInterval)
	}
	return ErrClusterJoin
}

func (r *RoutingTable) tryWithInterval(max int, interval time.Duration, f func() error) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var err error
	err = f()
	if err == nil {
		// Done. No need to try with interval
		return nil
	}

	var count = 1
loop:
	for count < max {
		select {
		case <-ticker.C:
			count++
			err = f()
			if err == nil {
				break loop
			}
		case <-r.ctx.Done():
			// the server is gone
			return ErrServerGone
		}
	}
	return err
}
