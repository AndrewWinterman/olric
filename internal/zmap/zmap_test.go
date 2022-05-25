// Copyright 2018-2022 Burak Sezer
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

package zmap

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/buraksezer/olric/internal/testzmap"
	"github.com/buraksezer/olric/internal/zmap/config"
	"github.com/stretchr/testify/require"
)

func testZMapConfig(t *testing.T) *config.Config {
	tmpdir, err := ioutil.TempDir("", "olric-zmap")
	require.NoError(t, err)
	c := config.DefaultConfig()
	c.DataDir = tmpdir

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(tmpdir))
	})
	return c
}

func TestZMap_Name(t *testing.T) {
	tz := testzmap.New(t, NewService)
	s := tz.AddStorageNode(nil).(*Service)

	zm, err := s.NewZMap("myzmap", nil)
	require.NoError(t, err)
	require.Equal(t, "myzmap", zm.Name())
}

func TestZMap_NewZMap(t *testing.T) {
	tz := testzmap.New(t, NewService)
	s := tz.AddStorageNode(nil).(*Service)

	c := testZMapConfig(t)

	zm1, err := s.NewZMap("myzmap", c)
	require.NoError(t, err)

	zm2, err := s.NewZMap("myzmap", c)
	require.NoError(t, err)

	require.Equal(t, zm1, zm2)
	require.Len(t, s.zmaps, 1)
}
