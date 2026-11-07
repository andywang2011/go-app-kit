package singleflight_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andywang2011/go-app-kit/singleflight"
)

func TestGroupCoalescesConcurrentCalls(t *testing.T) {
	g := singleflight.NewGroup()
	var runs, shared int32
	started := make(chan struct{})
	release := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err, sh := g.Do("key", func() (interface{}, error) {
				atomic.AddInt32(&runs, 1)
				close(started)
				<-release
				return "ok", nil
			})
			if err != nil || v != "ok" {
				t.Errorf("Do = (%v, %v)", v, err)
			}
			if sh {
				atomic.AddInt32(&shared, 1)
			}
		}()
	}
	<-started
	time.Sleep(20 * time.Millisecond) // let the other callers join the in-flight call
	close(release)
	wg.Wait()
	if runs != 1 {
		t.Fatalf("fn runs: want 1, got %d", runs)
	}
	if shared == 0 {
		t.Fatal("expected shared result")
	}
}

func TestGroupPropagatesError(t *testing.T) {
	g := singleflight.NewGroup()
	wantErr := errors.New("boom")
	_, err, _ := g.Do("err-key", func() (interface{}, error) { return nil, wantErr })
	if !errors.Is(err, wantErr) {
		t.Fatalf("err: %v", err)
	}
}

func TestGroupDistinctKeys(t *testing.T) {
	g := singleflight.NewGroup()
	var runs int32
	var wg sync.WaitGroup
	for _, key := range []string{"a", "b"} {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			_, err, _ := g.Do(k, func() (interface{}, error) {
				atomic.AddInt32(&runs, 1)
				return k, nil
			})
			if err != nil {
				t.Errorf("Do(%s): %v", k, err)
			}
		}(key)
	}
	wg.Wait()
	if runs != 2 {
		t.Fatalf("runs: want 2, got %d", runs)
	}
}

func TestGroupsAreIndependent(t *testing.T) {
	a, b := singleflight.NewGroup(), singleflight.NewGroup()
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _, _ = a.Do("k", func() (interface{}, error) {
			close(started)
			<-release
			return 1, nil
		})
	}()
	<-started
	v, _, shared := b.Do("k", func() (interface{}, error) { return 2, nil })
	if v != 2 || shared {
		t.Fatalf("separate group must not share in-flight call: v=%v shared=%v", v, shared)
	}
	close(release)
	<-done
}

func TestForgetDoChanAndZeroValue(t *testing.T) {
	var g singleflight.Group // zero value is usable
	var d singleflight.Doer = &g
	var runs int32
	fn := func() (interface{}, error) { atomic.AddInt32(&runs, 1); return "v", nil }
	_, _, _ = d.Do("k", fn)
	g.Forget("k")
	_, _, _ = d.Do("k", fn)
	if runs != 2 {
		t.Fatalf("runs = %d", runs)
	}
	res := <-g.DoChan("c", fn)
	if res.Err != nil || res.Val != "v" {
		t.Fatalf("DoChan = %+v", res)
	}
}
