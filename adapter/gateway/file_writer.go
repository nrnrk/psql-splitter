package gateway

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/nrnrk/psql-splitter/config"

	"github.com/nrnrk/psql-splitter/domain/split"
	"github.com/nrnrk/psql-splitter/domain/split/order"
)

func Write(
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
	file := path.Join(config.OutputDir, output)
	fmt.Printf("Writing %s\n", file)
	f, err := os.Create(file)
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
