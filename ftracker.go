package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"crypto/md5"
	"fmt"
	"gopkg.in/yaml.v1"
	"log"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

// A representation of an inode, i.e. a unique file object
type node struct {
	size    int64 // Size in bytes
	mode    os.FileMode
	modTime time.Time
	seLinux string         // SELinux context, if one exists
	hash    [md5.Size]byte // MD5 hash of file data
	links   []*link        // Back-references to files.  This is ignored if node is a directory
}

// Represents a file, i.e. a link to an inode
type link struct {
	// Base name
	name string
	*node
}

func (l link) isSymlink() bool {
	if l.node.mode|os.ModeSymlink != 0 {
		return true
	}
	return false
}

func (l *link) GetYAML() (string, interface{}) {
	value := make(map[string]interface{})
	value["Size"] = l.node.size
	value["Hash"] = fmt.Sprintf("%x", l.node.hash)
	value["Modified"] = l.node.modTime
	// TODO:  file mode
	return "", value
}

type dir struct {
	// A directory is also just a link, but most systems don't support multi-linking directories and neither will we
	name string
	*node
	// List of files
	files []link
	// List of symlinks within this directory, mapping the symlink files themselves to a pathname
	symlinks map[link]string
	// List of directories nested within this directory
	dirs []dir
}

func (d *dir) GetYAML() (string, interface{}) {
	value := make(map[string]interface{})
	for _, l := range d.files {
		value[l.name] = l
	}
	for l, v := range d.symlinks {
		value[l.name] = []interface{}{"-> " + v, l}
	}
	for _, ds := range d.dirs {
		value[ds.name] = []interface{}{ds.node, ds}
	}
	return "", value
}

// Per-instance file tree object
type Instance struct {
	// Absolute pathname to root
	path string
	// Root directory
	root *dir
	// Absolute pathname mapping to fsnotify watcher objects
	watchers map[string]fsnotify.Watcher
}

func (i *Instance) GetYAML() (string, interface{}) {
	tag, value := i.root.GetYAML()
	return tag, value
}

// NewInstance sets up the data structure for the given pathname by reading the configuration files in that path
func NewInstance(path string) *Instance {
	i := new(Instance)
	i.path = path
	i.root = new(dir)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Print(err)
	}
	var cfg map[string]interface{}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("%#v\n", cfg)
	return i
}

func main() {
	root := "/home/victor"
	i := NewInstance(root)
	println(i)
}
