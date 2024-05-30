package repos

import "time"

type Crypt struct {
	ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	UpdatedAt              time.Time `json:"updated_at"`
	CreatedAt              time.Time `json:"created_at"`
	TotalActiveCredentials int       `json:"total_active_credentials"`
	HashedPassword         []byte    `json:"hashed_password"`
	Signature              []byte    `json:"signature"`
}
