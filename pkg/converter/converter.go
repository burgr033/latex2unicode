// Package converter
package converter

import (
	"regexp"
	"strings"
)

// Converter handles LaTeX to Unicode conversion
type Converter struct {
	db *symbolDB
}

// New creates a new Converter by loading the symbol database
func New() (*Converter, error) {
	db, err := loadEmbeddedSymbolDB()
	if err != nil {
		return nil, err
	}
	return &Converter{db: db}, nil
}

// ConvertString converts LaTeX mathematical notation to Unicode
func (c *Converter) ConvertString(input string) string {
	result := convertFrac(input)
	result = convertMathFunctions(result)

	result = c.db.texPattern.ReplaceAllStringFunc(result, func(match string) string {
		cmd := strings.TrimPrefix(match, "\\")
		if char, ok := c.db.texToChar[cmd]; ok {
			return string(char)
		}
		return match
	})

	result = convertSuperscripts(result)
	result = convertSubscripts(result)

	return result
}

// ConvertMarkdown strips markdown math delimiters ($ and $$) and converts the LaTeX content
func (c *Converter) ConvertMarkdown(input string) string {
	input = strings.TrimPrefix(input, "$$")
	input = strings.TrimSuffix(input, "$$")
	input = strings.TrimPrefix(input, "$")
	input = strings.TrimSuffix(input, "$")

	return c.ConvertString(input)
}

func convertFrac(input string) string {
	re := regexp.MustCompile(`\\frac\{([^}]+)\}\{([^}]+)\}`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) == 3 {
			numerator := submatches[1]
			denominator := submatches[2]
			return numerator + "⁄" + denominator
		}
		return match
	})
}

func convertMathFunctions(input string) string {
	functions := []string{
		"arcsin", "arccos", "arctan", "arccot", "arcsec", "arccsc",
		"sinh", "cosh", "tanh", "coth",
		"log", "ln", "exp", "sin", "cos", "tan", "cot", "sec", "csc",
		"lim", "max", "min", "sup", "inf", "det", "dim", "ker", "deg",
		"gcd", "lcm", "Pr", "hom", "arg",
	}

	for _, fn := range functions {
		pattern := regexp.MustCompile(`\\` + fn + `\{([^}]*)\}`)
		input = pattern.ReplaceAllStringFunc(input, func(match string) string {
			re := regexp.MustCompile(`\\` + fn + `\{([^}]*)\}`)
			submatches := re.FindStringSubmatch(match)
			if len(submatches) == 2 {
				content := submatches[1]
				if content == "" {
					return fn + " "
				}
				return fn + " " + content + " "
			}
			return match
		})

		pattern = regexp.MustCompile(`\\` + fn + `([^a-zA-Z{])`)
		input = pattern.ReplaceAllString(input, fn+" $1")

		pattern = regexp.MustCompile(`\\` + fn + `$`)
		input = pattern.ReplaceAllString(input, fn)
	}

	return input
}

func convertSuperscripts(s string) string {
	re := regexp.MustCompile(`\^(\{[^}]+\}|[^{}\s\\])`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		content := match[1:]
		content = strings.Trim(content, "{}")
		return superscriptTable.applyStyle(content)
	})
}

func convertSubscripts(s string) string {
	re := regexp.MustCompile(`_(\{[^}]+\}|[^{}\s\\])`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		content := match[1:]
		content = strings.Trim(content, "{}")
		return subscriptTable.applyStyle(content)
	})
}
