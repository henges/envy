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
		{"multiline parsing with comments", `BASIC=basic

# previous line intentionally left blank
AFTER_LINE=after_line`, map[string]string{"BASIC": "basic", "AFTER_LINE": "after_line"}},
		{"defaults to empty", "EMPTY=", map[string]string{"EMPTY": ""}},
		{"empty single quotes", "EMPTY_SINGLE_QUOTES=''", map[string]string{"EMPTY_SINGLE_QUOTES": ""}},
		{"empty double quotes", `EMPTY_DOUBLE_QUOTES=""`, map[string]string{"EMPTY_DOUBLE_QUOTES": ""}},
		{"empty backticks", "EMPTY_BACKTICKS=\x60\x60", map[string]string{"EMPTY_BACKTICKS": ""}},
		{"single quotes are removed", "SINGLE_QUOTES='single_quotes'", map[string]string{"SINGLE_QUOTES": "single_quotes"}},
		{"spaces are not trimmed in single quotes", "SINGLE_QUOTES_SPACED='    single quotes    '", map[string]string{"SINGLE_QUOTES_SPACED": "    single quotes    "}},
		{"double quotes are removed", `DOUBLE_QUOTES="double_quotes"`, map[string]string{"DOUBLE_QUOTES": `double_quotes`}},
		{"spaces are not trimmed in double quotes", `DOUBLE_QUOTES_SPACED="    double quotes    "`, map[string]string{"DOUBLE_QUOTES_SPACED": `    double quotes    `}},
		{"double quotes work inside single quotes", `DOUBLE_QUOTES_INSIDE_SINGLE='double "quotes" work inside single quotes'`, map[string]string{"DOUBLE_QUOTES_INSIDE_SINGLE": `double "quotes" work inside single quotes`}},
		{"not sure what this is meant to test", `DOUBLE_QUOTES_WITH_NO_SPACE_BRACKET="{ port: $MONGOLAB_PORT}"`, map[string]string{"DOUBLE_QUOTES_WITH_NO_SPACE_BRACKET": "{ port: $MONGOLAB_PORT}"}},
		{"single quotes work inside double quotes", `SINGLE_QUOTES_INSIDE_DOUBLE="single 'quotes' work inside double quotes"`, map[string]string{"SINGLE_QUOTES_INSIDE_DOUBLE": "single 'quotes' work inside double quotes"}},
		{"backticks inside single works", "BACKTICKS_INSIDE_SINGLE='\x60backticks\x60 work inside single quotes'", map[string]string{"BACKTICKS_INSIDE_SINGLE": "\x60backticks\x60 work inside single quotes"}},
		{"backticks inside double works", "BACKTICKS_INSIDE_DOUBLE=\"\x60backticks\x60 work inside double quotes\"", map[string]string{"BACKTICKS_INSIDE_DOUBLE": "\x60backticks\x60 work inside double quotes"}},
		{"backticks are removed", "BACKTICKS=`backticks`", map[string]string{"BACKTICKS": "backticks"}},
		{"spaces are not trimmed in backticks", "BACKTICKS_SPACED=`    backticks    `", map[string]string{"BACKTICKS_SPACED": "    backticks    "}},
		{"double quotes preserved within backticks", "DOUBLE_QUOTES_INSIDE_BACKTICKS=\x60double \"quotes\" work inside backticks\x60", map[string]string{"DOUBLE_QUOTES_INSIDE_BACKTICKS": `double "quotes" work inside backticks`}},
		{"single quotes preserved within backticks", "SINGLE_QUOTES_INSIDE_BACKTICKS=\x60single 'quotes' work inside backticks\x60", map[string]string{"SINGLE_QUOTES_INSIDE_BACKTICKS": `single 'quotes' work inside backticks`}},
		{"double and single quotes are preserved within backticks", "DOUBLE_AND_SINGLE_QUOTES_INSIDE_BACKTICKS=\x60double \"quotes\" and single 'quotes' work inside backticks\x60", map[string]string{"DOUBLE_AND_SINGLE_QUOTES_INSIDE_BACKTICKS": `double "quotes" and single 'quotes' work inside backticks`}},
		{"expand newlines in double quotes", `EXPAND_NEWLINES="expand\nnew\nlines"`, map[string]string{"EXPAND_NEWLINES": `expand
new
lines`}},
		{"don't expand unquoted newlines", `DONT_EXPAND_UNQUOTED=dontexpand\nnewlines`, map[string]string{"DONT_EXPAND_UNQUOTED": `dontexpand\nnewlines`}},
		{"don't expand single quoted newlines", `DONT_EXPAND_SQUOTED='dontexpand\nnewlines'`, map[string]string{"DONT_EXPAND_SQUOTED": `dontexpand\nnewlines`}},
		{"comments", `# COMMENTS=work
INLINE_COMMENTS=inline comments # work #very #well`, map[string]string{"INLINE_COMMENTS": "inline comments"}},
		{"inline comments outside of single quotes work", `INLINE_COMMENTS_SINGLE_QUOTES='inline comments outside of #singlequotes' # work`, map[string]string{"INLINE_COMMENTS_SINGLE_QUOTES": "inline comments outside of #singlequotes"}},
		{"inline comments outside of double quotes work", `INLINE_COMMENTS_DOUBLE_QUOTES="inline comments outside of #doublequotes" # work`, map[string]string{"INLINE_COMMENTS_DOUBLE_QUOTES": "inline comments outside of #doublequotes"}},
		{"inline comments outside of backticks work", "INLINE_COMMENTS_BACKTICKS=`inline comments outside of #backticks` # work", map[string]string{"INLINE_COMMENTS_BACKTICKS": "inline comments outside of #backticks"}},
		{"inline comments outside of quotes start as soon as # sign is seen", `INLINE_COMMENTS_SPACE=inline comments start with a#number sign. no space required.`, map[string]string{"INLINE_COMMENTS_SPACE": "inline comments start with a"}},
		{"equals signs are preserved", `EQUAL_SIGNS=equals==`, map[string]string{"EQUAL_SIGNS": "equals=="}},
		{"inner quotes are preserved", `RETAIN_INNER_QUOTES={"foo": "bar"}`, map[string]string{"RETAIN_INNER_QUOTES": `{"foo": "bar"}`}},
		{"inner quotes are preserved in single quotes", `RETAIN_INNER_QUOTES_AS_STRING='{"foo": "bar"}'`, map[string]string{"RETAIN_INNER_QUOTES_AS_STRING": `{"foo": "bar"}`}},
		{"inner quotes are preserved in backticks", "RETAIN_INNER_QUOTES_AS_BACKTICKS=`{\"foo\": \"bar's\"}`", map[string]string{"RETAIN_INNER_QUOTES_AS_BACKTICKS": `{"foo": "bar's"}`}},
		{"trim space from unquoted string", `TRIM_SPACE_FROM_UNQUOTED=    some spaced out string`, map[string]string{"TRIM_SPACE_FROM_UNQUOTED": "some spaced out string"}},
		{"username", `USERNAME=therealnerdybeast@example.tld`, map[string]string{"USERNAME": `therealnerdybeast@example.tld`}},
		{"spaced key", `    SPACED_KEY = parsed`, map[string]string{"SPACED_KEY": `parsed`}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := EnvFile(tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
