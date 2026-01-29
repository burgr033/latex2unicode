package converter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadEmbeddedSymbolDB(t *testing.T) {
	db, err := loadEmbeddedSymbolDB()
	require.NoError(t, err)
	require.NotNil(t, db)

	// Assert something meaningful
	require.NotEmpty(t, db.charToSymbol)
	require.NotEmpty(t, db.texToChar)
	require.NotEmpty(t, db.unicodeMathToChar)
	require.NotEmpty(t, db.texPattern)
}
