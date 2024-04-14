/*
Copyright 2022 The deepauto-io LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ifs

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

type testFS struct {
	fstest.MapFS
}

func (t testFS) MkdirAll(path string, perm fs.FileMode) error {
	return errors.New("not implemented")
}

func (t testFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	return nil, errors.New("not implemented")
}

func (t testFS) Remove(name string) error {
	return errors.New("not implemented")
}

func TestOSAccess(t *testing.T) {
	var fs FS = testFS{}

	require.False(t, IsOS(fs))

	fs = OS()

	require.True(t, IsOS(fs))
}
