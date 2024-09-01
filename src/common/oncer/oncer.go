package oncer

import (
	"fmt"
	"runtime"
	"sync"
)

// 通用注册Once

var mu sync.Mutex
var caller = map[string]*sync.Once{}

func Do(fn func(), withline ...bool) {
	mu.Lock()
	_, file, line, _ := runtime.Caller(1)
	if len(withline) > 0 && withline[0] {
		file = fmt.Sprintf("%s:%d", file, line)
	}
	if _, ok := caller[file]; !ok {
		var once sync.Once
		caller[file] = &once
	}
	mu.Unlock()
	caller[file].Do(fn)
}
