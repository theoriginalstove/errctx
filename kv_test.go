package errctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	kv1 := KV{"key1": "value1"}
	kv2 := KV{"key2": 2}

	output := Merge(kv1, kv2)
	assert.Equal(t, 2, len(output))
	assert.Equal(t, KV{"key1": "value1", "key2": 2}, output)

	kv3 := KV{"key3": false}
	output = Merge(output, kv3)
	assert.Equal(t, 3, len(output))
	assert.Equal(t, KV{"key1": "value1", "key2": 2, "key3": false}, output)
}
