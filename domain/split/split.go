package split

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

func Split(fileName string, splitBy int) error {
	r, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	contC := make(chan SplittedStatements)
	terminateC := make(chan bool, 1)
	errC := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Write(defaultPrefix(fileName), contC, terminateC, errC)
	}()

	splitter := NewSplitter(r, splitBy)

	splitter.Split(contC, errC)

	terminateC <- true
	wg.Wait()
	log.Debug("Read done")
	return nil
}
