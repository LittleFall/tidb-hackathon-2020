//
// Created by littlefall on 2021/1/15.
//

#include "AggFunctionHandleBase.h"

class SumAggFunctionHandle : public AggFunctionHandleBase{

    Key accumulator{};
    uint64_t count = 0;

    void createAccumulator() override {
        accumulator = 0;
    }
    void accumulate(Key input) override {
        accumulator += input;
        ++count;
    }
    void retract(Key input) override {
        accumulator -= input;
        --count;
    }
    std::optional<Key> getValue() override {
        if (!count) {
            return std::nullopt;
        } else {
            return accumulator;
        }
    }
};