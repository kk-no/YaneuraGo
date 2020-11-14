package usi

import (
	"bufio"
	"context"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/kk-no/YaneuraGo/state/engine"
)

type ReadWriteProcessor interface {
	Start(ctx context.Context)
	Stop() error
	Write(ctx context.Context)
	Read(ctx context.Context)
	SendCommand(ctx context.Context, command string)
}

type process struct {
	cmd       *exec.Cmd
	wg        *sync.WaitGroup
	cancel    context.CancelFunc
	procIn    io.WriteCloser
	procOut   io.ReadCloser
	sendQueue chan string
}

func NewReadWriteProcessor(ctx context.Context) (ReadWriteProcessor, error) {
	p := new(process)
	p.wg = &sync.WaitGroup{}
	p.cmd = exec.CommandContext(ctx, engine.Binary)
	p.sendQueue = make(chan string)

	var err error
	if p.procIn, err = p.cmd.StdinPipe(); err != nil {
		log.Println("Failed to get std in pipe:", err)
		return nil, err
	}

	if p.procOut, err = p.cmd.StdoutPipe(); err != nil {
		log.Println("Failed to get std out pipe:", err)
		return nil, err
	}

	if err := p.cmd.Start(); err != nil {
		log.Println("Failed to start process:", err)
		return nil, err
	}

	return p, nil
}

func (p *process) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	p.wg.Add(1)
	go p.Read(ctx)

	p.wg.Add(1)
	go p.Write(ctx)
}

func (p *process) Stop() error {
	if p.cancel != nil {
		p.cancel()
		// Wait for read and write goroutine close.
		p.wg.Wait()
	}
	if p.sendQueue != nil {
		close(p.sendQueue)
	}
	if p.procIn != nil {
		if err := p.procIn.Close(); err != nil {
			return err
		}
	}
	if p.procOut != nil {
		if err := p.procOut.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *process) Write(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.wg.Done()
			return
		default:
			command := <-p.sendQueue
			log.Println(">", command)
			if _, err := p.procIn.Write([]byte(command + "\n")); err != nil {
				log.Println("Failed to write std in:", err)
				return
			}
		}
	}
}

func (p *process) Read(ctx context.Context) {
	scanner := bufio.NewScanner(p.procOut)
	for {
		select {
		case <-ctx.Done():
			p.wg.Done()
			return
		default:
			if scanner.Scan() {
				log.Println("<", scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				log.Println("Failed to read std out:", err)
				return
			}
		}
	}
}

func (p *process) SendCommand(ctx context.Context, command string) {
	p.sendQueue <- command
}
