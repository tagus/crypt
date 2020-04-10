package utils

import (
	"fmt"
	"os"

	"github.com/sugatpoudel/crypt/internal/env"
)

// FatalIf will panic if the given error is not nil
func FatalIf(err error) {
	if err != nil {
		fmt.Printf("an unexpected error occurred")
		if env.IsDev() {
			panic(err)
		}
		os.Exit(1)
	}
}
