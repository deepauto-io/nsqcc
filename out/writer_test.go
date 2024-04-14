/*
Copyright 2024 The nsqcc Authors.

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

package out

import (
	"context"
	"github.com/deepauto-io/nsqcc/filepath/ifs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNSQWriter(t *testing.T) {
	cfg := NewConfig()
	write, err := NewNSQWriter(cfg, ifs.OS())
	assert.NoError(t, err)

	err = write.Connect(context.Background())
	assert.NoError(t, err)

	err = write.WriteWithContext(context.Background(), "hello", []byte("world"))
	assert.NoError(t, err)
}
