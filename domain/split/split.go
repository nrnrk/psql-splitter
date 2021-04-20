package split

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Split(fileName string, splitBy int) error {
	r, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	contC := make(chan SplittedStatements, 1000)
	errC := make(chan error)

	splitter := NewSplitter(r, splitBy)
	splitter.Split(contC, errC)

	StartWriting(defaultPrefix(fileName), contC, errC)
	log.Debug("Done")
	return nil
}
