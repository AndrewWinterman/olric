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

package dmap

import "github.com/buraksezer/olric/internal/protocol/resp"

func (s *Service) RegisterHandlers() {
	s.respServer.ServeMux().HandleFunc(resp.PutCmd, s.putCommandHandler)
	s.respServer.ServeMux().HandleFunc(resp.GetCmd, s.getCommandHandler)
	s.respServer.ServeMux().HandleFunc(resp.DelCmd, s.delCommandHandler)
	s.respServer.ServeMux().HandleFunc(resp.DelEntryCmd, s.delEntryCommandHandler)
	s.respServer.ServeMux().HandleFunc(resp.GetEntryCmd, s.getEntryCommandHandler)
	s.respServer.ServeMux().HandleFunc(resp.PutReplicaCmd, s.putReplicaCommandHandler)
	s.respServer.ServeMux().HandleFunc(resp.ExpireCmd, s.expireCommandHandler)
	s.respServer.ServeMux().HandleFunc(resp.DestroyCmd, s.destroyCommandHandler)
}
