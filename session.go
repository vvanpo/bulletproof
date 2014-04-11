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
	// Map of all watched paths to events
	event map[string]*fsnotify.FileEvent
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

func (s *Session) start() error {
	c, err := dbConn()
	if err != nil { return err }
	q, err := c.Query("SELECT path FROM folder")
	folders := make([]string, 10)
	for ; err == nil; err = q.Next() {

	}
	// err = filepath.Walk(s.root, walk)
}

func (s *Session) addWatch(path string) error {
	err := s.watcher.Watch(s.absolutePath(path))
	if err != nil {
		return fmt.Errorf("fsnotify: %s", err)
	}
	return nil
}
