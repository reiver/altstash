package cfg

import (
	"os"
)

var userHomeDir string

func init() {
	var err error
	userHomeDir, err = os.UserHomeDir()
	if nil != err {
		panic(err)
	}
}
