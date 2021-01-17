package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/btree"
	"github.com/pingcap/ticdc/cdc/model"
)



const (
	Sum   = iota // sum
	Count        // count
	Avg          // sum, count
	Min          // tree
	Max          // tree

	SumDistinct   // distinct_sum occurs
	CountDistinct // occurs
	AvgDistinct   // occurs
	MinDistinct
	MaxDistinct
)

type AggFuncHandler struct {
	sum         Value
	count       Value
	distinctSum Value
	occurs      map[Value]Value
	tree        *btree.BTree
}

type uint64_t struct {
	v Value
}

func (m uint64_t) Less(item btree.Item) bool {
	return m.v < (item.(uint64_t).v)
}

func (aggHandler *AggFuncHandler) insert(key Value) {
	aggHandler.sum += key
	aggHandler.count++
	aggHandler.occurs[key]++
	if aggHandler.occurs[key] == 1 {
		aggHandler.distinctSum++
		aggHandler.tree.ReplaceOrInsert(uint64_t{key})
	}
}

func (aggHandler *AggFuncHandler) retract(key Value) {
	aggHandler.sum -= key
	aggHandler.count--
	aggHandler.occurs[key]--
	if aggHandler.occurs[key] == 0 {
		aggHandler.distinctSum--
		delete(aggHandler.occurs, key)
		aggHandler.tree.Delete(uint64_t{key})
	}
}

func (aggHandler *AggFuncHandler) getSum() Value {
	return aggHandler.sum
}

type MVHandler struct {
	cols     []uint16
	funs     []uint16
	handlers []AggFuncHandler
}

func newMVHandler() *MVHandler {
	handler := &MVHandler{}
	handler.createMVHandler([]uint16{0}, []uint16{0})
	return handler
}

func (mvHandler *MVHandler) createMVHandler(funs, cols []uint16) {
	if len(funs) != len(cols) {
		panic("funs len != cols len, crawl")
	}
	copy(mvHandler.cols, cols)
	copy(mvHandler.funs, funs)
	mvHandler.handlers = make([]AggFuncHandler, len(cols))
	for i := range mvHandler.handlers {
		mvHandler.handlers[i].occurs = make(map[Value]Value)
		mvHandler.handlers[i].tree = btree.New(2)
	}
}

func (mvHandler *MVHandler) OnRowChanged(row *model.RowChangedEvent) Value {
	if row.PreColumns != nil {
		//fmt.Printf("precolumn type = %v, value = %v, real type = %v\n", row.PreColumns[0].Type, row.PreColumns[0].Value, reflect.TypeOf(row.PreColumns[0].Value).String())
		if t, ok := row.PreColumns[0].Value.(json.Number); ok {
		//	fmt.Printf("will retract an value %v\n", t)
			if i, err := t.Int64(); err == nil {
				mvHandler.handlers[0].retract(Value(i))
			}
		} else {
			fmt.Println("not a json number, crawl!")
		}
	}
	if row.Columns != nil {
		//fmt.Printf("newcolumn type = %v, value = %v, real type = %v\n", row.Columns[0].Type, row.Columns[0].Value, reflect.TypeOf(row.Columns[0].Value).String())
		if t, ok := row.Columns[0].Value.(json.Number); ok {
			//fmt.Printf("will insert an value %v\n", t)
			if i, err := t.Int64(); err == nil {
				mvHandler.handlers[0].insert(Value(i))
			}
		} else {
			//fmt.Println("not a json number, crawl!")
		}
	}
	fmt.Printf("value in this time: %v\n",mvHandler.handlers[0].getSum())
	return Value(mvHandler.handlers[0].getSum())
}
