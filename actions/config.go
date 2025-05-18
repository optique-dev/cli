package actions

import (
	"fmt"

	"github.com/optique-dev/optique"
)

func GetConfig() string {
	root, err := optique.FindOptiqueJson()
	if err != nil {
		optique.Error("You are not in an optique project. To create a new one, run `optique init`")
	}
	_ = fmt.Sprint("%s/config/config.go", root)
	// return parseConfigFile(filename)
	return ""
}

// func parseConfigFile(path string) string {
//
// }
