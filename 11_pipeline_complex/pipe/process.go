package pipe

import (
	"context"
	"errors"
)

// All our processes will use this as a base.
type Process struct {
	context context.Context
}

// A default implementation with an error nag. Not required, but serves as an example.
func (p *Process) Run() (<-chan error, error) {
	return nil, errors.New("process needs to implement Run()")
}

// We set our context so that it can be used in the processes to terminate goroutines
// if it needs to.
func (p *Process) SetContext(ctx context.Context) {
	p.context = ctx
}

// Any process needs to implement the Processor interface.
type Processor interface {
	Run() (<-chan error, error)
	SetContext(ctx context.Context)
}