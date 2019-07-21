package errorHandling

import (
	"testing"
)

type record struct {
	U string
	V int
}

func doSomething0() (r *record, e error) {
	v, e := foo()
	var u string
	fe := []func(){
		func() { u, e = bar() },
		func() { r = &record{V: v, U: u} },
	}
	n := 0
	for e == nil && n != len(fe) {
		fe[n]()
		n = n + 1
	}
	return
}

func doSomething1() (r *record, e error) {
	v, e := foo()
	var u string
	fs := []func(){
		func() { u, e = bar() },
		func() { r = &record{V: v, U: u} },
	}
	bLnSearch(
		func(i int) bool { b := e != nil; fs[i](); return b },
		len(fs),
	)
	return
}

func bLnSearch(bf func(int) bool, n int) (ok bool, i int) {
	i = 0
	for i != n && !bf(i) {
		i = i + 1
	}
	ok = i != n
	return
}

func doSomething() (r *record, e error) {
	var v int
	var u string
	fs := []func(){
		func() { v, e = foo() },
		func() { u, e = bar() },
		func() { r = &record{V: v, U: u} },
	}
	bLnSearch(
		func(i int) bool { fs[i](); return e != nil },
		len(fs),
	)
	return
}

func foo() (v int, e error) {
	return
}

func bar() (u string, e error) {
	return
}

func TestDoSomething(t *testing.T) {
	r, e := doSomething()
	if e != nil || r.U != "" || r.V != 0 {
		t.Fatalf("Failed: %v %v", e, r)
	}
}
