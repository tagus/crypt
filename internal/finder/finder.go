// Package finder exposes an interface to query for credentials based on
// varying params such as service name, description, questions etc.
// Since we're not dealing with a lot of data, we will use linear search
// for everything however future improvements can include creating some
// search indices.
package finder

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/tagus/crypt/internal/crypt"
)

type Finder struct {
	cr  *crypt.Crypt
	idx bleve.Index
}

func New(cr *crypt.Crypt) (*Finder, error) {
	mapping := bleve.NewIndexMapping()
	idx, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, err
	}

	for _, cred := range cr.Credentials {
		idx.Index(cred.Id, cred)
	}
	return &Finder{
		cr:  cr,
		idx: idx,
	}, err
}

func (f *Finder) Filter(query string) ([]*crypt.Credential, error) {
	search := bleve.NewSearchRequest(bleve.NewMatchQuery(query))
	results, err := f.idx.Search(search)
	if err != nil {
		return nil, err
	}

	var matches []*crypt.Credential
	for _, cred := range results.Hits {
		matches = append(matches, f.cr.GetCredential(cred.ID))
	}
	return matches, err
}

func (f *Finder) Find(query string) (*crypt.Credential, error) {
	matches, err := f.Filter(query)
	if err != nil {
		return nil, err
	}
	if len(matches) > 0 {
		return matches[0], nil
	}
	return nil, nil
}
