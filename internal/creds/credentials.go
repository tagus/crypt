package creds

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/sugatpoudel/crypt/internal/utils"
)

// Crypt represents contents of a crypt file
type Crypt struct {
	Credentials map[string]Credential `json:"credentials"`
	UpdatedAt   int64                 `json:"updated_at"`
	CreatedAt   int64                 `json:"created_at"`
}

// Credential houses all pertinent information for a given service
type Credential struct {
	Service           string             `json:"service"`
	Email             string             `json:"email"`
	Username          string             `json:"username"`
	Password          string             `json:"password"`
	Description       string             `json:"description"`
	SecurityQuestions []SecurityQuestion `json:"security_questions"`
	CreatedAt         int64              `json:"created_at"`
	UpdatedAt         int64              `json:"updated_at"`
}

// SecurityQuestion holds a security question and its answer
type SecurityQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// GetCreatedAt retrieves the credential's creation time
func (c *Credential) GetCreatedAt() time.Time {
	return time.Unix(c.CreatedAt, 0)
}

// GetUpdatedAt retrieves the credential's last update time
func (c *Credential) GetUpdatedAt() time.Time {
	return time.Unix(c.UpdatedAt, 0)
}

// PrintCredential prints the credentials while redacting the password
func (c Credential) PrintCredential() {
	data := [][]string{
		[]string{"Email", normalizeField(c.Email)},
		[]string{"Username", normalizeField(c.Username)},
		[]string{"Password", "[redacted]"},
	}

	fmt.Println()
	caption := fmt.Sprintf("%s: %s", c.Service, c.Description)
	utils.PrintTable(data, []string{"field", "value"}, caption)
}

// Returns 'N/A' for empty string
func normalizeField(field string) string {
	if field == "" {
		return color.WhiteString("N/A")
	}
	return field
}

// SetCredential is used to add/update the list of credentials
func (c *Crypt) SetCredential(cred Credential) {
	key := utils.NormalizeString(cred.Service)
	c.Credentials[key] = cred
}

// FindCredential finds the Credential struct corresponding to the given service name.
func (c *Crypt) FindCredential(service string) *Credential {
	key := utils.NormalizeString(service)
	if cred, ok := c.Credentials[key]; ok {
		return &cred
	}
	return nil
}

// RemoveCredential removes the given service from the Store
func (c *Crypt) RemoveCredential(service string) {
	key := utils.NormalizeString(service)
	delete(c.Credentials, key)
}

// GetCreatedAt retrieves the crypt's creation time
func (c *Crypt) GetCreatedAt() time.Time {
	return time.Unix(c.CreatedAt, 0)
}

// GetUpdatedAt retrieves the crypt's last update time
func (c *Crypt) GetUpdatedAt() time.Time {
	return time.Unix(c.UpdatedAt, 0)
}

// GetJSON marshals the Crypt object into json
func (c *Crypt) GetJSON() ([]byte, error) {
	str, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return str, nil
}

// IsValid determines if the given service exists the Crypt
func (c *Crypt) IsValid(service string) bool {
	key := utils.NormalizeString(service)
	_, ok := c.Credentials[key]
	return ok
}

// FromJSON unmarshalls a Crypt object from json
func FromJSON(data []byte) (*Crypt, error) {
	var crypt Crypt
	if err := json.Unmarshal(data, &crypt); err != nil {
		return nil, err
	}
	return &crypt, nil
}

// GetSuggestions get the closest services to the given input service
// based on the edit distance of the service name
func (c *Crypt) GetSuggestions(service string) []string {
	var suggestions []string
	for _, v := range c.Credentials {
		ld := utils.CalculateLevenshteinDistance(service, v.Service)
		if ld < 5 {
			suggestions = append(suggestions, v.Service)
		}
	}
	return suggestions
}

// Len determines the number of credentials that are part of this crypt
func (c *Crypt) Len() int {
	return len(c.Credentials)
}
