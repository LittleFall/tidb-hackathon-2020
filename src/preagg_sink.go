package main

import (
	"container/heap"
	"context"

	"github.com/pingcap/log"
	"github.com/pingcap/ticdc/cdc/model"
	"github.com/pingcap/ticdc/pkg/config"
	"github.com/pingcap/ticdc/pkg/filter"
	"go.uber.org/zap"
)

// An Item is something we manage in a priority queue.
type Item struct {
	row *model.RowChangedEvent // The value of the item; arbitrary.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].row.CommitTs < pq[j].row.CommitTs
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type preAggSink struct {
	filter *filter.Filter
	config *config.ReplicaConfig
	preAgg *PreAggregate
	pq     PriorityQueue
}

func newPreAggSink(filter *filter.Filter, config *config.ReplicaConfig, preAgg *PreAggregate) (*preAggSink, error) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	p := &preAggSink{
		filter,
		config,
		preAgg,
		pq,
	}
	return p, nil
}

func (p *preAggSink) Initialize(ctx context.Context, tableInfo []*model.SimpleTableInfo) error {
	return nil
}

// EmitRowChangedEvents sends Row Changed Event to Sink
// EmitRowChangedEvents may write rows to downstream directly;
func (p *preAggSink) EmitRowChangedEvents(ctx context.Context, rows ...*model.RowChangedEvent) error {
	for _, row := range rows {
		if p.filter.ShouldIgnoreDMLEvent(row.StartTs, row.Table.Schema, row.Table.Table) {
			log.Info("Row changed event ignored", zap.Uint64("start-ts", row.StartTs))
			continue
		}
		heap.Push(&p.pq, &Item{
			row,
			0,
		})
	}
	return nil
}

// EmitDDLEvent sends DDL Event to Sink
// EmitDDLEvent should execute DDL to downstream synchronously
func (p *preAggSink) EmitDDLEvent(ctx context.Context, ddl *model.DDLEvent) error {
	return nil
}

// FlushRowChangedEvents flushes each row which of commitTs less than or equal to `resolvedTs` into downstream.
// TiCDC guarantees that all of Event which of commitTs less than or equal to `resolvedTs` are sent to Sink through `EmitRowChangedEvents`
func (p *preAggSink) FlushRowChangedEvents(ctx context.Context, resolvedTs uint64) (uint64, error) {
	for {
		if p.pq.Len() == 0 {
			break
		}
		item := heap.Pop(&p.pq).(*Item)
		if item.row.CommitTs > resolvedTs {
			heap.Push(&p.pq, item)
			break
		}
		p.preAgg.rowChange(item.row)
	}

	p.preAgg.flushResolvedTs(resolvedTs)
	return resolvedTs, nil
}

// EmitCheckpointTs sends CheckpointTs to Sink
// TiCDC guarantees that all Events **in the cluster** which of commitTs less than or equal `checkpointTs` are sent to downstream successfully.
func (p *preAggSink) EmitCheckpointTs(ctx context.Context, ts uint64) error {
	log.Info("emit checkpoint ts", zap.Uint64("checkpoint-ts", ts))
	return nil
}

// Close closes the Sink
func (p *preAggSink) Close() error {
	log.Info("close")
	return nil
}
