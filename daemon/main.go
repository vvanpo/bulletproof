
package main

import (
	"log"
)

func main() {
	s := newSession("/tmp")
	err := s.start()
	if err != nil {
		log.Fatal(err)
	}
	return
}
