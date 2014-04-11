package bp

import (
	"os"
	"os/exec"
	"testing"
)

var root string
var file string
var s *Session

func init() {
	root = "/tmp"
	file = "testfile"
}

func TestNewSession(t *testing.T) {
	t.Logf("Using dir '%s'.", root)
	s = NewSession(root)
}

func TestCreateDatabase(t *testing.T) {
	os.Remove(root + "/.bp/object.db")
	err := s.createDatabase()
	if err != nil {
		t.Errorf("Database could not be created:\n%s", err)
	}
	path := s.absolutePath(".bp/object.db")
	cmd := exec.Command("/usr/bin/sqlite3", path, ".dump")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Error opening sqlite3: %s", err)
		t.Logf("Schema output:\n%s", out)
	}
}

func TestVerifyObject(t *testing.T) {
	cmd := exec.Command("/bin/dd", "if=/dev/urandom", "of=" + root + "/" + file,
				"count=1", "bs=10K")
	err := cmd.Run()
	if err != nil {
		t.Errorf("Could not create temporary test file:\n%s", err)
	}
	obj, err := s.getFile(file)
	if err != nil {
		t.Errorf("Could not stat file '%s':\n%s", file, err)
	}
	err = s.addObject(file, 0, obj)
	if err != nil {
		t.Errorf("Failed to add object:\n%s", err)
	} else {
		o, err := s.getObject(file)
		if err != nil {
			t.Errorf("Failed to retrieve object:\n%s", err)
		} else {
			if !o.equal(obj) {
				t.Errorf("Object inconsistent across add/retrieve.\n" +
						"Values added:\n\t%v\nValues Retrieved:\n\t%v", obj, o)
			}
		}
	}
}

func TestNewWatch(t *testing.T) {
	err := s.addWatch(file)
	if err != nil {
		t.Errorf("Watcher could not be set for '%s'.\n%s", file, err)
	}
}
