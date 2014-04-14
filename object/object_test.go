package object

import (
	"os"
	"os/exec"
	"testing"
)

var root string
var file string
var s ObjectStore

func init() {
	root = "/tmp"
	file = "testfile"
	os.Chdir(root)
}

func testFile() {
	t := new(testing.T)
	cmd := exec.Command("/bin/dd", "if=/dev/urandom", "of="+root+"/"+file,
		"count=1", "bs=10K")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Could not create temporary test file:\n%s", err)
	}
}

func TestCreateSqlite(t *testing.T) {
	var err error
	s, err = CreateSqlite()
	if err != nil {
		t.Fatalf("Object store could not be created:\n%s", err)
	}
	path := root + "/.bp/object.db"
	cmd := exec.Command("/usr/bin/sqlite3", path, ".dump")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Error opening sqlite3: %s", err)
		t.Logf("Schema output:\n%s", out)
	}
}

func TestAddViewRemove(t *testing.T) {
	testFile()
	obj, err := s.StatObject(file)
	if err != nil {
		t.Errorf("Could not stat file '%s':\n%s", file, err)
	}
	err = s.AddObject(file, 0, obj)
	if err != nil {
		t.Errorf("Failed to add object:\n%s", err)
	}
	o, err := s.ViewObject(file)
	if err != nil {
		t.Error(err)
	} else {
		if !o.Equal(obj) {
			t.Errorf("Object inconsistent across add/retrieve.\n"+
				"Values added:\n\t%v\nValues Retrieved:\n\t%v", obj, o)
		}
	}
}

func TestVerifyObject(t *testing.T) {
	obj, _ := s.StatObject(file)
	err := s.AddObject(file, 0, obj)
	if err != nil {
		t.Errorf("Failed to add object:\n%s", err)
	}
	v, err := s.VerifyObject(file)
	if !v || err != nil {
		t.Error()
	}
	testFile()
	v, err = s.VerifyObject(file)
	if v || err != nil {
		t.Error()
	}
}

func TestRemoveObject(t *testing.T) {
	err := s.RemoveObject(file)
	if err != nil { t.Errorf("Failed to remove object: %s", err) }
}
