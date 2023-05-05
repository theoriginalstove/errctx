package errctx

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	. "testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type key int

func TestErrCtx(t *T) {
	err := errors.New("foo")

	assert.Equal(t, err, Base(err))

	err1 := Set(err, key(0), "a")
	assert.Equal(t, err.Error(), err1.Error())
	assert.Equal(t, err, Base(err1))
	assert.True(t, errors.Is(err1, err))
	assert.Nil(t, Get(err, key(0)))
	assert.Equal(t, "a", Get(err1, key(0)))

	err2 := Set(err, key(1), "b")
	assert.NotEqual(t, err1, err2)
	assert.Equal(t, err.Error(), err2.Error())
	assert.Equal(t, err, Base(err2))
	assert.True(t, errors.Is(err2, err))
	assert.Nil(t, Get(err, key(1)))
	assert.Nil(t, Get(err2, key(0)))
	assert.Equal(t, "b", Get(err2, key(1)))

	err3 := Set(err2, key(2), "c")
	assert.Equal(t, err.Error(), err3.Error())
	assert.Equal(t, err, Base(err3))
	assert.True(t, errors.Is(err3, err))
	assert.Nil(t, Get(err3, key(0)))
	assert.Nil(t, Get(err2, key(2)))
	assert.Equal(t, "b", Get(err3, key(1)))
	assert.Equal(t, "c", Get(err3, key(2)))

	assert.True(t, err3.(errctx).Is(err3))
	assert.True(t, err3.(errctx).Is(err2))
	assert.True(t, err3.(errctx).Is(err2.(errctx).err))
}

func TestMark(t *T) {
	err := errors.New("bar")

	l, ok := Line(err)
	assert.False(t, ok)
	assert.Empty(t, l)

	_, _, ln, ok := runtime.Caller(0)
	require.True(t, ok)
	err = Mark(err)
	l, ok = Line(err)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("errctx_test.go:%d", ln+2), l)

	// calling it again shouldn't do anything
	err = Mark(err)
	l, ok = Line(err)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("errctx_test.go:%d", ln+2), l)

	err = func() error {
		// 1 should return the anonymous function
		return MarkSkip(errors.New("bar"), 1)
	}()
	_, _, ln, ok = runtime.Caller(0)
	require.True(t, ok)
	l, ok = Line(err)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf("errctx_test.go:%d", ln-1), l)
}

func TestErrKV(t *T) {
	err := errors.New("foo")
	assert.Equal(t, KV{"err": err.Error()}, ErrKV(err))

	kv := KV{"a": "a"}
	err2 := ErrWithKV(err, kv)
	assert.Equal(t, KV{"err": err.Error()}, ErrKV(err))
	assert.Equal(t, KV{"err": err2.Error(), "a": "a", "source": "errctx_test.go:87"}, ErrKV(err2))

	// changing the kv now shouldn't do anything
	kv["a"] = "b"
	assert.Equal(t, KV{"err": err.Error()}, ErrKV(err))
	assert.Equal(t, KV{"err": err2.Error(), "a": "a", "source": "errctx_test.go:87"}, ErrKV(err2))

	// a new ErrWithKV shouldn't affect the previous one
	err3 := ErrWithKV(err2, KV{"b": "b"})
	assert.Equal(t, KV{"err": err.Error()}, ErrKV(err))
	assert.Equal(t, KV{"err": err2.Error(), "a": "a", "source": "errctx_test.go:87"}, ErrKV(err2))
	assert.Equal(t, KV{"err": err3.Error(), "a": "a", "b": "b", "source": "errctx_test.go:87"}, ErrKV(err3))

	// make sure precedence works
	err4 := ErrWithKV(err3, KV{"b": "bb"})
	assert.Equal(t, KV{"err": err.Error()}, ErrKV(err))
	assert.Equal(t, KV{"err": err2.Error(), "a": "a", "source": "errctx_test.go:87"}, ErrKV(err2))
	assert.Equal(t, KV{"err": err3.Error(), "a": "a", "b": "b", "source": "errctx_test.go:87"}, ErrKV(err3))
	assert.Equal(t, KV{"err": err4.Error(), "a": "a", "b": "bb", "source": "errctx_test.go:87"}, ErrKV(err4))

	err = nil
	assert.Equal(t, KV{}, ErrKV(err))
}

// copied from go-llog errctx_test.go

func TestCtxKV(t *T) {
	ctx := context.Background()
	assert.Equal(t, KV{}, CtxKV(ctx))

	kv := KV{"a": "a"}
	ctx2 := CtxWithKV(ctx, kv)
	assert.Equal(t, KV{}, CtxKV(ctx))
	assert.Equal(t, KV{"a": "a"}, CtxKV(ctx2))

	// changing the kv now shouldn't do anything
	kv["a"] = "b"
	assert.Equal(t, KV{}, CtxKV(ctx))
	assert.Equal(t, KV{"a": "a"}, CtxKV(ctx2))

	// a new CtxWithKV shouldn't affect the previous one
	ctx3 := CtxWithKV(ctx2, KV{"b": "b"})
	assert.Equal(t, KV{}, CtxKV(ctx))
	assert.Equal(t, KV{"a": "a"}, CtxKV(ctx2))
	assert.Equal(t, KV{"a": "a", "b": "b"}, CtxKV(ctx3))

	// make sure precedence works
	ctx4 := CtxWithKV(ctx3, KV{"b": "bb"})
	assert.Equal(t, KV{}, CtxKV(ctx))
	assert.Equal(t, KV{"a": "a"}, CtxKV(ctx2))
	assert.Equal(t, KV{"a": "a", "b": "b"}, CtxKV(ctx3))
	assert.Equal(t, KV{"a": "a", "b": "bb"}, CtxKV(ctx4))
}
