// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginKeyIsValid(t *testing.T) {
	kv := PluginKeyValue{PluginId: "someid", Key: "somekey", Value: []byte("somevalue")}
	assert.Nil(t, kv.IsValid())

	kv.PluginId = ""
	assert.NotNil(t, kv.IsValid())

	kv.PluginId = "someid"
	kv.Key = ""
	assert.NotNil(t, kv.IsValid())

	kv.Key = "this is an extremely long key and should be invalid and this is being verified in this test"
	assert.NotNil(t, kv.IsValid())
}
