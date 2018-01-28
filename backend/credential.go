package backend

import "time"

// Represents contents of a crypt file
type CryptFile struct {
	Credentials []Credential `json:"credentials"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CreatedAt   time.Time    `json:"created_at"`
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
