// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

// the package allows you to keep records of running goroutines
// it was useful when you need wait when all your goroutines will finished

package subprocmgr

import (
	"sync"
)

// Contains list of names of gorutines
type Goroutines struct {
	sync.RWMutex
	done chan bool
	List map[string]bool

	onceNew    sync.Once
	onceRemove sync.Once
}

// initializing a new goroutines list
func New() *Goroutines {
	g := &Goroutines{}
	g.List = make(map[string]bool)
	return g
}

// Add add name of goroutine to list
func (g *Goroutines) Add(name string) {
	// initialisation of Done channel
	g.onceNew.Do(func() {
		g.done = make(chan bool)
	})

	g.Lock()
	defer g.Unlock()

	g.List[name] = true
}

// Remove remove name of goroutine from list
func (g *Goroutines) Remove(name string) {
	g.Lock()
	delete(g.List, name)
	g.Unlock()

	// if all goroutines removed
	if g.Empty() {
		g.onceRemove.Do(func() {
			g.done <- true
		})
	}
}

// return channel wich set up true when all goroutines removed
func (g *Goroutines) Done() chan bool {
	return g.done
}

// Empty check list for empty
func (g *Goroutines) Empty() bool {
	g.RLock()
	defer g.RUnlock()
	return len(g.List) == 0
}

// Len return length of list
func (g *Goroutines) Len() int {
	g.RLock()
	defer g.RUnlock()
	return len(g.List)
}
