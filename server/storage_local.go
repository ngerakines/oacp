package server

import (
	"context"
	"sync"
)

type localStorage struct {
	Locations map[string]string

	mux sync.Mutex
}

var _ storage = &localStorage{}

func newLocalStorage() (storage, error) {
	return localStorage{
		Locations: make(map[string]string),
	}, nil
}

func (l localStorage) RecordLocation(_ context.Context, state, location string) error {
	l.mux.Lock()
	defer l.mux.Unlock()
	if _, ok := l.Locations[state]; ok {
		return errLocationExists
	}
	l.Locations[state] = location
	return nil
}

func (l localStorage) GetLocation(ctx context.Context, state string) (string, error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	if location, ok := l.Locations[state]; ok {
		delete(l.Locations, state)
		return location, nil
	}
	return "", errLocationNotFound
}

func (l localStorage) Close() error {
	return nil
}
