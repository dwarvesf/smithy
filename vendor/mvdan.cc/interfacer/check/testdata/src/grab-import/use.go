package foo

import (
	"grab-import/def"
	def2 "grab-import/def/nested"
)

type St struct{}

func (s *St) Foo(rc def.ReadCloser, i int) int {
	rc.Close()
	return def.SomeVar
}

func NonInterestingCall() {
	def.SomeFunc()
}

func Foo(f def2.Fooer) {
	f.Foo()
}

type st2 struct{}

func (s st2) Foo()

func FooWrong(s st2) { // WARN s can be grab-import/def/nested.Fooer
	s.Foo()
}
