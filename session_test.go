package bp

import (
	"os"
	"os/exec"
	"testing"
)

var root string

func init() {
	if len(os.Args) > 1 {
		root = os.Args[1]
	} else {
		root = "/tmp"
	}
	return
}

func TestNewSession(t *testing.T) {
	var path string
	//	path, _ := os.Getwd()
	if len(os.Args) == 3 {
		path = os.Args[2]
	}
	t.Logf("Using dir '%s'", root)
	s := NewSession(root)
	err := s.addWatch(path)
	if err != nil {
		t.Errorf("Watcher could not be set for '%s'.\n%s", path, err)
	}
}

func TestCreateDatabase(t *testing.T) {
	if len(os.Args) > 2 && root == os.Args[1] {
		t.Skip("Skipping new database test to avoid overwriting real data.")
	} else {
		os.Remove("/tmp/.bp/object.db")
	}
	s := NewSession(root)
	err := s.createDatabase()
	if err != nil {
		t.Errorf("Database could not be created:\n%s", err)
	}
	path := s.absolutePath(".bp/object.db")
	cmd := exec.Command("/usr/bin/sqlite3", path, ".dump")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Error opening sqlite3: %s", err)
	} else {
		t.Logf("Schema output:\n%s", out)
	}
}
