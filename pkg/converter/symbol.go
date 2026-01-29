package converter

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

type symbol struct {
	char        rune
	texCommand  string
	unicodeMath string
	description string
}

type symbolDB struct {
	texToChar         map[string]rune
	unicodeMathToChar map[string]rune
	charToSymbol      map[rune]*symbol
	texPattern        *regexp.Regexp
}

// loadEmbeddedSymbolDB loads from the embedded data
func loadEmbeddedSymbolDB() (*symbolDB, error) {
	return parseSymbolDB(strings.NewReader(embeddedSymbolData))
}

// parseSymbolDB parses the embedded symbol database
func parseSymbolDB(r io.Reader) (*symbolDB, error) {
	db := &symbolDB{
		texToChar:         make(map[string]rune),
		unicodeMathToChar: make(map[string]rune),
		charToSymbol:      make(map[rune]*symbol),
	}

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		fields := strings.Split(line, "^")

		if len(fields) < 8 {
			continue
		}

		var cp int64
		_, err := fmt.Sscanf(fields[0], "%X", &cp)
		if err != nil {
			continue
		}
		codepoint := rune(cp)

		texCmd := strings.TrimSpace(fields[2])
		unicodeMathCmd := strings.TrimSpace(fields[3])
		description := strings.TrimSpace(fields[7])

		sym := &symbol{
			char:        codepoint,
			texCommand:  texCmd,
			unicodeMath: unicodeMathCmd,
			description: description,
		}

		if texCmd != "" {
			cmdName := strings.TrimPrefix(texCmd, "\\")
			db.texToChar[cmdName] = codepoint
		}

		if unicodeMathCmd != "" {
			cmdName := strings.TrimPrefix(unicodeMathCmd, "\\")
			if texCmd == "" || cmdName != strings.TrimPrefix(texCmd, "\\") {
				db.unicodeMathToChar[cmdName] = codepoint
			}
		}

		db.charToSymbol[codepoint] = sym
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	db.buildTeXPattern()
	return db, nil
}

func (db *symbolDB) buildTeXPattern() {
	commands := make([]string, 0, len(db.texToChar))
	for cmd := range db.texToChar {
		commands = append(commands, regexp.QuoteMeta(cmd))
	}

	sort.Slice(commands, func(i, j int) bool {
		return len(commands[i]) > len(commands[j])
	})

	pattern := `\\(?:` + strings.Join(commands, "|") + `)`
	db.texPattern = regexp.MustCompile(pattern)
}
