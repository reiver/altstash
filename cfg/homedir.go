package cfg

import (
	"os"

	"codeberg.org/reiver/go-erorr"
)

var userHomeDir string

func init() {
	var err error
	userHomeDir, err = os.UserHomeDir()
	if nil != err {
		panic(erorr.Wrap(err, "could not determine where the user home directory is"))
	}
}
