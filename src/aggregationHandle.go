package main

import (
	"github.com/google/btree"
	"github.com/pingcap/ticdc/cdc/model"
)

const (
	Sum = iota // sum
	Count // count
	Avg // sum, count
	Min // tree
	Max // tree

	SumDistinct // distinct_sum occurs
	CountDistinct // occurs
	AvgDistinct // occurs
	MinDistinct
	MaxDistinct
)

type AggFuncHandler struct {
	sum         uint64
	count       uint64
	distinctSum uint64
	occurs 		map[uint64]uint64
	tree		btree.BTree
}

type uint64_t struct {
	v	uint64
}

func (m uint64_t) Less(item btree.Item) bool {
	return m.v < (item.(uint64_t).v)
}

func (aggHandler *AggFuncHandler) insert(key uint64) {
	aggHandler.sum += key
	aggHandler.count ++
	aggHandler.occurs[key] ++
	if aggHandler.occurs[key] == 1 {
		aggHandler.distinctSum++
		aggHandler.tree.ReplaceOrInsert(uint64_t{key})
	}
}

func (aggHandler *AggFuncHandler) retract(key uint64) {
	aggHandler.sum -= key
	aggHandler.count --
	aggHandler.occurs[key] --
	if aggHandler.occurs[key] == 0 {
		aggHandler.distinctSum --
		delete(aggHandler.occurs, key)
		aggHandler.tree.Delete(uint64_t{key})
	}
}

func (aggHandler *AggFuncHandler) getSum() uint64 {
	return aggHandler.sum
}

type MVHandler struct {
	cols	 []uint16
	funs 	 []uint16
	handlers []AggFuncHandler
}

func (mvHandler *MVHandler) createMVHandler(funs, cols []uint16) {
	if len(funs) != len(cols) {
		panic("funs len != cols len, crawl")
	}
	copy(mvHandler.cols, cols)
	copy(mvHandler.funs, funs)
	mvHandler.handlers = make([]AggFuncHandler, len(cols))
	for i := range mvHandler.handlers {
		mvHandler.handlers[i].occurs = make(map[uint64]uint64)
	}
}

func (mvHandler *MVHandler) OnRowChanged(row *model.RowChangedEvent) Value {
	if row.PreColumns != nil {
		if t, ok := row.PreColumns[0].Value.(uint64); ok {
			mvHandler.handlers[0].retract(t)
		}
	}
	if row.Columns != nil {
		if t, ok := row.PreColumns[0].Value.(uint64); ok {
			mvHandler.handlers[0].insert(t)
		}
	}
	return Value(mvHandler.handlers[0].getSum())
}

