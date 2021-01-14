package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/kk-no/YaneuraGo/usi"
)

func main() {
	ctx := context.Background()

	usiEngine := usi.NewEngine()
	if err := usiEngine.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer usiEngine.Disconnect(ctx)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			line := scanner.Text()
			usiEngine.SendCommand(ctx, line)
			if line == "quit" {
				break
			}
		}
	}
}
