package freepool

import "testing"

func TestName(t *testing.T) {
	fp := New(16)
	for i := 0; i < 2; i++ {
		page := fp.Page()
		s1 := page.String("hello world")
		s2 := page.String("abcdefg")
		s3 := page.String("HELLO WORLD")
		s4 := page.String("ABCDEFG")
		page.String("hello world")
		page.String("abcdefg")
		page.String("HELLO WORLD")
		page.String("ABCDEFG")
		t.Logf("s1: %#v, s2: %#v, s3: %#v, s4: %#v", s1.RefString(), s2.RefString(), s3.RefString(), s4.RefString())
		fpi := fp.(*freePool)
		t.Logf("%#v", fpi)
		t.Logf("%#v", fpi.chunks[0])
		t.Logf("%#v", fpi.chunks[0].page(0, fpi.wrappedPageSize))
		page.Release()
		t.Logf("%#v", fpi)
		t.Logf("%#v", fpi.chunks[0])
	}
}
