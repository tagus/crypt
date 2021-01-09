package utils

// FatalIf will panic if the given error is not nil
func FatalIf(err error) {
	if err != nil {
		panic(err)
	}
}
