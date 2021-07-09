package creds

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/tagus/crypt/internal/utils"
)

func PrintCredential(cred *Credential) {
	data := [][]string{
		{"id", cred.Id},
		{"email", normalizeField(cred.Email)},
		{"username", normalizeField(cred.Username)},
		{"pwd", "// redacted //"},
		{"tags", fmt.Sprintf("%v", cred.Tags)},
	}

	caption := fmt.Sprintf("%s: %s", cred.Service, cred.Description)
	utils.PrintTable(data, utils.TableOpts{
		Headers: []string{"field", "value"},
		Caption: caption,
	})
}

// Returns 'N/A' for empty string
func normalizeField(field string) string {
	if field == "" {
		return color.WhiteString("N/A")
	}
	return field
}
