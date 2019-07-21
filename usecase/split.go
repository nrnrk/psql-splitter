package usecase

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/nrnrk/psql-splitter/usecase/ordering"
)

func Split(fileName string, splitBy int) error {
	r, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	contC := make(chan writeContent)
	terminateC := make(chan bool, 1)
	errC := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		writer(defaultPrefix(fileName), contC, terminateC, errC)
	}()

	splitter := newSplitter(splitBy)

	for {
		err := splitter.readFrom(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			errC <- err
			panic(err)
		}

		if splitter.isEndStmt() {
			splitter.appendSql()
			if splitter.canSplit() {
				contC <- *splitter.cont
				splitter.flushStmts()
			}
			splitter.flushSql()
		}
	}

	if !splitter.isContentEmpty() {
		splitter.cont.statements += "\n"
		contC <- *splitter.cont
		splitter.flushStmts()
	}

	terminateC <- true
	wg.Wait()
	fmt.Println("\nread done")
	return nil
}

type writeContent struct {
	statements string
	order      int
}

func writer(
	prefix string,
	contC <-chan writeContent,
	terminateC <-chan bool,
	errC <-chan error,
) {
	for {
		select {
		case c := <-contC:
			write(prefix, c.statements, c.order)
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

func write(prefix string, statements string, order int) {
	// TODO: should be set by option
	output := fmt.Sprintf("%s-%s.sql", prefix, ordering.ByAlphabet(order))
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
