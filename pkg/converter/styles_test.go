package converter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStyleTableApply(t *testing.T) {
	tests := map[string]struct {
		st       styleTable
		input    string
		expected string
	}{
		"superscript digits": {
			st:       superscriptTable,
			input:    "123",
			expected: "¹²³",
		},
		"subscript digits": {
			st:       subscriptTable,
			input:    "123",
			expected: "₁₂₃",
		},
		"unmapped characters pass through": {
			st:       superscriptTable,
			input:    "hello world!",
			expected: "ʰᵉˡˡᵒ ʷᵒʳˡᵈ!",
		},
		"empty string": {
			st:       superscriptTable,
			input:    "",
			expected: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.st.applyStyle(tt.input))
		})
	}
}

