
package bp

import (
	"code.google.com/p/go.exp/fsnotify"
)

// Per-instance session object
type Session struct {
	// Absolute pathname to root
	root string
	// Configuration structure
	config
	// Absolute pathname mapping to fsnotify watcher objects
	watchers map[string]fsnotify.Watcher
}

func NewSession(pathname string) *Session {
	s := new(Session)
	s.config.update()
	return s
}

