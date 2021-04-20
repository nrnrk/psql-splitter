package split

import (
	"io"
	"os"
	"path"
	"testing"
)

func BenchmarkSplit1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := readerFromRelative(`../../test/unit/20000_statements.sql`)
		s := NewSplitter(r, 1)
		s.Split(make(chan<- SplittedStatements, 20000), make(chan<- error, 1))
	}
}

func BenchmarkSplit10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := readerFromRelative(`../../test/unit/20000_statements.sql`)
		s := NewSplitter(r, 10)
		s.Split(make(chan<- SplittedStatements, 20000), make(chan<- error, 1))
	}
}

func BenchmarkSplit100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := readerFromRelative(`../../test/unit/20000_statements.sql`)
		s := NewSplitter(r, 100)
		s.Split(make(chan<- SplittedStatements, 20000), make(chan<- error, 1))
	}
}

func BenchmarkSplit1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := readerFromRelative(`../../test/unit/20000_statements.sql`)
		s := NewSplitter(r, 1000)
		s.Split(make(chan<- SplittedStatements, 20000), make(chan<- error, 1))
	}
}

func BenchmarkSplit10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := readerFromRelative(`../../test/unit/20000_statements.sql`)
		s := NewSplitter(r, 10000)
		s.Split(make(chan<- SplittedStatements, 20000), make(chan<- error, 1))
	}
}

func readerFromRelative(filename string) io.Reader {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	r, err := os.Open(path.Join(dir, filename))
	if err != nil {
		panic(err)
	}
	return r
}
