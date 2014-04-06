package bp

import (
	//"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
//	pathname, _ := os.Getwd()
	pathname := "/home/victor"
	t.Logf("Using dir '%s'", pathname)
	c, err := GetConfig(pathname)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Configuration:\n%v", c)
}
