package main

import (
	"fmt"
	"sync"

	"github.com/google/btree"
	"github.com/pingcap/ticdc/cdc/model"
)

type Value int64

type PreAggregateResult struct {
	ts uint64
	vs  []Value
}

func (m *PreAggregateResult) Less(item btree.Item) bool {
	return m.ts < (item.(*PreAggregateResult)).ts
}

type PreAggregateMVCC struct {
	mutex      sync.Mutex
	results    *btree.BTree
	resolvedTs uint64
	chans      []struct {
		ch chan []Value
		ts uint64
	}
}

func NewPreAggregateMVCC() *PreAggregateMVCC {
	return &PreAggregateMVCC{
		results:    btree.New(2),
		resolvedTs: 0,
	}
}

func (m *PreAggregateMVCC) FindValue(readTs uint64) chan []Value {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ch := make(chan []Value, 1)
	if readTs > m.resolvedTs {
		m.chans = append(m.chans, struct {
			ch chan []Value
			ts uint64
		}{
			ch,
			readTs,
		})
		return ch
	}
	var v []Value
	v = nil
	m.results.DescendLessOrEqual(&PreAggregateResult{
		ts: readTs,
	},
		func(a btree.Item) bool {
			res := a.(*PreAggregateResult)
			if res.ts < readTs {
				v = res.vs
				return false
			}
			return true
		},
	)
	ch <- v
	return ch
}

func (m *PreAggregateMVCC) AddValue(result *PreAggregateResult) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	mx := m.results.Max()
	if mx != nil {
		if result.ts < mx.(*PreAggregateResult).ts {
			panic("commit ts smaller, crawl!")
		}
	}
	m.results.ReplaceOrInsert(result)
	if m.results.Len() > 10000 {
		m.results.DeleteMin()
	}
}

func (m *PreAggregateMVCC) UpdateResolveTs(resolvedTs uint64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.resolvedTs > resolvedTs {
		panic("resolve ts smaller, crawl!")
	}
	m.resolvedTs = resolvedTs
	for _, s := range m.chans {
		if m.resolvedTs < s.ts {
			continue
		}
		var v []Value
		v = nil
		m.results.DescendLessOrEqual(&PreAggregateResult{
			ts: s.ts,
		},
			func(a btree.Item) bool {
				res := a.(*PreAggregateResult)
				if res.ts < s.ts {
					v = res.vs
					return false
				}
				return true
			},
		)
		s.ch <- v
	}
}

type AggregateHandler interface {
	OnRowChanged(row *model.RowChangedEvent) []Value
}

type PreAggregate struct {
	preaggMVCC *PreAggregateMVCC
	handler    AggregateHandler
}

func NewPreAggregate(preaggMVCC *PreAggregateMVCC, handler AggregateHandler) *PreAggregate {
	return &PreAggregate{
		preaggMVCC,
		handler,
	}
}

func (p *PreAggregate) rowChange(row *model.RowChangedEvent) {
	vs := p.handler.OnRowChanged(row)
	fmt.Printf("CommitTs: %v: ( ", row.CommitTs)
	for _, v := range vs {
		fmt.Printf("%v ", v)
	}
	fmt.Printf(")\n")

	p.preaggMVCC.AddValue(&PreAggregateResult{
		ts: row.CommitTs,
		vs:  vs,
	})
}

func (p *PreAggregate) flushResolvedTs(resolvedTs uint64) {
	p.preaggMVCC.UpdateResolveTs(resolvedTs)
}
