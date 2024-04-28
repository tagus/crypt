package repos

import "time"

type Credential struct {
	ID            string     `json:"id"`
	Service       string     `json:"service"`
	Domains       []string   `json:"domains"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	Password      string     `json:"password"`
	Description   string     `json:"description"`
	Details       *Details   `json:"details"`
	Tags          []string   `json:"tags"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	AccessedAt    *time.Time `json:"accessed_at"`
	AccessedCount int        `json:"accessed_count"`
	Version       int        `json:"version"`
}

type Details struct {
	SecurityQuestions []SecurityQuestion `json:"security_questions"`
}

type SecurityQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
