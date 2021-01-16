// TxnSolver.cpp
// Created by littlefall on 2021/1/13.

#pragma once
#include <string>
#include <queue>
#include <iostream>
#include <unordered_map>
#include <cassert>
#include "../TobeImplement.cpp"

class TxnSolver
{
    using Key = std::string;
    using TimeStamp = uint64_t;
    using Event = std::string;

    struct Lock {
        Key key;
        TimeStamp prewrite_ts;
        Event event;

        bool operator<(const Lock & another) const {
            return prewrite_ts > another.prewrite_ts;
        }
    };

public:
    void prewrite(Key key, TimeStamp ts, Event event);
    void commit(Key key, TimeStamp ts);
    void apply(TimeStamp ts, Event event);

private:
    std::priority_queue<Lock> lock_heap;
    std::unordered_map<Key, TimeStamp> commit_hashtable;
};
