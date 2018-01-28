package backend

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
