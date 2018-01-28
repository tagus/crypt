package backend

import (
	"os"

	"github.com/fatih/color"
)

func PrintErrorAndExit(err error) {
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}
