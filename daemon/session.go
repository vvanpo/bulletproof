package main

import (
	"github.com/vvanpo/bulletproof/object"
	"gopkg.in/fsnotify.v0"
	"fmt"
	"os"
	"path/filepath"
	"log"
)

// Per-instance session object
type session struct {
	// Absolute pathname to root
	root string
	// The object store instance
	store object.ObjectStore
	// File watcher instance
	watcher *fsnotify.Watcher
}

func newSession(root string) *session {
	s := new(session)
	s.root = root
	os.Chdir(root)
	var err error
	s.store, err = object.CreateSqlite()
	if err != nil {
		log.Fatalf("Failed to initialize object store:\n%s", err)
	}
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to initialize inotify instance:\n%s", err)
	}
	return s
}

func (s *session) start() error {
	// p = map[<path>]<consistent?>
	p, err := s.store.VerifyAllObjects()
	if err != nil { return err }
	for k := range(p) {
		if !p[i] {
			err = s.inconsistent(k)
			if err != nil { return err }
		}
		err = s.addWatch(k)
		if err != nil { return err }
	}
	return nil
}

func (s *session) addWatch(path string) error {
	err := s.watcher.Watch(path)
	if err != nil {
		return fmt.Errorf("fsnotify: %s", err)
	}
	return nil
}

