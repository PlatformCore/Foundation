package supervisor

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Task func(context.Context) error

type Policy struct {
	MaxRestarts int
	Backoff     time.Duration
}

type Supervisor struct {
	mu     sync.Mutex
	tasks  map[string]Task
	policy Policy
}

func New(policy Policy) *Supervisor {
	if policy.Backoff <= 0 {
		policy.Backoff = time.Second
	}
	if policy.MaxRestarts <= 0 {
		policy.MaxRestarts = 3
	}
	return &Supervisor{tasks: map[string]Task{}, policy: policy}
}

func (s *Supervisor) Add(name string, task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[name] = task
}

func (s *Supervisor) Run(ctx context.Context) error {
	s.mu.Lock()
	tasks := make(map[string]Task, len(s.tasks))
	for k, v := range s.tasks {
		tasks[k] = v
	}
	s.mu.Unlock()
	if len(tasks) == 0 {
		return errors.New("no supervised tasks")
	}
	errCh := make(chan error, len(tasks))
	var wg sync.WaitGroup
	for name, task := range tasks {
		name, task := name, task
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- s.runTask(ctx, name, task)
		}()
	}
	go func() { wg.Wait(); close(errCh) }()
	for err := range errCh {
		if err != nil && !errors.Is(err, context.Canceled) {
			return err
		}
	}
	return ctx.Err()
}

func (s *Supervisor) runTask(ctx context.Context, name string, task Task) error {
	restarts := 0
	for {
		err := task(ctx)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err == nil {
			return nil
		}
		restarts++
		if restarts > s.policy.MaxRestarts {
			return errors.New(name + ": restart limit exceeded: " + err.Error())
		}
		timer := time.NewTimer(s.policy.Backoff * time.Duration(restarts))
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}
