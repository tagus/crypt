package creds

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/teris-io/shortid"
)

// Crypt represents contents of a crypt file
type Crypt struct {
	Credentials map[string]*Credential `json:"credentials"`
	UpdatedAt   int64                  `json:"updated_at"`
	CreatedAt   int64                  `json:"created_at"`
}

func FromJSON(data []byte) (*Crypt, error) {
	var crypt Crypt
	if err := json.Unmarshal(data, &crypt); err != nil {
		return nil, err
	}
	return &crypt, nil
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
	str, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return str, nil
}

// Credential houses all pertinent information for a given service
type Credential struct {
	Id          string   `json:"id"`
	Service     string   `json:"service"`
	Domains     []string `json:"domains"`
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Description string   `json:"description"`
	Details     []Detail `json:"details"`
	Tags        []string `json:"tags"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

// Detail models an arbitrary key value pair
type Detail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c *Credential) GetCreatedAt() time.Time {
	return time.Unix(c.CreatedAt, 0)
}

func (c *Credential) GetUpdatedAt() time.Time {
	return time.Unix(c.UpdatedAt, 0)
}

func (c *Credential) Merge(cred *Credential) *Credential {
	dirty := false
	if cred.Service != "" {
		val := strings.TrimSpace(cred.Service)
		if val != "" {
			dirty = true
			c.Service = val
		}
	}
	if cred.Email != "" {
		val := strings.TrimSpace(cred.Email)
		if val != "" {
			dirty = true
			c.Email = val
		}
	}
	if cred.Username != "" {
		val := strings.TrimSpace(cred.Username)
		if val != "" {
			dirty = true
			c.Username = val
		}
	}
	if cred.Password != "" {
		val := strings.TrimSpace(cred.Password)
		if val != "" {
			dirty = true
			c.Password = val
		}
	}
	if cred.Description != "" {
		val := strings.TrimSpace(cred.Description)
		if val != "" {
			dirty = true
			c.Description = val
		}
	}
	if dirty {
		c.UpdatedAt = time.Now().Unix()
	}
	return c
}
