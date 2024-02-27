package parse

import (
	"regexp"
	"strings"
)

// Parsing logic is based on: https://github.com/motdotla/dotenv/blob/fe58fef9e2aaf657b6371f7a2f2de8862042878e/lib/main.js

// This is a long and complicated regex. Here is the breakdown of the original:
// - (?:^|^) - match the start of a line
// - \s* - match any number of spaces
// - (?:export\s+)? - match the sequence 'export' followed by one or more spaces, zero or one times
// - ([\w.-]+) - match one or more 'wordy' characters (A-Z, a-z, 0-9, or some punctuation) or '.' or '-'
// - (?:\s*=\s*?|:\s+?) - match either: any number of spaces, followed by '=', followed by any number of spaces, OR a colon followed by one or more spaces
// - (\s*'(?:\\'|[^'])*'|\s*"(?:\\"|[^"])*"|\s*`(?:\\`|[^`])*`|[^#\r\n]+)? - zero or one times, match any of the following:
//   - any number of spaces, followed by a single quote, followed by any number of escaped single-quote or non-quote characters, followed by a single quote, OR
//   - as with the previous, but using double quotes instead, OR
//   - as with the previous, but using backticks instead, OR
//   - one or more characters that are not a '#' or a new line
//
// - \s* - match any number of spaces
// - (?:#.*)? - zero or one times, match a '#' followed by any number of characters
// - (?:$|$) - match the end of a line.
// - mg - match globally, across multiple lines.
var envFileRegex = regexp.MustCompile("(?m:^\\s*(?:export\\s+)?([\\w.-]+)(?:\\s*=\\s*?|:\\s+?)(\\s*'(?:\\\\'|[^'])*'|\\s*\"(?:\\\\\"|[^\"])*\"|\\s*\\x60(?:\\\\\\x60|[^\\x60])*\\x60|[^#\\r\\n]+)?\\s*(?:#.*)?$)")

var quoteReplaceRegex = regexp.MustCompile("^(['\"`])([\\s\\S]*)(['\"`])$")

func EnvFile(bs string) (map[string]string, error) {
	matches := envFileRegex.FindAllStringSubmatch(bs, -1)
	ret := make(map[string]string, len(matches))
	for _, submatch := range matches {
		var key, value string
		if len(submatch) >= 2 {
			key = submatch[1]
		}
		// Value defaults to "" if not present
		if len(submatch) >= 3 {
			value = strings.TrimSpace(submatch[2])
			if len(value) > 0 {
				firstValueChar := value[0]
				// Strip quotes from beginning and end of string
				value = quoteReplaceRegex.ReplaceAllString(value, "$2")
				if firstValueChar == '"' {
					value = strings.ReplaceAll(value, "\\n", "\n")
					value = strings.ReplaceAll(value, "\\r", "\r")
				}
			}
		}

		ret[key] = value
	}

	return ret, nil
}
