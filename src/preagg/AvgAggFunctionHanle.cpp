//
// Created by littlefall on 2021/1/15.
//

#include "AggFunctionHandleBase.h"

class AvgAggFunctionHandle : public AggFunctionHandleBase{

    Key sum{};
    uint64_t count = 0;

    void createAccumulator() override {
        sum = count = 0;
    }
    void accumulate(Key input) override {
        sum += input;
        ++count;
    }
    void retract(Key input) override {
        sum -= input;
        --count;
    }
    std::optional<Key> getValue() override {
        if (!count) {
            return std::nullopt;
        } else {
            return (double)sum/count;
        }
        // TODO: return a real double.
    }
};