package algorithm

import (
	"errors"
	"sync"
)

type Round struct {
	CurIndex int
	Rss      []string
	mutex    sync.Mutex // add a mutex to ensure concurrent safety
}

func (r *Round) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("at least 1 parameter is required")
	}
	addr := params[0]
	r.mutex.Lock() // lock the mutex before modifying shared state
	r.Rss = append(r.Rss, addr)
	r.mutex.Unlock() // unlock the mutex after modifying shared state
	return nil
}

func (r *Round) Next() (string, error) {
	if len(r.Rss) == 0 {
		return "", errors.New("no parameters exist")
	}
	r.mutex.Lock() // lock the mutex before modifying shared state
	curElement := r.Rss[r.CurIndex]
	r.CurIndex = (r.CurIndex + 1) % len(r.Rss)
	r.mutex.Unlock() // unlock the mutex after modifying shared state
	return curElement, nil
}
