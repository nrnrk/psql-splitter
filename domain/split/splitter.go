package split

import (
	"bufio"
	"io"

	pg_query "github.com/lfittl/pg_query_go"
	log "github.com/sirupsen/logrus"

	"github.com/nrnrk/psql-splitter/domain/split/head"
)

type Splitter interface {
	Split(contC chan<- SplittedStatements, errC chan<- error)
	AppendSql()
	CanSplit() bool
	FlushSql()
	FlushStmts()
	IsEndStmt() bool
	IsContentEmpty() bool
}

type splitter struct {
	r        io.Reader
	splitNum int
	splitCnt int
	Cont     *SplittedStatements
	sql      []byte
	buf      []byte
	head     head.Head
}

func NewSplitter(r io.Reader, splitNum int) *splitter {
	return &splitter{
		r:        bufio.NewReader(r),
		splitNum: splitNum,
		splitCnt: 0,
		Cont: &SplittedStatements{
			Statements: ``,
			Order:      0,
		},
		sql:  make([]byte, 0, 500*1024),
		buf:  make([]byte, 100*1024),
		head: head.NewHead(),
	}
}

func (s *splitter) Split(contC chan<- SplittedStatements, errC chan<- error) {
	for {
		n, err := s.r.Read(s.buf)
		if n == 0 {
			if err != nil && err != io.EOF {
				errC <- err
				panic(err)
			}
			break
		}

		for i := 0; i < n; i++ {
			s.sql = append(s.sql, s.buf[i])
			s.head.Continue(s.buf[i])

			if s.IsEndStmt() {
				s.AppendSql()
				if s.CanSplit() {
					contC <- *s.Cont
					s.FlushStmts()
				}
				s.FlushSql()
			}
		}
	}

	if !s.IsContentEmpty() {
		contC <- *s.Cont
		s.FlushStmts()
	}
	close(contC)
}

func (s *splitter) AppendSql() {
	s.splitCnt++
	s.Cont.Statements += string(s.sql)
}

func (s *splitter) CanSplit() bool {
	if s.splitCnt < s.splitNum {
		return false
	}

	if err := s.parseCheck(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Parse failed (Please check all SQLs are valid)")
		panic(err)
	}

	return true
}

func (s *splitter) IsEndStmt() bool {
	return s.head.IsEndStmt()
}

func (s *splitter) FlushSql() {
	s.head.Restart()
	// keep capacity
	s.sql = s.sql[:0]
}

func (s *splitter) FlushStmts() {
	s.splitCnt = 0
	s.Cont.Order++
	s.Cont.Statements = ``
}

func (s *splitter) IsContentEmpty() bool {
	return s.Cont.Statements == ``
}

func (s *splitter) parseCheck() error {
	if _, err := pg_query.Parse(string(s.sql)); err != nil {
		return err
	}
	return nil
}
