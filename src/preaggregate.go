package main

import (
	"sync"

	"github.com/google/btree"
	"github.com/pingcap/ticdc/cdc/model"
)

type Value uint64

type PreAggregateResult struct {
	ts uint64
	v  Value
}

func (m *PreAggregateResult) Less(item btree.Item) bool {
	return m.ts < (item.(*PreAggregateResult)).ts
}

type PreAggregateMVCC struct {
	mutex      sync.Mutex
	results    *btree.BTree
	resolvedTs uint64
}

func NewPreAggregateMVCC() *PreAggregateMVCC {
	return &PreAggregateMVCC{
		results: btree.New(2),
	}
}

// Retry until success
func (m *PreAggregateMVCC) FindValue(readTs uint64) *Value {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if readTs > m.resolvedTs {
		return nil
	}
	var v *Value
	v = nil
	m.results.DescendLessOrEqual(&PreAggregateResult{
		ts: readTs,
	},
		func(a btree.Item) bool {
			res := a.(*PreAggregateResult)
			if res.ts < readTs {
				*v = res.v
				return false
			}
			return true
		},
	)
	return v
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
	another := m.results.ReplaceOrInsert(result)
	if another != nil {
		panic("commit ts same, result %v, crawl!")
	}
}

func (m *PreAggregateMVCC) UpdateResolveTs(resolvedTs uint64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.resolvedTs > resolvedTs {
		panic("resolve ts smaller, crawl!")
	}
	m.resolvedTs = resolvedTs
}

type AggregateHandler interface {
	OnRowChanged(row *model.RowChangedEvent) Value
}

type PreAggregate struct {
	preaggMVCC *PreAggregateMVCC
	handler    *AggregateHandler
}

func NewPreAggregate(preaggMVCC *PreAggregateMVCC, handler *AggregateHandler) *PreAggregate {
	return &PreAggregate{
		preaggMVCC,
		handler,
	}
}

func (p *PreAggregate) rowChange(row *model.RowChangedEvent) {
	// 算聚合答案
	// 调用这个 interface
	v := p.handler.OnRowChanged(row)

	p.preaggMVCC.AddValue(&PreAggregateResult{
		ts: row.CommitTs,
		v:  v,
	})
}

func (p *PreAggregate) flushResolvedTs(resolvedTs uint64) {
	p.preaggMVCC.UpdateResolveTs(resolvedTs)
}
