//
// Created by littlefall on 2021/1/15.
//

#include "AggFunctionHandleBase.h"

class CountAggFunctionHandle : public AggFunctionHandleBase{

    Key accumulator{};

    void createAccumulator() override {
        accumulator = 0;
    }
    void accumulate(Key input) override {
        accumulator ++;
    }
    void retract(Key input) override {
        accumulator --;
    }
    std::optional<Key> getValue() override {
        return accumulator;
    }
};