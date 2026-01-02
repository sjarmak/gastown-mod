package llm

import (
	"path/filepath"
	"strings"
	"sync"
)

var (
	registryMu sync.RWMutex
	registry           = map[string]Adapter{}
	fallback   Adapter = shellAdapter{}
)

// RegisterAdapter registers an adapter implementation by name.
func RegisterAdapter(name string, adapter Adapter) {
	key := canonicalName(name)
	registryMu.Lock()
	defer registryMu.Unlock()
	if adapter == nil {
		delete(registry, key)
		return
	}
	registry[key] = adapter
}

func getAdapter(options LaunchOptions) Adapter {
	key := canonicalName(options.Provider)
	if key == "" {
		key = canonicalName(options.Runtime.Command)
	}

	registryMu.RLock()
	adapter, ok := registry[key]
	registryMu.RUnlock()
	if ok {
		return adapter
	}
	return fallback
}

func canonicalName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}
	base := filepath.Base(trimmed)
	return strings.ToLower(base)
}
