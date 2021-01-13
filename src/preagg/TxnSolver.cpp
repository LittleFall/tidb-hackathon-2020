// TxnSolver.cpp
// Created by littlefall on 2021/1/13.


#include "TxnSolver.h"

void TxnSolver::preWrite(Key key, TimeStamp ts, Event event)
{
    lock_heap.push({key, ts, event});
}

void TxnSolver::commit(Key key, TimeStamp commit_ts)
{
    commit_hashtable[key] = commit_ts;
    while (!lock_heap.empty() && commit_hashtable[lock_heap.top().key]) {
        apply(commit_ts, lock_heap.top().event);
        commit_hashtable[lock_heap.top().key] = 0;
        lock_heap.pop();
    }
}

void TxnSolver::apply(TimeStamp ts, Event event)
{
    static TimeStamp applied_ts = -1;
    assert(applied_ts < ts);
    applied_ts = ts;

    std::cout << "apply event:(" << event << "), commit_ts:" << ts << std::endl;
}