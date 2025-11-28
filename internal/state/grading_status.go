package state

import (
	"context"
	"sync"
)

type GradingState struct {
	isGrading bool
	status    string
	mu        *sync.RWMutex
	cancel    context.CancelFunc
	Ctx       context.Context
}

func NewGradingLock() *GradingState {
	return &GradingState{
		mu: &sync.RWMutex{},
	}
}

func (l *GradingState) Lock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.Ctx, l.cancel = context.WithCancel(context.Background())

	l.isGrading = true
}

func (l *GradingState) Unlock() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.isGrading = false
}

func (l *GradingState) IsGrading() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.isGrading
}

func (l *GradingState) SetStatus(status string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.status = status
}

func (l *GradingState) GetStatus() string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.status
}

func (l *GradingState) Cancel() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cancel()
}
