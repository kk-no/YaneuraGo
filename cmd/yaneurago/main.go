package main

import (
	"context"
	"log"
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

	// TODO: Receiving Commands from Outside.
	cmd := []string{"usi", "isready", "quit"}
	for _, v := range cmd {
		if err := usiEngine.SendCommand(ctx, v); err != nil {
			log.Fatal(err)
		}
	}
	usiEngine.Disconnect(ctx)
}
