package service

import (
	"redirect/utils/algorithm"
	"sync"
	"time"
)

//RoundBalanceMap 轮询全局锁
type RoundBalanceMap struct {
	Data      map[string]*algorithm.Round
	Lock      *sync.RWMutex
	UpdatedAt time.Time
}

//WeightBalanceMap 权重全局锁
type WeightBalanceMap struct {
	Data      map[string]*algorithm.WeightRoundRobinBalance
	Lock      *sync.RWMutex
	UpdatedAt time.Time
}

func NewRoundBalanceMap() *RoundBalanceMap {
	return &RoundBalanceMap{
		Data:      make(map[string]*algorithm.Round),
		Lock:      &sync.RWMutex{},
		UpdatedAt: time.Now(),
	}
}

func (d RoundBalanceMap) Get(k string) *algorithm.Round {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	return d.Data[k]
}

func (d RoundBalanceMap) Set(k string, v *algorithm.Round) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Data[k] = v
}

func (d RoundBalanceMap) Del(k string) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	delete(d.Data, k)
}

func NewWeightBalanceMap() *WeightBalanceMap {
	return &WeightBalanceMap{
		Data:      make(map[string]*algorithm.WeightRoundRobinBalance),
		Lock:      &sync.RWMutex{},
		UpdatedAt: time.Now(),
	}
}

func (d WeightBalanceMap) Get(k string) *algorithm.WeightRoundRobinBalance {
	d.Lock.RLock()
	defer d.Lock.RUnlock()
	return d.Data[k]
}

func (d WeightBalanceMap) Set(k string, v *algorithm.WeightRoundRobinBalance) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Data[k] = v
}

func (d WeightBalanceMap) Del(k string) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	delete(d.Data, k)
}
