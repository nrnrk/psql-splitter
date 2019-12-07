package presenter

import (
	"fmt"
	"os"
	"strings"

	"github.com/nrnrk/psql-splitter/domain/splitter/order"
)

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
