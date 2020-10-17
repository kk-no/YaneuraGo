package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/kk-no/YaneuraGo/protocol/state/engine"
	"github.com/kk-no/YaneuraGo/protocol/usi"
)

func main() {
	ctx := context.Background()

	usiEngine := usi.New()
	if err := usiEngine.Connect(ctx, filepath.Join(engine.Dir, engine.Path)); err != nil {
		log.Fatal(err)
	}
	defer usiEngine.Disconnect(ctx)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			if err := usiEngine.SendCommand(ctx, scanner.Text()); err != nil {
				log.Fatal(err)
			}
		}
	}
}
