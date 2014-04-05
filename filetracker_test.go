package filetracker

import (
	//"os"
	"testing"
)

func TestNewSession(t *testing.T) {
//	pathname, _ := os.Getwd()
	pathname := "/home/victor"
	t.Logf("Using dir '%s'", pathname)
	NewSession(pathname)
}
