package goview

import (
	"net/http"
	"sync"
	"sync/atomic"
)

var instance atomic.Pointer[ViewEngine]
var defaultInstanceOnce sync.Once

// Use setting default instance engine
func Use(engine *ViewEngine) {
	instance.Store(engine)
}

// Render render view template with default instance
func Render(w http.ResponseWriter, status int, name string, data interface{}) error {
	inst := instance.Load()
	if inst == nil {
		defaultInstanceOnce.Do(func() {
			instance.CompareAndSwap(nil, Default())
		})
		inst = instance.Load()
	}
	return inst.Render(w, status, name, data)
}
