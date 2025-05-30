package utils

import (
	"fmt"
	"github.com/tagus/crypt/internal/repos"
	"github.com/tagus/mango"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	spaces = regexp.MustCompile(`\s+`)
)

// NormalizeString normalizes the provided string to remove spaces
// and be lowercase
func NormalizeString(str string) string {
	normalized := strings.ToLower(str)
	normalized = spaces.ReplaceAllString(normalized, "_")
	return normalized
}

// TableOpts defines options for the table writer
type TableOpts struct {
	Headers []string
	Caption string
}

// PrintTable prints the given table in a formatted table.
func PrintTable(data [][]string, opts TableOpts) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoMergeCells(true)
	if len(opts.Headers) > 0 {
		table.SetHeader(opts.Headers)
	}
	if opts.Caption != "" {
		table.SetCaption(true, opts.Caption)
	}
	table.SetAutoMergeCells(false)
	table.AppendBulk(data)
	table.Render()
}

func PrintCredential(cred *repos.Credential) {
	data := [][]string{
		{"id", cred.ID},
		{"service", getFallbackString(cred.Service)},
		{"email", getFallbackString(cred.Email)},
		{"username", getFallbackString(cred.Username)},
		{"created_at", mango.FormatSimpleDate(&cred.CreatedAt)},
		{"updated_at", mango.FormatSimpleDate(&cred.UpdatedAt)},
		{"accessed_at", mango.FormatSimpleDate(cred.AccessedAt)},
		{"accessed_count", strconv.Itoa(cred.AccessedCount)},
		{"tags", fmt.Sprintf("%v", cred.Tags)},
	}

	caption := fmt.Sprintf("%s: %s", cred.Service, cred.Description)
	PrintTable(data, TableOpts{
		Headers: []string{"field", "value"},
		Caption: caption,
	})
}

func getFallbackString(field string) string {
	if field == "" {
		return mango.ColorizeWhite("N/A")
	}
	return field
}
