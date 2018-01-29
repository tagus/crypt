package creds

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

var (
	spaces = regexp.MustCompile(`\s+`)
)

// Represents contents of a crypt file
type Crypt struct {
	Credentials map[string]Credential `json:"credentials"`
	UpdatedAt   time.Time             `json:"updated_at"`
	CreatedAt   time.Time             `json:"created_at"`
}

// Houses all pertinent information for a given service
type Credential struct {
	Service           string             `json:"service"`
	Email             string             `json:"email"`
	Username          string             `json:"username"`
	Password          string             `json:"password"`
	Description       string             `json:"description"`
	SecurityQuestions []SecurityQuestion `json:"security_questions"`
}

// Security question tuple
type SecurityQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// This is used to add/update the list of credentials
func (c *Crypt) SetCredential(cred Credential) {
	key := normalizeName(cred.Service)
	c.Credentials[key] = cred
}

// Returns the Credential struct corresponding to the given service name.
func (c *Crypt) FindCredential(service string) *Credential {
	key := normalizeName(service)
	if cred, ok := c.Credentials[key]; ok {
		return &cred
	}
	return nil
}

func (c *Crypt) RemoveCredential(service string) {
	delete(c.Credentials, service)
}

func (c *Crypt) GetCreatedAt() string {
	return c.CreatedAt.Format("Jan 2 15:04:05, 2006")
}

func (c *Crypt) GetUpdatedAt() string {
	return c.UpdatedAt.Format("Jan 2 15:04:05, 2006")
}

func (c *Crypt) GetJson() ([]byte, error) {
	str, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return str, nil
}

func FromJson(data []byte) (*Crypt, error) {
	var crypt Crypt
	if err := json.Unmarshal(data, &crypt); err != nil {
		return nil, err
	}
	return &crypt, nil
}

func normalizeName(name string) string {
	name = strings.ToLower(name)
	name = spaces.ReplaceAllString(name, "_")
	return name
}
