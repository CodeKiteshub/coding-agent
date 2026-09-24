package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	flag.Parse()

	if *verbose {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Verbose logging enabled")
	} else {
		log.SetOutput(os.Stdout)
		log.SetFlags(0)
		log.SetPrefix("")
	}

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	agent := NewAgent(getUserMessage, *verbose)
	err := agent.Run(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
