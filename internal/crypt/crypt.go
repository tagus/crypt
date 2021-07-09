package crypt

import (
	"encoding/json"
	"time"

	"github.com/teris-io/shortid"
)

// Version represents the current version of the cryptfile + cli
const Version = "v1.3"

// Crypt represents contents of a crypt file
type Crypt struct {
	Credentials map[string]*Credential `json:"credentials"`
	UpdatedAt   int64                  `json:"updated_at"`
	CreatedAt   int64                  `json:"created_at"`
	Version     string                 `json:"version"`
}

func FromJSON(data []byte) (*Crypt, error) {
	var cr Crypt
	if err := json.Unmarshal(data, &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

func (c *Crypt) Len() int {
	return len(c.Credentials)
}

func (c *Crypt) SetCredential(cred *Credential) (*Credential, error) {
	if cred.Id == "" {
		var err error
		cred.Id, err = shortid.Generate()
		if err != nil {
			return nil, err
		}
	}
	c.Credentials[cred.Id] = cred
	return cred, nil
}

func (c *Crypt) GetCredential(id string) *Credential {
	return c.Credentials[id]
}

func (c *Crypt) RemoveCredential(id string) {
	delete(c.Credentials, id)
}

func (c *Crypt) GetCreatedAt() time.Time {
	return time.Unix(c.CreatedAt, 0)
}

func (c *Crypt) GetUpdatedAt() time.Time {
	return time.Unix(c.UpdatedAt, 0)
}

func (c *Crypt) GetJSON() ([]byte, error) {
	c.Version = Version
	str, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return str, nil
}
