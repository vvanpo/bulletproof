package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"crypto/md5"
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

func (n *node) GetYAML() (string, interface{}) {
	var value [3]interface{}
	value[0] = n.size
	value[1] = fmt.Sprintf("%x", n.hash)
	value[2] = n.modTime
	// TODO:  file mode
	return "!!seq", value
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
	value := make([]interface{}, len(d.files) + len(d.symlinks) + len(d.dirs))
	i := 0
	for _, l := range d.files {
		value[i] = map[string]interface{}{l.name: l.node}
		i++
	}
	for l, v := range d.symlinks {
		value[i] = map[string]interface{}{l.name + " -> " + v: l.node}
		i++
	}
	for _, ds := range d.dirs {
		value[i] = map[string]interface{}{ds.name: ds}
		i++
	}
	return "!!seq", value
}

func (d *dir) SetYAML(tag string, value interface{}) bool {
	if v, ok := value.([]interface{}); ok {
		for _, i := range v {
			for k, t := range i.(map[interface{}]interface{}) {
				fmt.Printf("%v\t%v\n", k, t)
			}
		}
		return true
	}
	return false
}

// Per-instance file tree object
type Session struct {
	// Absolute pathname to root
	pathname string
	// Root directory
	root *dir
	// Absolute pathname mapping to fsnotify watcher objects
	watchers map[string]fsnotify.Watcher
}

func (s *Session) GetYAML() (tag string, value interface{}) {
	return "", s.root
}

func (s *Session) SetYAML(tag string, value interface{}) bool {
	return s.root.SetYAML(tag, value)
}

// NewSession sets up the data structure for the given pathname by reading the configuration files in that path
func NewSession(pathname string) *Session {
	configPath, err := filepath.Abs("config.yaml")
	//configPath := filepath.Join(pathname, ".bp/config")
	if err != nil {
		log.Fatal(err)
	}
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	s := new(Session)
	s.pathname = pathname
	s.root = new(dir)
	err = yaml.Unmarshal(configFile, s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n\n%#v\n", s.root)
	return s
}

func main() {
	root := "/home/victor"
	i := NewSession(root)
	println(i)
}
