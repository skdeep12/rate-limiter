package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

type leakyBucket struct {
	mu                           sync.Mutex
	maxCapacity                  int
	capacity                     int
	rateOfFillingTokensPerSecond int
	lastUpdatedAt                time.Time
}

type Bucket interface {
	ConsumeTokens(count int, at time.Time) bool
}

func (b *leakyBucket) ConsumeTokens(utilisation int, at time.Time) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	fmt.Printf("Consuming tokens %d\n", utilisation)
	var diff time.Duration
	if b.lastUpdatedAt.Before(at) {
		diff = at.Sub(b.lastUpdatedAt)
		b.lastUpdatedAt = at
	}
	b.capacity = b.capacity + int(b.tokenFromDuration(diff))

	fmt.Printf("capacity is %d\n", b.capacity)
	if b.capacity > b.maxCapacity {
		b.capacity = b.maxCapacity
	}
	if utilisation <= b.capacity {
		b.capacity -= utilisation
		return true
	}
	return false
}

func (b *leakyBucket) tokenFromDuration(duration time.Duration) float64 {
	return duration.Seconds() * float64(b.rateOfFillingTokensPerSecond)
}

func NewLeakyBucket(maxCapacity int, rate int) Bucket {
	return &leakyBucket{
		mu:                           sync.Mutex{},
		maxCapacity:                  maxCapacity,
		capacity:                     0,
		rateOfFillingTokensPerSecond: rate,
		lastUpdatedAt:                time.Now(),
	}
}
