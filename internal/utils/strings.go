package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/tagus/crypt/internal/repos"
	"os"
	"regexp"
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

// CalculateLevenshteinDistance measures the similarity between the two
// given strings.
func CalculateLevenshteinDistance(a, b string) int {
	aLen, bLen := len(a), len(b)

	if aLen == 0 {
		return bLen
	} else if bLen == 0 {
		return aLen
	}

	grid := make([][]int, aLen+1)
	for i := range grid {
		grid[i] = make([]int, bLen+1)
	}

	for i := 1; i <= aLen; i++ {
		grid[i][0] = i
	}

	for j := 1; j <= bLen; j++ {
		grid[0][j] = j
	}

	for j := 1; j <= bLen; j++ {
		for i := 1; i <= aLen; i++ {
			var cost int
			if a[i-1] == b[j-1] {
				cost = 0
			} else {
				cost = 1
			}
			grid[i][j] = Min(grid[i-1][j]+1, grid[i][j-1]+1, grid[i-1][j-1]+cost)
		}
	}

	return grid[aLen][bLen]
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
		{"created_at", FormatDate(cred.CreatedAt)},
		{"updated_at", FormatDate(cred.UpdatedAt)},
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
		return color.WhiteString("N/A")
	}
	return field
}
