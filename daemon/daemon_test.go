package main

import (
	"testing"
	"github.com/vvanpo/bulletproof/object"
)

var _ object.Object		// FOR DEBUG

var root string
var file string
var s *session

func init() {
	root = "/tmp"
	file = "testfile"
}

func TestNewSession(t *testing.T) {
	t.Logf("Using dir '%s'.", root)
	s = newSession(root)
}

func TestAddWatch(t *testing.T) {
	err := s.addWatch(file)
	if err != nil {
		t.Errorf("Watcher could not be set for '%s'.\n%s", file, err)
	}
}
