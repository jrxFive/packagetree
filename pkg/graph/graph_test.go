package graph

import (
	"testing"
)

func newGraph(t *testing.T) *TreeGraph {
	g, err := NewGraph()

	if err != nil {
		t.Error("Failed to create new graph")
	}

	return g
}

func TestExists(t *testing.T) {

	g := newGraph(t)

	g.Add("boo")
	if exists := g.Exists("boo"); !exists {
		t.Errorf("Node does not exist after adding to graph")
	}

	if exists := g.Exists("bar"); exists {
		t.Errorf("Node should not exist in graph")
	}
}

func TestAdd(t *testing.T) {

	g := newGraph(t)
	var errResult bool

	var tests = []struct {
		name           string
		edges          []string
		errShouldBeNil bool
	}{
		{name: "boo", errShouldBeNil: true},
		{name: "boo", edges: []string{"foo,bar"}, errShouldBeNil: false},
		{name: "foo", errShouldBeNil: true},
		{name: "bar", errShouldBeNil: true},
		{name: "boo", edges: []string{"foo"}, errShouldBeNil: true},
		{name: "boo", edges: []string{"foo", "bar"}, errShouldBeNil: true},
	}

	for _, test := range tests {
		err := g.Add(test.name, test.edges...)

		if err == nil {
			errResult = true
		} else {
			errResult = false
		}

		if errResult != test.errShouldBeNil {
			t.Errorf("Add(%q, %q) = %v", test.name, test.edges, err)
		}

	}

}

func TestRemoveEmpty(t *testing.T) {
	g := newGraph(t)
	var errResult bool

	var tests = []struct {
		name        string
		shouldBeNil bool
	}{
		{name: "boo", shouldBeNil: true},
		{name: "foo", shouldBeNil: true},
		{name: "bar", shouldBeNil: true},
	}

	for _, test := range tests {
		err := g.Remove(test.name)

		if err == nil {
			errResult = true
		} else {
			errResult = false
		}

		if errResult != test.shouldBeNil {
			t.Errorf("Remove(%q) = %v", test.name, err)
		}

	}
}

func TestRemove(t *testing.T) {
	g := newGraph(t)

	var errResult bool
	var addPackages = []struct {
		name  string
		edges []string
	}{
		{name: "boo"},
		{name: "foo"},
		{name: "bar"},
		{name: "bar2"},
		{name: "zap", edges: []string{"bar", "bar2"}},
	}

	var removePackages = []struct {
		name           string
		errShouldBeNil bool
	}{
		{name: "boo", errShouldBeNil: true},
		{name: "foo", errShouldBeNil: true},
		{name: "zap", errShouldBeNil: false},
		{name: "bar", errShouldBeNil: true},
		{name: "bar2", errShouldBeNil: true},
		{name: "zap", errShouldBeNil: true},
	}

	for _, p := range addPackages {
		g.Add(p.name, p.edges...)
	}

	for _, test := range removePackages {
		err := g.Remove(test.name)

		if err == nil {
			errResult = true
		} else {
			errResult = false
		}

		if errResult != test.errShouldBeNil {
			t.Errorf("Remove(%q) = %v", test.name, err)
		}
	}
}
