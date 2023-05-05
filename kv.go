package errctx

import (
	"fmt"
	"sort"
	"strings"
)

// KV is used to provide context to a log entry in the form of a dynamic set of
// key/value pairs which can be different for every entry.
type KV map[string]interface{}

// Copy returns a copy of the KV being called on. This method will never return
// nil
func (kv KV) Copy() KV {
	nkv := make(KV, len(kv))
	for k, v := range kv {
		nkv[k] = v
	}
	return nkv
}

// Merge takes in multiple KVs and returns a single KV which is the union of all
// the passed in ones. Key/vals on the rightmost of the set take precedence over
// conflicting ones to the left. This function will never return nil
func Merge(kvs ...KV) KV {
	kv := make(KV, len(kvs))
	for i := range kvs {
		for k, v := range kvs[i] {
			kv[k] = v
		}
	}
	return kv
}

// Set returns a copy of the KV being called on with the given key/val set on
// it. The original KV is unaffected
func (kv KV) Set(k string, v interface{}) KV {
	nkv := kv.Copy()
	nkv[k] = v
	return nkv
}

// StringSlice converts the KV into a slice of [2]string entries (first index is
// the key, second is the string form of the value).
func (kv KV) StringSlice() [][2]string {
	slice := make([][2]string, 0, len(kv))
	for kstr, v := range kv {
		vstr := fmt.Sprint(v)
		// TODO this is only here because logstash is dumb and doesn't
		// properly handle escaped quotes. Once
		// https://github.com/elastic/logstash/issues/1645
		// gets figured out this Replace can be removed
		vstr = strings.ReplaceAll(vstr, `"`, `'`)
		slice = append(slice, [2]string{kstr, vstr})
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i][0] < slice[j][0]
	})
	return slice
}
