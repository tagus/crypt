package backend

import "time"

// Represents contents of a crypt file
type CryptFile struct {
	Credentials []Credential
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// Houses all pertinent information for a given service
type Credential struct {
	Service     string
	Email       string
	Username    string
	Password    []byte
	Description string
	Securities  []SecurityQuestion
}

// Security question tuple
type SecurityQuestion struct {
	Question string
	Answer   string
}
