package repos

import (
	"fmt"
	"time"
)

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

func (c *Credential) String() string {
	return fmt.Sprintf("%s.%s", c.ID, c.Service)
}

func (c *Credential) Clone() *Credential {
	return &Credential{
		ID:            c.ID,
		Service:       c.Service,
		Domains:       append([]string{}, c.Domains...),
		Email:         c.Email,
		Username:      c.Username,
		Password:      c.Password,
		Description:   c.Description,
		Details:       c.Details.Clone(),
		Tags:          append([]string{}, c.Tags...),
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		AccessedAt:    c.AccessedAt,
		AccessedCount: c.AccessedCount,
		Version:       c.Version,
	}
}

type Details struct {
	SecurityQuestions []SecurityQuestion `json:"security_questions"`
}

func (d *Details) Clone() *Details {
	if d == nil {
		return nil
	}
	sqs := make([]SecurityQuestion, len(d.SecurityQuestions))
	for i := range d.SecurityQuestions {
		sqs[i] = *d.SecurityQuestions[i].Clone()
	}
	return &Details{
		SecurityQuestions: sqs,
	}
}

type SecurityQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (s *SecurityQuestion) Clone() *SecurityQuestion {
	return &SecurityQuestion{
		Question: s.Question,
		Answer:   s.Answer,
	}
}
