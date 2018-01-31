package creds

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	spaces = regexp.MustCompile(`\s+`)
)

// Represents contents of a crypt file
type Crypt struct {
	Credentials map[string]Credential `json:"credentials"`
	UpdatedAt   int64                 `json:"updated_at"`
	CreatedAt   int64                 `json:"created_at"`
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

// Command line output for a credential.
// Redacts the password
func (c Credential) PrintCredential() {
	color.Green("%s (%s)", c.Service, c.Description)

	fmt.Printf("%s: %s\n", color.BlueString("email"), normalizeField(c.Email))
	fmt.Printf("%s: %s\n", color.BlueString("username"), normalizeField(c.Username))
	fmt.Printf("%s: [redacted]\n", color.BlueString("password"))
}

// Returns 'N/A' for empty string
func normalizeField(field string) string {
	if field == "" {
		return color.WhiteString("N/A")
	} else {
		return field
	}
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

func (c *Crypt) GetCreatedAt() time.Time {
	return time.Unix(c.CreatedAt, 0)
}

func (c *Crypt) GetUpdatedAt() time.Time {
	return time.Unix(c.UpdatedAt, 0)
}

func (c *Crypt) GetJson() ([]byte, error) {
	str, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return str, nil
}

func (c *Crypt) IsValid(service string) bool {
	key := normalizeName(service)
	_, ok := c.Credentials[key]
	return ok
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
