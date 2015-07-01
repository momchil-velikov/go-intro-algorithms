package io

import (
	"fmt"
	"io"
)

type GraphBuilder interface {
	CreateNode(string) int
	SetNodeAttrs(int, []string)
	CreateEdge(int, int, []string)
}

type parseError struct {
	name      string
	line, col uint
	msg       string
}

func (e parseError) Error() string {
	return fmt.Sprintf("parse error: %s:%d:%d: %s", e.name, e.line, e.col, e.msg)
}

type parser struct {
	graph GraphBuilder
	nmap  map[string]int
	name  string
	scn   scanner
	token tokenT
}

func Read(name string, r io.RuneReader, g GraphBuilder) error {
	p := parser{graph: g, nmap: make(map[string]int), name: name, scn: newScanner(r)}
	return p.parseGraph()
}

func (p *parser) createNode(name string) int {
	if id, ok := p.nmap[name]; !ok {
		id = p.graph.CreateNode(name)
		p.nmap[name] = id
		return id
	} else {
		return id
	}
}

func (p *parser) error(err interface{}) error {
	return parseError{p.name, p.token.line, p.token.col, fmt.Sprint(err)}
}

func (p *parser) next() {
	p.token = p.scn.getToken()
}

func (p *parser) expect(k kindT) error {
	if p.token.kind == k {
		return nil
	} else {
		if p.token.kind == tErr {
			return p.error(p.token.val)
		} else {
			return p.error(fmt.Sprintf("got token %v, expected %v", p.token.kind, k))
		}
	}
}

func (p *parser) match(k kindT) error {
	if err := p.expect(k); err == nil {
		p.next()
		return nil
	} else {
		return err
	}
}

// graph ::= "digraph" sym "{" stmt-list> "}"
func (p *parser) parseGraph() error {
	p.next()
	if err := p.expect(tSym); err != nil {
		return err
	} else if str, ok := p.token.val.(string); !ok || str != "digraph" {
		return p.error("expected 'digraph'")
	}
	p.next()
	if err := p.match(tSym); err != nil {
		return err
	}
	if err := p.match(tLBrace); err != nil {
		return err
	}
	if p.token.kind == tSym {
		if err := p.parseStmtList(); err != nil {
			return err
		}
	}
	if err := p.match(tRBrace); err != nil {
		return err
	}
	if err := p.match(tEOF); err != nil {
		return err
	}
	return nil
}

// stmt-list ::= stmt ";" [ stmt-list ]
// stmt ::= node-stmt | edge-stmt
func (p *parser) parseStmtList() error {
	for p.token.kind == tSym {
		src := p.token.val.(string)
		p.next()
		if p.token.kind == tArrow {
			if e := p.parseEdgeStmt(src); e != nil {
				return e
			}
		} else if e := p.parseNodeStmt(src); e != nil {
			return e
		}
		if e := p.match(tSemi); e != nil {
			return e
		}
	}
	return nil
}

// node-stmt ::= sym [ "[" attr-list [","] "]" ]
func (p *parser) parseNodeStmt(name string) error {
	if p.token.kind == tLBracket {
		if attrs, err := p.parseAttrList(); err != nil {
			return err
		} else {
			id := p.createNode(name)
			p.graph.SetNodeAttrs(id, attrs)
		}
	}
	return nil
}

// edge-stmt ::= sym "->" sym [ "[" attr-list [","] "]"]
func (p *parser) parseEdgeStmt(src string) error {
	p.next()
	if e := p.expect(tSym); e != nil {
		return e
	}
	dst := p.token.val.(string)
	p.next()
	var attrs []string = nil
	if p.token.kind == tLBracket {
		if a, e := p.parseAttrList(); e != nil {
			return e
		} else {
			attrs = a
		}
	}
	s, d := p.createNode(src), p.createNode(dst)
	p.graph.CreateEdge(s, d, attrs)
	return nil
}

// attr-list ::= attr [ "," attr-list ]
// attr ::= sym "=" sym
func (p *parser) parseAttrList() ([]string, error) {
	p.next()
	var attrs []string = nil
	for p.token.kind == tSym {
		name := p.token.val.(string)
		p.next()
		if e := p.match(tEQ); e != nil {
			return nil, e
		}
		if e := p.expect(tSym); e != nil {
			return nil, e
		}
		val := p.token.val.(string)
		attrs = append(attrs, name, val)
		p.next()
		if p.token.kind == tComma {
			p.next()
		}
	}
	e := p.match(tRBracket)
	return attrs, e
}
