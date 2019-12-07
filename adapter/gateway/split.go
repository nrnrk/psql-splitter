package gateway

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/nrnrk/psql-splitter/domain/split"
)

func Split(fileName string, splitBy int) error {
	r, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	contC := make(chan split.SplittedStatements)
	terminateC := make(chan bool, 1)
	errC := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Write(defaultPrefix(fileName), contC, terminateC, errC)
	}()

	splitter := split.NewSplitter(splitBy)

	for {
		err := splitter.ReadFrom(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			errC <- err
			panic(err)
		}

		if splitter.IsEndStmt() {
			splitter.AppendSql()
			if splitter.CanSplit() {
				contC <- *splitter.Cont
				splitter.FlushStmts()
			}
			splitter.FlushSql()
		}
	}

	if !splitter.IsContentEmpty() {
		splitter.Cont.AddNewLine()
		contC <- *splitter.Cont
		splitter.FlushStmts()
	}

	terminateC <- true
	wg.Wait()
	fmt.Println("\nread done")
	return nil
}
