package main

import (
	"code.google.com/p/rsc/fuse"
	"flag"
	"log"
	"os"
	"path"
)

/*
 * main parses the command-line in the form of:
 * ./fs <file> <mount-point>
 */
func main() {
	flag.Parse()
	f := new(FS)
	master := path.Clean(flag.Arg(1))
	if master == "." {
		log.Fatal("Missing file")
	}
	mountPoint := flag.Arg(2)
	if mountPoint == "" {
		log.Fatal("Missing mount point")
	}
	if !f.validate() {
		log.Fatal("Invalid file")
	}
	var err error
	f.fd, err = os.OpenFile(master, os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	c, err := fuse.Mount(mountPoint)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Serve(FS{}); err != nil {
		log.Fatal(err)
	}
}

type FS struct{
	fd *os.File			// master file descriptor
}

// validate ensures the master file is a valid filesystem
func (f FS) validate() bool {
	return true
}

func (FS) Root() (fuse.Node, fuse.Error) {
	return nil, nil
}

/*
func (FS) Statfs(r *fuse.StatfsResponse, intr fuse.Intr) fuse.Error {
}
*/

/*
type Handle struct{}

func (Handle) Flush() fuse.Error {
}
*/

type Node struct {
	attr fuse.Attr
}

func (n *Node) Attr() fuse.Attr {
	return n.attr
}

