package creds

import (
	"encoding/json"
	"time"

	"github.com/teris-io/shortid"
)

// Crypt represents contents of a crypt file
type Crypt struct {
	Credentials map[string]*Credential `json:"credentials"`
	UpdatedAt   int64                  `json:"updated_at"`
	CreatedAt   int64                  `json:"created_at"`
}

// Credential houses all pertinent information for a given service
type Credential struct {
	Id                string             `json:"id"`
	Service           string             `json:"service"`
	Domains           []string           `json:"domains"`
	Email             string             `json:"email"`
	Username          string             `json:"username"`
	Password          string             `json:"password"`
	Description       string             `json:"description"`
	SecurityQuestions []SecurityQuestion `json:"security_questions"`
	CreatedAt         int64              `json:"created_at"`
	UpdatedAt         int64              `json:"updated_at"`
}

// SecurityQuestion models a question/answer pair
type SecurityQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (c *Crypt) Len() int {
	return len(c.Credentials)
}

func (c *Credential) GetCreatedAt() time.Time {
	return time.Unix(c.CreatedAt, 0)
}

func (c *Credential) GetUpdatedAt() time.Time {
	return time.Unix(c.UpdatedAt, 0)
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

func FromJSON(data []byte) (*Crypt, error) {
	var crypt Crypt
	if err := json.Unmarshal(data, &crypt); err != nil {
		return nil, err
	}
	return &crypt, nil
}
