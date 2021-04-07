package split

import (
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/nrnrk/psql-splitter/config"
	"github.com/nrnrk/psql-splitter/domain/split/order"
)

func Write(
	prefix string,
	contC <-chan SplittedStatements,
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
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error catched when writing")
			return
		}
	}
}

func write(prefix string, statements string, index int) {
	// TODO: should be set by option
	output := fmt.Sprintf("%s-%s.sql", prefix, order.ByAlphabet(index))
	file := path.Join(config.OutputDir, output)
	log.WithFields(log.Fields{
		"file": file,
	}).Debug("Writing file")
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
	fn := strings.Split(path.Base(fileName), `.`)
	if len(fn) == 0 {
		panic(`this cannot happen`)
	}

	return fn[0]
}
