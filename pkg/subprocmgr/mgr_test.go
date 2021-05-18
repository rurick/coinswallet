package subprocmgr

import (
	"context"
	"sync"
	"testing"
	"time"
)

// wait 100 Millisecond when all goroutines finished and return
func Test_Usage(t *testing.T) {
	f1 := func() { time.Sleep(10 * time.Millisecond) }
	f2 := func() { time.Sleep(5 * time.Millisecond) }
	f3 := func() { time.Sleep(30 * time.Millisecond) }
	f4 := func() { time.Sleep(15 * time.Millisecond) }

	mgr := New()
	mgr.Add("f1")
	go func() { f1(); mgr.Remove("f1") }()

	mgr.Add("f2")
	go func() { f2(); mgr.Remove("f2") }()

	mgr.Add("f3")
	go func() { f3(); mgr.Remove("f3") }()

	mgr.Add("f4")
	go func() { f4(); mgr.Remove("f4") }()

	// emulation of end of programm
	terminateTimeout := 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), terminateTimeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			t.Errorf("Test_Usage() terminate timeout")
			return
		default:
			if mgr.Empty() {
				return
			}
		}
		time.Sleep(50 * time.Microsecond) // CPU hold
	}
}

// test suck as Test_Usage using Done() function
func Test_UsageDone(t *testing.T) {
	f1 := func() { time.Sleep(10 * time.Millisecond) }
	f2 := func() { time.Sleep(5 * time.Millisecond) }
	f3 := func() { time.Sleep(30 * time.Millisecond) }
	f4 := func() { time.Sleep(15 * time.Millisecond) }

	mgr := New()
	mgr.Add("f1")
	go func() { f1(); mgr.Remove("f1") }()

	mgr.Add("f2")
	go func() { f2(); mgr.Remove("f2") }()

	mgr.Add("f3")
	go func() { f3(); mgr.Remove("f3") }()

	mgr.Add("f4")
	go func() { f4(); mgr.Remove("f4") }()

	// emulation of end of programm
	terminateTimeout := 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), terminateTimeout)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Errorf("Test_Usage() terminate timeout")
		return
	case <-mgr.Done():
		return
	}
}

// wait 10 Millisecond when all goroutines finished and return by timeout
func Test_UsageTimeoutErr(t *testing.T) {
	f1 := func() { time.Sleep(10 * time.Millisecond) }
	f2 := func() { time.Sleep(5 * time.Millisecond) }
	f3 := func() { time.Sleep(30 * time.Millisecond) }
	f4 := func() { time.Sleep(15 * time.Millisecond) }

	mgr := New()
	mgr.Add("f1")
	go func() { f1(); mgr.Remove("f1") }()

	mgr.Add("f2")
	go func() { f2(); mgr.Remove("f2") }()

	mgr.Add("f3")
	go func() { f3(); mgr.Remove("f3") }()

	mgr.Add("f4")
	go func() { f4(); mgr.Remove("f4") }()

	// emulation of end of programm
	terminateTimeout := 10 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), terminateTimeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if mgr.Empty() {
				t.Errorf("Test_UsageTimeoutErr() no timeout")
				return
			}
		}
		time.Sleep(50 * time.Microsecond) // CPU hold
	}
}

// Test_Sync testing sync list manipulations
func Test_Sync(t *testing.T) {
	mgr := New()
	mgr.Add("func1")
	mgr.Add("func2")
	mgr.Add("func1")

	t.Run("sync add 2 different func 3 times", func(t *testing.T) {
		if mgr.Len() != 2 {
			t.Errorf("Add() error: add 2 diffrrent func, got: %d", mgr.Len())
		}
	})

	mgr.Remove("func2")
	mgr.Remove("func1")
	mgr.Remove("func1")

	t.Run("sync remove 2 func 3 times", func(t *testing.T) {
		if !mgr.Empty() {
			t.Errorf("Remove() error: delete all func, got: %d", mgr.Len())
		}
	})

}

// Test_ASync testing async list manipulations
func Test_Async(t *testing.T) {
	mgr := New()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() { mgr.Add("func1"); wg.Done() }()
	wg.Add(1)
	go func() { mgr.Add("func1"); wg.Done() }()
	wg.Add(1)
	go func() { mgr.Add("func1"); wg.Done() }()
	wg.Add(1)
	go func() { mgr.Add("func2"); wg.Done() }()
	wg.Add(1)
	go func() { mgr.Add("func1"); wg.Done() }()
	wg.Add(1)
	go func() { mgr.Add("func2"); wg.Done() }()

	go t.Run("async length of list", func(t *testing.T) {
		t.Log("Len():", mgr.Len())
	})

	wg.Wait()
	t.Run("async add 2 different func 6 times", func(t *testing.T) {
		if mgr.Len() != 2 {
			t.Errorf("Add() error: add 2 diffrrent func, got: %d", mgr.Len())
		}
	})

	wg.Add(1)
	go func() { mgr.Remove("func2"); wg.Done() }()
	wg.Add(1)
	go func() { mgr.Remove("func2"); wg.Done() }()
	wg.Add(1)
	go func() { mgr.Remove("func1"); wg.Done() }()

	wg.Wait()
	t.Run("async remove 2 func 3 times", func(t *testing.T) {
		if !mgr.Empty() {
			t.Errorf("Remove() error: delete all func, got: %d", mgr.Len())
		}
	})

}

// Test_RemoveBeforeAdd
func Test_RemoveBeforeAdd(t *testing.T) {
	mgr := New()
	mgr.Remove("func1")
	<-mgr.Done()
}

// Test_RemoveBeforeAddAndWait
func Test_RemoveBeforeAddAndWait(t *testing.T) {
	mgr := New()
	mgr.Remove("func1")
	<-mgr.Done()
}

// Test_DontAddAndWait
func Test_DontAddAndWait(t *testing.T) {
	mgr := New()
	<-mgr.Done()
}
