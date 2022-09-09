package serverless

import (
	"crypto/md5"
	"testing"
)

func TestHashDirectory(t *testing.T) {
	// No this is not a real test and is only here to help run the function
	if out, err := hashDirectory("..", md5.New()); err != nil {
		t.Logf("error: %s\n", err)
	} else {
		t.Log(out)
	}

}
