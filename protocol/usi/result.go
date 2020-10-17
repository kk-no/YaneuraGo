package usi

import (
	"context"
	"strings"
)

type ThinkResult struct {
	BestMove    string
	Ponder      string
	LastReceive string
	Pvs         []string
}

func NewResult() *ThinkResult {
	return &ThinkResult{}
}

func (tr *ThinkResult) HandleBestMove(ctx context.Context, message string) {
	messages := strings.Split(message, " ")
	if len(messages) >= 4 && messages[2] == "ponder" {
		tr.Ponder = messages[3]
		tr.BestMove = messages[1]
		return
	}

	if len(messages) >= 2 {
		tr.BestMove = messages[1]
		return
	}

	tr.Ponder = "none"
	tr.BestMove = "none"
}

func (tr *ThinkResult) HandleInfo(ctx context.Context, message string) {}
