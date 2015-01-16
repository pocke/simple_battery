package main

import (
	"testing"

	"github.com/pocke/goevent"
)

func TestNewBatteries(t *testing.T) {
	e := goevent.NewTable()
	_ = NewBatteries(e)
}

func TestList(t *testing.T) {
	e := goevent.NewTable()
	b := NewBatteries(e)

	l := b.list()
	t.Log(l)
}

func TestGet(t *testing.T) {
	e := goevent.NewTable()
	b := NewBatteries(e)

	i, err := b.get(0)
	t.Log(i)
	if err != nil {
		t.Error(err)
	}

	i, err = b.get(1)
	t.Log(i)
	if err != nil {
		t.Error(err)
	}

	i, err = b.get(100)
	if err == nil {
		t.Errorf("expected: error, but got: nil")
	}
}

func TestInclude(t *testing.T) {
	l := []int{5, 4, 3, 2, 1}
	i := include(l, 3)
	if i != 2 {
		t.Errorf("expected: 2, but got %d", i)
	}

	i = include(l, 6)
	if i != -1 {
		t.Errorf("expected: -1, but got %d", i)
	}
}

func TestSliceDiff(t *testing.T) {
	a := []int{38, 4, 11, 24, 9}
	b := []int{22, 4, 9}

	s := sliceDiff(a, b)
	if !(include(s, 38) != -1 && include(s, 11) != -1 && include(s, 24) != -1 && include(s, 4) == -1 && include(s, 9) == -1) {
		t.Errorf("expected: %v, but got: %v", []int{38, 11, 24}, s)
	}

	s = sliceDiff(b, a)
	if !(include(s, 22) != -1 && include(s, 4) == -1 && include(s, 9) == -1) {
		t.Errorf("expected: %v, but got: %v", []int{22}, s)
	}
}
