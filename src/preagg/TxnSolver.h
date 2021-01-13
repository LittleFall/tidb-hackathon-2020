// TxnSolver.cpp
// Created by littlefall on 2021/1/13.

#pragma once
#include <string>
#include <queue>
#include <iostream>
#include <unordered_map>
#include "../TobeImplement.cpp"

class TxnSolver
{
    using Key = std::string;
    using TimeStamp = int;
    using Event = std::string;

    struct Lock {
        Key key;
        TimeStamp pre_write_ts;
        Event event;

        bool operator<(const Lock & another) const {
            return pre_write_ts > another.pre_write_ts;
        }
    };

public:
    void preWrite(Key key, TimeStamp ts, Event event);
    void commit(Key key, TimeStamp ts);
    void apply(TimeStamp ts, Event event);

private:
    std::priority_queue<Lock> lock_heap;
    std::unordered_map<Key, TimeStamp> commit_hashtable;
};

