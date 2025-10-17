package format

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/rotisserie/eris"

	"github.com/zq-xu/gotools/logx"
)

var skipDirs = []string{"vendor", "testdata"}

func FormatGoCodeInDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if slices.Contains(skipDirs, info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		return FormatGoCodeInFile(path)
	})
}

// FormatGoCodeInFile formats the import block of the specified Go file.
// It prints the file path if the file content is modified.
func FormatGoCodeInFile(path string) error {
	if !strings.HasSuffix(path, ".go") {
		return nil
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return eris.Wrapf(err, "read file %s", path)
	}

	changed, formatted, err := formatGoCode(src)
	if err != nil {
		return eris.Wrapf(err, "format %s", path)
	}

	if !changed {
		return nil
	}

	if err := os.WriteFile(path, formatted, 0o644); err != nil {
		return eris.Wrapf(err, "write file %s", path)
	}

	logx.Logger.Infof("Formatted: %s", path)
	return nil
}

//
// ---------------- Core Logic ----------------
//

// formatGoCode rearranges the import block and applies gofmt formatting.
// Returns whether the content was changed, the formatted content, and an error if any.
func formatGoCode(src []byte) (changed bool, formatted []byte, err error) {
	original := string(src)

	organized, err := organizeImports(original)
	if err != nil {
		return false, nil, err
	}

	formatted, err = format.Source([]byte(organized))
	if err != nil {
		return false, nil, eris.Wrap(err, "go/format source")
	}

	if bytes.Equal(src, formatted) {
		return false, formatted, nil
	}
	return true, formatted, nil
}

// organizeImports extracts the import block, organizes it into groups, and returns the updated source.
func organizeImports(src string) (string, error) {
	start, end := findImportRange(src)
	if start == -1 {
		return src, nil // No import block
	}

	lines := parseImportLines(src[start:end])
	std, third, local := classifyImports(lines)
	block := buildImportBlock(std, third, local)

	// Reassemble source with the new import block
	return src[:start] + block + src[end:], nil
}

//
// ---------------- Import Parsing and Construction ----------------
//

// findImportRange returns the start and end indices of the import block (including parentheses).
func findImportRange(src string) (start, end int) {
	start = strings.Index(src, "import (")
	if start == -1 {
		return -1, -1
	}
	endRel := strings.Index(src[start:], ")")
	if endRel == -1 {
		return -1, -1
	}
	return start, start + endRel + 1
}

// parseImportLines extracts individual import lines from an import block, trimming whitespace.
func parseImportLines(block string) []string {
	block = block[len("import (") : len(block)-1]
	var lines []string
	for _, l := range strings.Split(block, "\n") {
		if trimmed := strings.TrimSpace(l); trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}

// classifyImports groups import lines into standard library, third-party, and local project imports.
func classifyImports(lines []string) (std, third, local []string) {
	for _, l := range lines {
		switch {
		case strings.Contains(l, `"`+GoProject):
			local = append(local, l)
		case strings.Contains(l, "."):
			third = append(third, l)
		default:
			std = append(std, l)
		}
	}
	sort.Strings(std)
	sort.Strings(third)
	sort.Strings(local)
	return
}

// buildImportBlock constructs a new import block from the three groups.
func buildImportBlock(std, third, local []string) string {
	var b strings.Builder
	b.WriteString("import (\n")
	writeLines(&b, std)
	writeLines(&b, third)
	writeLines(&b, local)
	b.WriteString(")\n")
	return b.String()
}

// writeLines writes the import lines into the buffer with indentation and a separating newline.
func writeLines(b *strings.Builder, lines []string) {
	for _, l := range lines {
		b.WriteString("\t" + l + "\n")
	}
	if len(lines) > 0 {
		b.WriteString("\n")
	}
}
