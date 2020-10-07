package main

import (
	"context"
	"log"

	"github.com/kk-no/YaneuraGo/protocol/usi"
)

func main() {
	engine := usi.New()
	if err := engine.Connect(context.TODO(), "/engine"); err != nil {
		log.Fatal(err)
	}
}
