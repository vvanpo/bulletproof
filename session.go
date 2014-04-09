package bp

import (
	"code.google.com/p/go.exp/fsnotify"
	"fmt"
	"log"
	"path/filepath"
)

// Per-instance session object
type Session struct {
	// Absolute pathname to root
	root string
	// File watcher instance
	watcher *fsnotify.Watcher
	// List of watched paths
	watch []string
}

func NewSession(root string) *Session {
	s := new(Session)
	s.root = root
	var err error
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to initialize inotify instance:\n%s", err)
	}
	return s
}

// Convenience method to return absolute pathname
func (s *Session) absolutePath(path string) string {
	return filepath.Join(s.root, path)
}

func (s *Session) addWatch(path string) error {
	err := s.watcher.Watch(s.absolutePath(path))
	if err != nil {
		return fmt.Errorf("fsnotify: %s", err)
	}
	return nil
}
