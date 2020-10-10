package usi

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os/exec"

	"github.com/kk-no/YaneuraGo/dir"

	"github.com/kk-no/YaneuraGo/protocol/state/engine"
)

// TODO: Add SetState() function
type Engine interface {
	Connect(ctx context.Context, path string) error
	Disconnect(ctx context.Context) error
	IsConnected(ctx context.Context) bool
	SendCommand(ctx context.Context, command string) error
	ReadResult(ctx context.Context) error
}

type usi struct {
	state   engine.State
	options map[string]string
	// TODO: Define the parts of the subprocess involved separately
	process *exec.Cmd
	procIn  io.WriteCloser
	procOut io.ReadCloser
}

func New() Engine {
	return &usi{}
}

func (u *usi) Connect(ctx context.Context, path string) error {
	log.Println("Call connect")
	log.Printf("Specify engine directory: %v\n", path)

	if err := u.Disconnect(ctx); err != nil {
		return err
	}

	u.state = engine.WaitConnecting

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

	u.state = engine.Connected

	return nil
}

func (u *usi) Disconnect(ctx context.Context) error {
	log.Println("Call disconnect")
	u.process = nil
	u.state = engine.Disconnected
	return nil
}

func (u *usi) IsConnected(ctx context.Context) bool {
	return u.process != nil
}

func (u *usi) SendCommand(ctx context.Context, command string) error {
	log.Println("Call SendCommand")
	log.Println(">", command)
	if !u.IsConnected(ctx) {
		return errors.New("process is not started")
	}
	if _, err := u.procIn.Write([]byte(command + "\n")); err != nil {
		log.Println("Failed to write std in, cause by", err)
		return err
	}
	return u.ReadResult(ctx)
}

func (u *usi) ReadResult(ctx context.Context) error {
	log.Println("Call ReadResult")
	if !u.IsConnected(ctx) {
		return errors.New("process is not started")
	}
	scanner := bufio.NewScanner(u.procOut)
	for scanner.Scan() {
		log.Println("<", scanner.Text())

		// FIXME: Unable to finish reading the output
		//  Explicitly specify the last line as a temporary implementation
		if line := scanner.Text(); line == "usiok" || line == "readyok" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Failed to read std out, cause by", err)
		return err
	}
	return nil
}
