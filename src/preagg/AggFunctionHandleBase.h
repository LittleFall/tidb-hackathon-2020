//
// Created by littlefall on 2021/1/15.
//

#pragma once
#include <string>
#include <set>

class AggFunctionHandleBase {

public:
    using Key = uint64_t;

    virtual void createAccumulator() = 0;
    virtual void accumulate(Key input) = 0;
    virtual void retract(Key input) = 0;
    virtual std::optional<Key> getValue() = 0;
};

