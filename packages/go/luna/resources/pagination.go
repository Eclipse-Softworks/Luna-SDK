package resources

import (
	"context"
)

// Paginator provides an iterator interface for pagination
type Paginator[T any] struct {
	fetchNext func(context.Context, string) (*ListResponse[T], error)
	ctx       context.Context
	buffer    []T
	cursor    string
	hasMore   bool
	done      bool
	err       error
}

type ListResponse[T any] struct {
	Data       []T     `json:"data"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}

func NewPaginator[T any](ctx context.Context, fetchNext func(context.Context, string) (*ListResponse[T], error)) *Paginator[T] {
	return &Paginator[T]{
		fetchNext: fetchNext,
		ctx:       ctx,
		hasMore:   true, // Start with true to allow first fetch
	}
}

// Next advances the iterator to the next item
func (p *Paginator[T]) Next() bool {
	if len(p.buffer) > 0 {
		return true
	}

	if p.done || !p.hasMore {
		return false
	}

	page, err := p.fetchNext(p.ctx, p.cursor)
	if err != nil {
		p.err = err
		p.done = true
		return false
	}

	if len(page.Data) == 0 {
		p.done = true
		return false
	}

	p.buffer = page.Data
	p.hasMore = page.HasMore
	if page.NextCursor != nil {
		p.cursor = *page.NextCursor
	} else {
		p.cursor = ""
	}
	return true
}

// Value returns the current item
func (p *Paginator[T]) Value() T {
	if len(p.buffer) == 0 {
		var zero T
		return zero
	}
	item := p.buffer[0]
	p.buffer = p.buffer[1:]
	return item
}

// Err returns any error occurred during iteration
func (p *Paginator[T]) Err() error {
	return p.err
}
