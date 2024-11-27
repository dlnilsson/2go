package togo

import (
	"strings"
	"unicode"
)

// See:
// https://github.com/dominikh/go-tools/blob/915b568982be0ad65a98e822471748b328240ed0/config/example.conf#L2-L8
// https://google.github.io/styleguide/go/decisions.html#initialisms
var il = []string{"ACL", "API", "ASCII", "CPU", "CSS", "DNS",
	"EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID",
	"IP", "JSON", "QPS", "RAM", "RPC", "SLA",
	"SMTP", "SQL", "SSH", "TCP", "TLS", "TTL",
	"UDP", "UI", "GID", "UID", "UUID", "URI",
	"URL", "UTF8", "VM", "XML", "XMPP", "XSRF",
	"XSS", "SIP", "RTP", "AMQP", "DB", "TS",
	"XMLAPI", "IOS", "GRPC", "DDoS", "Txn",
}

func initialisms() map[string]bool {
	initialisms := make(map[string]bool, len(il))
	for _, word := range il {
		initialisms[word] = true
	}
	return initialisms
}

// lintName based on il
// See: https://github.com/dominikh/go-tools/blob/915b568982be0ad65a98e822471748b328240ed0/stylecheck/st1003/st1003.go#L228-L291
func lintName(name string, initialisms map[string]bool) (should string) {
	// A large part of this function is copied from
	// github.com/golang/lint, Copyright (c) 2013 The Go Authors,
	// licensed under the BSD 3-clause license.

	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	if strings.IndexFunc(name, func(r rune) bool { return !unicode.IsLower(r) }) == -1 {
		return name
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' && i+1 != len(runes)-1 {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); initialisms[u] {
			// Keep consistent case, which is lowercase only at the start.
			if w == 0 && unicode.IsLower(runes[w]) {
				u = strings.ToLower(u)
			}
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			// TODO(dh): this won't be true once we allow custom initialisms
			copy(runes[w:], []rune(u))
		} else if w > 0 && strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}
