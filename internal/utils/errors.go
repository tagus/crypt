package utils

import (
	"fmt"
	"os"

	"github.com/sugatpoudel/crypt/internal/env"
)

// FatalIf will panic if the given error is not nil
func FatalIf(err error) {
	if err != nil {
		fmt.Println("an error occurred: ", err)
		if env.IsDev() {
			panic(err)
		}
		os.Exit(1)
	}
}
