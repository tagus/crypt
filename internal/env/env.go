package env

var (
	dev bool = false
)

// IsDev determines if the dev flag is true in which case callers can
// perform dev specific action or short circuit from unnecessary
// operations
func IsDev() bool {
	return dev
}

// SetDev sets the dev flag to the given value
func SetDev(val bool) {
	dev = val
}
