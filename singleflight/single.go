// Package singleflight provides in-process duplicate call suppression.
//
// Group is an injectable type (no package-level state): construct one per
// component, or embed the zero value, so tests and multiple instances stay isolated.
package singleflight

import xsf "golang.org/x/sync/singleflight"

// Result holds the results of DoChan.
type Result = xsf.Result

// Doer is the behaviour consumers depend on; *Group implements it.
type Doer interface {
	Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)
}

// Group coalesces concurrent calls that share a key. The zero value is ready
// to use. A Group must not be copied after first use.
type Group struct {
	g xsf.Group
}

var _ Doer = (*Group)(nil)

// NewGroup returns a new empty Group.
func NewGroup() *Group { return &Group{} }

// Do runs fn once for concurrent callers using the same key and shares the result.
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	return g.g.Do(key, fn)
}

// DoChan is like Do but returns a channel that receives the result.
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
	return g.g.DoChan(key, fn)
}

// Forget tells the group to forget key so the next call runs fn again.
func (g *Group) Forget(key string) { g.g.Forget(key) }
