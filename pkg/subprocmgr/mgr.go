// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// yuriyiovkov@gmail.com; telegram: @yuriyiovkov

// the package allows you to monitor of running goroutines
// it was useful when you need wait when all your goroutines will finished
// look Test_Usage* in mgr_test.go file for usage

package subprocmgr

import (
	"sync"
)

// Goroutines Contains list of names of goroutines
type Goroutines struct {
	sync.RWMutex
	done chan bool
	List map[string]bool

	onceNew   sync.Once
	onceClose sync.Once
}

// New initializing a new goroutines list
func New() *Goroutines {
	g := &Goroutines{}

	// initialisation of Done channel
	g.onceNew.Do(func() {
		g.done = make(chan bool)
	})

	g.Lock()
	g.List = make(map[string]bool)
	g.Unlock()
	return g
}

// Add add name of goroutine to list
func (g *Goroutines) Add(name string) {

	g.Lock()
	defer g.Unlock()

	g.List[name] = true
}

// Remove remove name of goroutine from list
func (g *Goroutines) Remove(name string) {
	g.Lock()
	delete(g.List, name)
	g.Unlock()

	// empty test. Not in lock block because Empty user RLock(else deadlock)
	if g.Empty() {
		g.closeChan()
	}
}

// return channel witch set up true when all goroutines removed
func (g *Goroutines) Done() chan bool {
	// if all goroutines removed
	if g.Empty() {
		g.closeChan()
	}
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

// closeChan close channel done one time
func (g *Goroutines) closeChan() {
	g.onceClose.Do(func() {
		close(g.done)
	})
}
