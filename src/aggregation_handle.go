package main

import (
	"encoding/json"
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

func (aggHandler *AggFuncHandler) getSum() Value { return aggHandler.sum }
func (aggHandler *AggFuncHandler) getCount() Value { return aggHandler.count}
func (aggHandler *AggFuncHandler) getAvg() float64 { return float64(aggHandler.sum)/float64(aggHandler.count) }
func (aggHandler *AggFuncHandler) getMin() Value {
	if aggHandler.count>0 {
		return aggHandler.tree.Min().(uint64_t).v
	} else {
		return Value(0)
	}
}
func (aggHandler *AggFuncHandler) getMax() Value {
	if aggHandler.count>0 {
		return aggHandler.tree.Max().(uint64_t).v
	} else {
		return Value(0)
	}
}
func (aggHandler *AggFuncHandler) getDistinctSum() Value {return aggHandler.distinctSum}
func (aggHandler *AggFuncHandler) getDistinctCount() Value {return Value(len(aggHandler.occurs))}
func (aggHandler *AggFuncHandler) getDistinctAvg() float64 { return float64(aggHandler.distinctSum)/float64(len(aggHandler.occurs)) }
func (aggHandler *AggFuncHandler) getDistinctMin() Value { return aggHandler.tree.Min().(uint64_t).v }
func (aggHandler *AggFuncHandler) getDistinctMax() Value { return aggHandler.tree.Max().(uint64_t).v }

type MVHandler struct {
	cols     []uint16
	funs     []uint16
	handlers []AggFuncHandler
}

// sum(a), max(b), distinct count(c)
func newMVHandler() *MVHandler {
	handler := &MVHandler{}
	handler.createMVHandler([]uint16{0, 1, 2}, []uint16{0, 4, 6})
	return handler
}

func (mvHandler *MVHandler) createMVHandler(cols, funs []uint16) {
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

func (mvHandler *MVHandler) OnRowChanged(row *model.RowChangedEvent) []Value {
	for cid := range row.PreColumns {
		if t, ok := row.PreColumns[cid].Value.(json.Number); ok {
			if i, err := t.Int64(); err == nil {
				mvHandler.handlers[cid].retract(Value(i))
			}
		}
	}
	for cid := range row.Columns {

		if t, ok := row.Columns[cid].Value.(json.Number); ok {
			if i, err := t.Int64(); err == nil {
				//fmt.Printf("here cid %v, value %v\n", cid, i)
				mvHandler.handlers[cid].insert(Value(i))
			}
		}
	}
	//if row.PreColumns != nil {
	//	if t, ok := row.PreColumns[0].Value.(json.Number); ok {
	//		if i, err := t.Int64(); err == nil {
	//			mvHandler.handlers[0].retract(Value(i))
	//		}
	//	}
	//}
	//if row.Columns != nil {
	//	if t, ok := row.Columns[0].Value.(json.Number); ok {
	//		if i, err := t.Int64(); err == nil {
	//			mvHandler.handlers[0].insert(Value(i))
	//		}
	//	}
	//}
	//fmt.Printf("value in this time: %v\n", mvHandler.handlers[0].getSum())
	return []Value{mvHandler.handlers[2].getSum(), mvHandler.handlers[1].getMax(), mvHandler.handlers[0].getDistinctCount()}
}
