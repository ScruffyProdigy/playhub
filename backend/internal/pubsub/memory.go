package pubsub

import (
	"context"
	"sync"
)

// Memory is an in-process broker for local development and tests.
type Memory struct {
	mu   sync.RWMutex
	subs map[string]map[chan []byte]struct{}
}

// NewMemory creates an in-memory pub/sub broker.
func NewMemory() *Memory {
	return &Memory{subs: make(map[string]map[chan []byte]struct{})}
}

func (m *Memory) Publish(_ context.Context, channel string, payload []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for ch := range m.subs[channel] {
		data := append([]byte(nil), payload...)
		select {
		case ch <- data:
		default:
		}
	}
	return nil
}

func (m *Memory) Subscribe(_ context.Context, channel string) (<-chan []byte, func(), error) {
	ch := make(chan []byte, 16)

	m.mu.Lock()
	if m.subs[channel] == nil {
		m.subs[channel] = make(map[chan []byte]struct{})
	}
	m.subs[channel][ch] = struct{}{}
	m.mu.Unlock()

	unsubscribe := func() {
		m.mu.Lock()
		delete(m.subs[channel], ch)
		if len(m.subs[channel]) == 0 {
			delete(m.subs, channel)
		}
		m.mu.Unlock()
		close(ch)
	}

	return ch, unsubscribe, nil
}

func (m *Memory) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for channel, listeners := range m.subs {
		for ch := range listeners {
			close(ch)
		}
		delete(m.subs, channel)
	}
	return nil
}
