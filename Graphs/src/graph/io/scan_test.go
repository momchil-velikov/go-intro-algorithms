package io

import (
	"strings"
	"testing"
)

func TestScanEmpty(t *testing.T) {
	s := newScanner(strings.NewReader(""))
	tok := s.getToken()
	if tok.kind != tEOF {
		t.Error("got non EOF token from empty input")
	}
}

func TestScanPunctuators(t *testing.T) {
	input := "[ ] { } , : ; = ->"
	exp := []kindT{tLBracket, tRBracket, tLBrace, tRBrace, tComma, tColon, tSemi, tEQ, tArrow, tEOF}
	i := 0
	s := newScanner(strings.NewReader(input))
	tok := tokenT{}
	for tok = s.getToken(); i < len(exp); tok = s.getToken() {
		if tok.kind != exp[i] {
			t.Errorf("unexpected token %v, should be %v\n", tok.kind, exp[i])
		}
		i++
	}
}

func TestScanSymbols(t *testing.T) {
	input := `a αβγ
     12γ ->

      "συμ with\" spaces"`
	exp := []kindT{tSym, tSym, tSemi, tSym, tArrow, tSym, tEOF}
	val := []string{"a", "αβγ", "", "12γ", "", "συμ with\" spaces", ""}
	i := 0
	s := newScanner(strings.NewReader(input))
	tok := tokenT{}
	for tok = s.getToken(); i < len(exp); tok = s.getToken() {
		if tok.kind != exp[i] {
			t.Errorf("unexpected token %v, should be %v\n", tok.kind, exp[i])
		}
		if tok.kind == tSym {
			v, ok := tok.val.(string)
			if !ok {
				t.Errorf("invalid type of the semantic value of a %v token", tok.kind)
			}
			if v != val[i] {
				t.Errorf("invalid semantic value of a %v token", tok.kind)
			}
		}
		i++
	}
}
