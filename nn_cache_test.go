package nncache

import (
	"github.com/bmizerany/assert"
	"testing"
	"time"
)

func TestNewNNCache1(t *testing.T) {
	cache, _ := NewNNCache(DefaultConfig(10 * time.Minute))
	err := cache.Set("luwenchao", []byte("zhendishuai"))
	assert.Equal(t, err, nil)
	res, err := cache.Get("luwenchao")
	assert.Equal(t, err, nil)
	assert.Equal(t, res, []byte("zhendishuai"))
}
func TestNewNNCache2(t *testing.T) {
	cache, _ := NewNNCache(DefaultConfig(10 * time.Minute))
	err := cache.Set("luwenchao", []byte("zhendishuai"))
	assert.Equal(t, err, nil)
	res, err := cache.Get("someting")
	assert.Equal(t, string(res), "")
	assert.Equal(t, err.Error(), "Entry not found")
}
