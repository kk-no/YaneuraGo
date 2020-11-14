package usi

import (
	"context"
	"strings"
)

type ResultManager interface {
	ReceiveMessage(msg string)
	HandleBestMove(ctx context.Context, message string)
	HandleInfo(ctx context.Context, message string)
}

type result struct {
	BestMove    string
	Ponder      string
	LastReceive string
	Pvs         []string
}

func NewResultManager() ResultManager {
	return &result{}
}

func (r *result) ReceiveMessage(msg string) {
	r.LastReceive = msg
}

func (r *result) HandleBestMove(ctx context.Context, message string) {
	messages := strings.Split(message, " ")
	if len(messages) >= 4 && messages[2] == "ponder" {
		r.Ponder = messages[3]
		r.BestMove = messages[1]
		return
	}

	if len(messages) >= 2 {
		r.BestMove = messages[1]
		return
	}

	r.Ponder = "none"
	r.BestMove = "none"
}

func (r *result) HandleInfo(ctx context.Context, message string) {}
