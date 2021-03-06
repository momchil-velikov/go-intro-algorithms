package io

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	input := `digraph a {
	a -> "b 1"; c->d[color="green"];
	e; f->g[x = y]
	a -> g;
	"b 1" ->
	e; d -> a
	a[label=start, color="bright yellow"]
	S
	e -> d[label=α]
	 }`
	g := new(graph)
	err := Read("x.dot", strings.NewReader(input), g)
	if err != nil {
		t.Error(err)
	}
}
