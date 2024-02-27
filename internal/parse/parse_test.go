package parse

import (
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestParseEnvFile(t *testing.T) {

	type testStruct struct {
		name     string
		input    string
		expected map[string]string
	}
	testCases := []testStruct{
		{"multiline", `ABCD=ABCD
1234=1234`, map[string]string{"ABCD": "ABCD", "1234": "1234"}},
		{"quotes", `1234="A quoted string"`, map[string]string{"1234": "A quoted string"}},
		{"backticks", "1234=\x60A backticked string\x60", map[string]string{"1234": "A backticked string"}},
		{"quoted newlines", `1234="a string
with a newline in it"`, map[string]string{"1234": `a string
with a newline in it`}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := EnvFile(tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
