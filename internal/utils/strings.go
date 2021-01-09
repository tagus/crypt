package utils

import (
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

// FallbackStr returns the fallback if the given val is empty otherwise
// returns the normalized val
func FallbackStr(val, fallback string) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return fallback
	}
	return val
}
