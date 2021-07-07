// Package finder exposes an interface to query for credentials based on
// varying params such as service name, description, questions etc.
// Since we're not dealing with a lot of data, we will use linear search
// for everything however future improvements can include creating some
// search indices.
package finder

import (
	"strings"

	"github.com/tagus/crypt/internal/creds"
)

type Finder struct {
	crypt *creds.Crypt
}

func New(crypt *creds.Crypt) *Finder {
	return &Finder{
		crypt: crypt,
	}
}

func (f *Finder) Filter(query string) []*creds.Credential {
	var matches []*creds.Credential
	for _, cred := range f.crypt.Credentials {
		if strings.EqualFold(query, cred.Service) {
			matches = append(matches, cred)
		}
	}
	return matches
}

func (f *Finder) Find(query string) *creds.Credential {
	matches := f.Filter(query)
	if len(matches) > 0 {
		return matches[0]
	}
	return nil
}
