// TxnSolver.cpp
// Created by littlefall on 2021/1/13.


#include "TxnSolver.h"

void TxnSolver::prewrite(Key key, TimeStamp ts, Event event)
{
    lock_heap.push({key, ts, event});
}

void TxnSolver::commit(Key key, TimeStamp commit_ts)
{
    commit_hashtable[key] = commit_ts;
    while (!lock_heap.empty() && commit_hashtable.count(lock_heap.top().key)) {
        apply(commit_ts, lock_heap.top().event);
        commit_hashtable.erase(lock_heap.top().key);
        lock_heap.pop();
    }
}

void TxnSolver::apply(TimeStamp ts, Event event)
{
    static TimeStamp applied_ts = 0;
    assert(applied_ts <= ts);
    applied_ts = ts;

    std::cout << "apply event:(" << event << "), commit_ts:" << ts << std::endl;
}
