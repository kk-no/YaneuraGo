package usi

import (
	"bufio"
	"context"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/kk-no/YaneuraGo/dir"

	"github.com/kk-no/YaneuraGo/protocol/state/engine"
)

type Engine interface {
	SetState(ctx context.Context, state engine.State)
	Connect(ctx context.Context, path string) error
	Disconnect(ctx context.Context) error
	IsConnected(ctx context.Context) bool
	SendCommand(ctx context.Context, command string)
	WriteProcess(ctx context.Context)
	ReadProcess(ctx context.Context)
}

type usi struct {
	state   engine.State
	options map[string]string
	// TODO: Define the parts of the subprocess involved separately
	process   *exec.Cmd
	procIn    io.WriteCloser
	procOut   io.ReadCloser
	sendQueue chan string
	result    *ThinkResult
	isDebug   bool
}

func New() Engine {
	return &usi{
		state:     engine.Disconnected,
		options:   nil,
		process:   nil,
		procIn:    nil,
		procOut:   nil,
		sendQueue: make(chan string),
		result:    NewResult(),
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

	f, err := dir.ChangeDir(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := f(); err != nil {
			log.Println("Failed to return original directory, cause by", err)
		}
	}()

	u.process = exec.CommandContext(ctx, engine.Binary)

	if u.procIn, err = u.process.StdinPipe(); err != nil {
		log.Println("Failed to get std in pipe, cause by", err)
		return err
	}

	if u.procOut, err = u.process.StdoutPipe(); err != nil {
		log.Println("Failed to get std out pipe, cause by", err)
		return err
	}

	if err := u.process.Start(); err != nil {
		log.Println("Failed to start process, cause by", err)
		return err
	}

	u.SetState(ctx, engine.Connected)

	go u.WriteProcess(ctx)
	go u.ReadProcess(ctx)

	return nil
}

func (u *usi) Disconnect(ctx context.Context) error {
	if u.procIn != nil {
		if err := u.procIn.Close(); err != nil {
			return err
		}
	}
	if u.procOut != nil {
		if err := u.procOut.Close(); err != nil {
			return err
		}
	}
	u.process = nil
	u.SetState(ctx, engine.Disconnected)
	return nil
}

func (u *usi) IsConnected(ctx context.Context) bool {
	return u.process != nil
}

func (u *usi) SendCommand(ctx context.Context, command string) {
	u.sendQueue <- command
}

func (u *usi) WriteProcess(ctx context.Context) {
	for {
		command := <-u.sendQueue
		// FIXME: Move to HandleCommand() or other
		var token string
		if index := strings.Index(command, " "); index == -1 {
			token = command
		} else {
			token = command[0:index]
		}

		switch token {
		case Stop:
			if u.state != engine.WaitBestMove {
				continue
			}
		case Go:
			u.SetState(ctx, engine.WaitBestMove)
		case Position:
			u.SetState(ctx, engine.WaitCommand)
		case Moves, Side:
			u.SetState(ctx, engine.WaitOneLine)
		case NewGame, GameOver:
			u.SetState(ctx, engine.WaitCommand)
		}

		if _, err := u.procIn.Write([]byte(command + "\n")); err != nil {
			log.Println("Failed to write std in, cause by", err)
			break
		}

		if u.isDebug {
			log.Println(">", command)
		}

		if token == Quit {
			u.SetState(ctx, engine.Disconnected)
			break
		}
	}
}

func (u *usi) ReadProcess(ctx context.Context) {
	scanner := bufio.NewScanner(u.procOut)
	for scanner.Scan() {
		if u.isDebug {
			log.Println("<", scanner.Text())
		}
		u.HandleMessage(ctx, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("Failed to read std out, cause by", err)
	}
}

func (u *usi) HandleMessage(ctx context.Context, message string) {
	u.result.LastReceive = message

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
