package usi

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kk-no/YaneuraGo/protocol/state/engine"
)

type Engine interface {
	Connect(ctx context.Context, path string) error
	Disconnect(ctx context.Context) error
}

type usi struct {
	state    engine.State
	options  map[string]string
	received string
}

func New() Engine {
	return &usi{}
}

func (u usi) Connect(ctx context.Context, path string) error {
	log.Println("call Connect")

	if err := u.Disconnect(ctx); err != nil {
		return err
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}

	e := filepath.Join(root, path)
	fmt.Println(e)

	return nil
}

func (u usi) Disconnect(ctx context.Context) error {
	log.Println("call Disconnect")
	return nil
}
