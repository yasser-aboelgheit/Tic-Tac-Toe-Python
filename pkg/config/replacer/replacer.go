package replacer

import (
	"strings"
)

type Replacer struct {
	prefix             string
	keyDelimiter       string
	secondaryDelimiter string
}

func NewReplacer(prefix, keyDelimiter, secondaryDelimiter string) *Replacer {
	return &Replacer{
		prefix:             prefix,
		keyDelimiter:       keyDelimiter,
		secondaryDelimiter: secondaryDelimiter,
	}
}

func (r *Replacer) Replace(s string) string {
	return r.prefix + r.keyDelimiter + strings.ToUpper(
		strings.ReplaceAll(
			s,
			r.secondaryDelimiter,
			r.keyDelimiter,
		),
	)
}
