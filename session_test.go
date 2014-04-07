package bp

import (
	//"os"
	"testing"
)

func TestNewSession(t *testing.T) {
//	pathname, _ := os.Getwd()
	pathname := "/home/victor/Dropbox/workspace/"
	t.Logf("Using dir '%s'", pathname)
	NewSession(pathname)
}
