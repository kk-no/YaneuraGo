package usi

import (
	"context"
	"log"
	"strings"

	"github.com/kk-no/YaneuraGo/dir"
	"github.com/kk-no/YaneuraGo/state/engine"
)

type Engine interface {
	SetState(ctx context.Context, state engine.State)
	Connect(ctx context.Context, path string) error
	Disconnect(ctx context.Context) error
	IsConnected(ctx context.Context) bool
	SendCommand(ctx context.Context, command string)
}

type usi struct {
	state engine.State
	// options map[string]string
	process ReadWriteProcessor
	result  ResultManager
	isDebug bool
}

func New() Engine {
	return &usi{
		state: engine.Disconnected,
		// FIXME: change debug setting
		isDebug: true,
	}
}

func (u *usi) SetState(ctx context.Context, state engine.State) {
	if u.isDebug {
		log.Printf("Engine status change %v -> %v\n", u.state, state)
	}
	u.state = state
}

func (u *usi) Connect(ctx context.Context, path string) error {
	u.SetState(ctx, engine.WaitConnecting)

	// Need to move the directory in order for the engine to read the eval.
	_, err := dir.ChangeDir(path)
	if err != nil {
		return err
	}

	u.process, err = NewReadWriteProcessor(ctx)
	if err != nil {
		log.Println("Failed to init process:", err)
	}
	u.process.Start(ctx)

	u.result = NewResultManager()

	u.SetState(ctx, engine.Connected)

	return nil
}

func (u *usi) Disconnect(ctx context.Context) error {
	if err := u.process.Stop(); err != nil {
		return err
	}
	u.process = nil
	u.SetState(ctx, engine.Disconnected)
	return nil
}

func (u *usi) IsConnected(ctx context.Context) bool {
	return u.process != nil
}

func (u *usi) SendCommand(ctx context.Context, command string) {
	u.process.SendCommand(ctx, command)
}

func (u *usi) HandleMessage(ctx context.Context, message string) {
	u.result.ReceiveMessage(message)

	var token string
	if index := strings.Index(message, " "); index == -1 {
		token = message
	} else {
		token = message[0:index]
	}

	switch token {
	case ReadyOK:
		u.SetState(ctx, engine.WaitCommand)
	case BestMove:
		u.result.HandleBestMove(ctx, message)
		u.SetState(ctx, engine.WaitCommand)
	case Info:
		u.result.HandleInfo(ctx, message)
	}
}
