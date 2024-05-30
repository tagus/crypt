package repos

type QueryCryptsFilter struct {
	Name string
	ID   string
}

type QueryCredentialsFilter struct {
	ID                   string
	CryptID              string
	Service              string
	IncrementAccessCount bool
	Limit                int
	Tag                  string
}
