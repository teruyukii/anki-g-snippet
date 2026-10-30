package libs

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func PrintAnkiCodeHTML(filepath string, startLine int, endLine int) (string, error) {
	if startLine > endLine || startLine < 1 {
		return "", fmt.Errorf("invalid line range: %d to %d", startLine, endLine)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %w", filepath, err)
	}
	defer file.Close()

	var codeBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	currentLine := 1

	for scanner.Scan() {
		if currentLine >= startLine && currentLine <= endLine {
			codeBuilder.WriteString(scanner.Text() + "\n")
		}
		if currentLine >= endLine {
			break
		}
		currentLine++
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning file %s: %w", filepath, err)
	}

	// Detect language from filename
	lexer := lexers.Match(filepath)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("github-dark")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(
		html.WithClasses(false),        // Force inline CSS styles so Anki needs no external CSS setup
		html.WithLineNumbers(true),     // Show line numbers
		html.BaseLineNumber(startLine), // Start at your specific line number
		html.LineNumbersInTable(false), // Use spans instead of tables (tables often break in Anki)
		html.TabWidth(4),
	)

	iterator, err := lexer.Tokenise(nil, codeBuilder.String())
	if err != nil {
		return "", err
	}

	var htmlBuffer bytes.Buffer
	if err := formatter.Format(&htmlBuffer, style, iterator); err != nil {
		return "", err
	}

	// Wrap in an Anki-friendly container that forces left alignment, rounded corners, and horizontal scroll
	ankiWrapper := fmt.Sprintf(
		`<div style="text-align: left; font-family: monospace; font-size: 14px; border-radius: 8px; overflow-x: auto; padding: 12px; background-color: #282a36;">%s</div>`,
		htmlBuffer.String(),
	)

	return ankiWrapper, nil
}
