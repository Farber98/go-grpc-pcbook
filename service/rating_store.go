package service

import "sync"

type RatingStore interface {
	Add(laptopId string, score float64) (*Rating, error)
}

type MemoryRatingStore struct {
	mutex   sync.RWMutex
	ratings map[string]*Rating
}

type Rating struct {
	count int32
	score float64
}

func NewMemoryRatingStore() *MemoryRatingStore {
	return &MemoryRatingStore{ratings: make(map[string]*Rating)}
}

func (m *MemoryRatingStore) Add(laptopId string, score float64) (*Rating, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	rating, ok := m.ratings[laptopId]
	if !ok {
		rating = &Rating{count: 1, score: score}
	} else {
		rating.count++
		rating.score += score
	}

	m.ratings[laptopId] = rating

	return rating, nil
}
