package bp

import (
	"os"
	"testing"
)

func TestNewSession(t *testing.T) {
	var root string
	var path string
//	path, _ := os.Getwd()
	if len(os.Args) == 3 {
		root = os.Args[1]
		path = os.Args[2]
	} else {
		root = "/home/victor/Dropbox/workspace/"
	}
	t.Logf("Using dir '%s'", root)
	s := NewSession(root)
	err := s.addWatch(path)
	if err != nil {
		t.Errorf("Watcher could not be set for '%s'.\n%s", path, err)
	}
}
