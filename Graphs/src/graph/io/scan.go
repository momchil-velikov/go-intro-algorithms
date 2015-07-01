package io

import (
	"io"
	"unicode"
)

type kindT int

//go:generate stringer -type=kindT
const (
	tEOF kindT = iota
	tSym
	tArrow
	tLBracket
	tRBracket
	tLBrace
	tRBrace
	tEQ
	tComma
	tColon
	tSemi
	tErr
)

type tokenT struct {
	kind      kindT
	val       interface{}
	line, col uint
}

type scanner struct {
	rd          io.RuneReader
	tLine, tCol uint
	line, col   uint
	needSemi    bool
	ch          rune
	err         error
}

func newScanner(r io.RuneReader) scanner {
	s := scanner{rd: r, line: 1, col: 0}
	s.next()
	return s
}

func (s *scanner) startToken() {
	s.tLine, s.tCol = s.line, s.col
}

func (s *scanner) newToken(kind kindT, val interface{}) tokenT {
	return tokenT{kind: kind, val: val, line: s.tLine, col: s.tCol}
}

func (s *scanner) next() {
	s.ch, _, s.err = s.rd.ReadRune()
	if s.err == nil {
		s.col++
	}
}

func (s *scanner) errorToken() tokenT {
	s.startToken()
	if s.err == io.EOF {
		return s.newToken(tEOF, nil)
	} else {
		return s.newToken(tErr, s.err)
	}
}

func (s *scanner) getToken() tokenT {
	// We may have an error from the last call to next()
	if s.err != nil {
		return s.errorToken()
	}
	// Skip whitespace; insert semicolons as needed
	for unicode.IsSpace(s.ch) {
		// Do not advance the line number until the newline is actually consumed.
		if s.ch == '\n' {
			if s.needSemi {
				s.needSemi = false
				s.startToken()
				return s.newToken(tSemi, nil)
			} else {
				s.line++
				s.col = 1
			}
		}
		s.next()
		if s.err != nil {
			return s.errorToken()
		}
	}
	// Scan token
	s.startToken()
	s.needSemi = false
	switch s.ch {
	case '[':
		s.next()
		return s.newToken(tLBracket, nil)
	case ']':
		s.needSemi = true
		s.next()
		return s.newToken(tRBracket, nil)
	case '{':
		s.next()
		return s.newToken(tLBrace, nil)
	case '}':
		s.next()
		return s.newToken(tRBrace, nil)
	case '=':
		s.next()
		return s.newToken(tEQ, nil)
	case ',':
		s.next()
		return s.newToken(tComma, nil)
	case ':':
		s.next()
		return s.newToken(tColon, nil)
	case ';':
		s.next()
		return s.newToken(tSemi, nil)
	case '-':
		s.next()
		if s.err != nil {
			return s.errorToken()
		}
		if s.ch == '>' {
			s.next()
			return s.newToken(tArrow, nil)
		}
	case '"':
		s.needSemi = true
		return s.scanQuotedSymbol()
	}
	if s.ch == '_' || unicode.IsLetter(s.ch) || unicode.IsDigit(s.ch) {
		s.needSemi = true
		return s.scanSymbol()
	}
	return s.newToken(tErr, "invalid character")
}

func (s *scanner) scanQuotedSymbol() tokenT {
	sym := []rune(nil)
	for {
		s.next()
		if s.err != nil {
			return s.errorToken()
		}
		switch {
		case s.ch == '\n':
			return s.newToken(tErr, "newline inside quotes")
		case s.ch == '"':
			s.next()
			return s.newToken(tSym, string(sym))
		case s.ch == '\\':
			s.next()
			if s.err != nil {
				return s.errorToken()
			}
			fallthrough
		default:
			sym = append(sym, s.ch)
		}
	}
}

func (s *scanner) scanSymbol() tokenT {
	sym := []rune(nil)
	for s.ch == '_' || unicode.IsLetter(s.ch) || unicode.IsDigit(s.ch) {
		sym = append(sym, s.ch)
		s.next()
		if s.err != nil {
			return s.errorToken()
		}
	}
	return s.newToken(tSym, string(sym))
}
