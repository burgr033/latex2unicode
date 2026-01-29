package converter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConverterNotEmpty(t *testing.T) {
	c, err := New()
	require.NoError(t, err)
	require.NotNil(t, c)
	require.NotEmpty(t, c.db)
}

func TestConvertString(t *testing.T) {
	tests := map[string]struct {
		st       styleTable
		input    string
		expected string
	}{
		"runtime complexity": {
			st:       superscriptTable,
			input:    `O(n\log{}n)`,
			expected: "O(nlog n)",
		},
		"integral": {
			input:    `\int g(x^2)dx = \pi e^{ix}`,
			expected: "∫ g(x²)dx = 𝜋 eⁱˣ",
		},
		"pi": {
			input:    `\pi`,
			expected: "𝜋",
		},
		"frac": {
			input:    `\frac{1}{4}`,
			expected: "1⁄4",
		},
		"greek": {
			input:    `\alpha + \beta = \gamma`,
			expected: `𝛼 + 𝛽 = 𝛾`,
		},
		"superscript digits": {
			st:       superscriptTable,
			input:    "^{123}",
			expected: "¹²³",
		},
		"subscript digits": {
			st:       subscriptTable,
			input:    "_{123}",
			expected: "₁₂₃",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c, err := New()
			require.NoError(t, err)
			require.Equal(t, tt.expected, c.ConvertString(tt.input))
		})
	}
}

func TestConvertMarkdownString(t *testing.T) {
	tests := map[string]struct {
		st       styleTable
		input    string
		expected string
	}{
		"runtime complexity": {
			st:       superscriptTable,
			input:    `$O(n\log{}n)$`,
			expected: "O(nlog n)",
		},
		"integral": {
			input:    `$$\int g(x^2)dx = \pi e^{ix}$$`,
			expected: "∫ g(x²)dx = 𝜋 eⁱˣ",
		},
		"pi": {
			input:    `$\pi$`,
			expected: "𝜋",
		},
		"frac": {
			input:    `$\frac{1}{4}$`,
			expected: "1⁄4",
		},

		"greek": {
			input:    `$$\alpha + \beta = \gamma$$`,
			expected: `𝛼 + 𝛽 = 𝛾`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c, err := New()
			require.NoError(t, err)
			require.Equal(t, tt.expected, c.ConvertMarkdown(tt.input))
		})
	}
}
