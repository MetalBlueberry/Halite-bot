package main

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// NewLocalConf uses stdin as soruce and stdout as response
func NewLocalConf() GameConfig {
	done := make(chan struct{})
	stdin := make(chan string)
	stdout := make(chan string)

	go func() {
		defer close(done)
		for {
			message, ok := <-stdout
			if !ok {
				log.Println("stdout closed")
				return
			}
			_, err := fmt.Fprintf(os.Stdout, "%s\n", message)
			if err != nil {
				log.Panic(err)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	go func(scanner *bufio.Scanner) {
		defer close(done)

		for scanner.Scan() {
			stdin <- scanner.Text()
		}

		if scanner.Err() != nil {
			panic(scanner.Err())
		}

	}(scanner)

	return NewConf(stdin, stdout)
}
