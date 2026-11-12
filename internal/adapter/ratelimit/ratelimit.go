// Package ratelimit provides the in-memory sliding-window limiter.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter allows at most Max requests per Window per key.
type Limiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	max    int
	window time.Duration
	now    func() time.Time
}

func New(max int, window time.Duration) *Limiter {
	if window <= 0 {
		window = time.Minute
	}
	return &Limiter{hits: make(map[string][]time.Time), max: max, window: window, now: time.Now}
}

func (l *Limiter) Allow(key string) bool {
	if l.max <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)
	entries := l.hits[key]
	kept := entries[:0]
	for _, entry := range entries {
		if entry.After(cutoff) {
			kept = append(kept, entry)
		}
	}
	if len(kept) >= l.max {
		l.hits[key] = kept
		return false
	}
	l.hits[key] = append(kept, now)
	if len(l.hits) > 10_000 {
		l.sweep(cutoff)
	}
	return true
}

func (l *Limiter) sweep(cutoff time.Time) {
	for key, entries := range l.hits {
		if len(entries) == 0 || entries[len(entries)-1].Before(cutoff) {
			delete(l.hits, key)
		}
	}
}
