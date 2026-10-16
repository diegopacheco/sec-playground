package taint

import (
	"hash/fnv"
	"strings"
	"unicode"
)

const shingleSize = 4

type token struct {
	text  string
	start int
	end   int
}

func tokenize(s string) []token {
	var out []token
	start := -1
	for i, r := range s {
		word := unicode.IsLetter(r) || unicode.IsDigit(r)
		if word && start < 0 {
			start = i
		}
		if !word && start >= 0 {
			out = append(out, token{strings.ToLower(s[start:i]), start, i})
			start = -1
		}
	}
	if start >= 0 {
		out = append(out, token{strings.ToLower(s[start:]), start, len(s)})
	}
	return out
}

func shingles(toks []token) []uint64 {
	if len(toks) < shingleSize {
		return nil
	}
	out := make([]uint64, 0, len(toks)-shingleSize+1)
	for i := 0; i+shingleSize <= len(toks); i++ {
		h := fnv.New64a()
		for _, t := range toks[i : i+shingleSize] {
			h.Write([]byte(t.text))
			h.Write([]byte{0})
		}
		out = append(out, h.Sum64())
	}
	return out
}
