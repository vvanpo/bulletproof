package main

import (
	"code.google.com/p/rsc/fuse"
	"flag"
	"log"
)

func main() {
	flag.Parse()
	mountPoint := flag.Arg(1)
	if mountPoint == "" {
		log.Fatal("Missing mount point")
	}
	c, err := fuse.Mount(mountPoint)
	if err != nil {
		log.Fatal(err)
	}
	c.Serve(FS{})
}

type FS struct{}

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
