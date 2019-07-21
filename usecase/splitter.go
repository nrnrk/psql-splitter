package usecase

import (
	"fmt"
	"io"

	pg_query "github.com/lfittl/pg_query_go"
	"github.com/nrnrk/psql-splitter/usecase/state"
)

type splitter struct {
	splitNum int
	splitCnt int
	cont     *writeContent
	sql      []byte
	buf      []byte
	head     state.Head
}

func newSplitter(splitNum int) *splitter {
	return &splitter{
		splitNum: splitNum,
		splitCnt: 0,
		cont: &writeContent{
			statements: ``,
			order:      0,
		},
		sql:  make([]byte, 0, 50),
		buf:  make([]byte, 1),
		head: state.NewHead(),
	}
}

func (s *splitter) readFrom(r io.Reader) error {
	_, err := r.Read(s.buf)
	if err != nil {
		return err
	}

	s.sql = append(s.sql, s.buf[0])
	s.head.Continue(s.buf[0])

	return nil
}

func (s *splitter) appendSql() {
	s.splitCnt++
	s.cont.statements += string(s.sql)
}

func (s *splitter) canSplit() bool {
	if s.splitCnt < s.splitNum {
		return false
	}

	if err := s.parseCheck(); err != nil {
		fmt.Printf("err: %v", err)
		panic(err)
	}

	return true
}

func (s *splitter) isEndStmt() bool {
	return s.head.IsEndStmt()
}

func (s *splitter) flushSql() {
	s.head.Restart()
	// keep capacity
	s.sql = s.sql[:0]
}

func (s *splitter) flushStmts() {
	s.splitCnt = 0
	s.cont.order++
	s.cont.statements = ``
}

func (s *splitter) isContentEmpty() bool {
	return s.cont.statements == ``
}

func (s *splitter) parseCheck() error {
	if _, err := pg_query.Parse(string(s.sql)); err != nil {
		return err
	}
	return nil
}
