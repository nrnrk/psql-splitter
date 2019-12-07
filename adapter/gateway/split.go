package gateway

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/nrnrk/psql-splitter/domain/split"
	"github.com/nrnrk/psql-splitter/domain/split/order"
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
		writer(defaultPrefix(fileName), contC, terminateC, errC)
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

func writer(
	prefix string,
	contC <-chan split.SplittedStatements,
	terminateC <-chan bool,
	errC <-chan error,
) {
	for {
		select {
		case c := <-contC:
			write(prefix, c.Statements, c.Order)
		case t := <-terminateC:
			if t {
				return
			}
		case err := <-errC:
			fmt.Printf("err: %v", err)
			return
		}
	}
}

func write(prefix string, statements string, index int) {
	// TODO: should be set by option
	output := fmt.Sprintf("%s-%s.sql", prefix, order.ByAlphabet(index))
	fmt.Printf("Writing %s\n", output)
	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(statements))
	if err != nil {
		panic(err)
	}
}

func defaultPrefix(fileName string) string {
	fn := strings.Split(fileName, `.`)
	if len(fn) == 0 {
		panic(`this cannot happen`)
	}

	return fn[0]
}
