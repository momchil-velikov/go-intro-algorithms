package io

import (
	"fmt"
	"io"
)

type GraphIO interface {
	NodeCount() int
	NodeAttrs(id int) []string
	EdgeCount(id int) int
	EdgeAttrs(id int, i int) (dst int, attrs []string)
}

func Write(w io.Writer, g GraphIO) error {
	if _, e := w.Write([]byte("digraph G {\n")); e != nil {
		return e
	}
	for i := 0; i < g.NodeCount(); i++ {
		if e := writeNode(w, i, g.NodeAttrs(i)); e != nil {
			return e
		}
		for j := 0; j < g.EdgeCount(i); j++ {
			d, attrs := g.EdgeAttrs(i, j)
			if e := writeEdge(w, i, d, attrs); e != nil {
				return e
			}
		}
	}
	_, e := w.Write([]byte("}\n"))
	return e
}

func writeNode(w io.Writer, id int, attrs []string) error {
	s := fmt.Sprintf("n%d", id)
	if _, e := w.Write([]byte(s)); e != nil {
		return e
	}
	return writeAttrs(w, attrs)
}

func writeEdge(w io.Writer, src int, dst int, attrs []string) error {
	s := fmt.Sprintf("n%d -> n%d", src, dst)
	if _, e := w.Write([]byte(s)); e != nil {
		return e
	}
	return writeAttrs(w, attrs)
}

func writeAttrs(w io.Writer, attrs []string) error {
	if len(attrs) > 1 {
		if _, e := w.Write([]byte{'['}); e != nil {
			return e
		}
		if e := writeAttrList(w, attrs); e != nil {
			return e
		}
		if _, e := w.Write([]byte{']', '\n'}); e != nil {
			return e
		}
		return nil
	} else {
		_, e := w.Write([]byte{'\n'})
		return e
	}
}

func writeAttrList(w io.Writer, attrs []string) error {
	n := len(attrs)
	for i := 0; i < n; i += 2 {
		s := fmt.Sprintf("\"%s\" = \"%s\", ", attrs[i], attrs[i+1])
		if _, e := w.Write([]byte(s)); e != nil {
			return e
		}
	}
	return nil
}
