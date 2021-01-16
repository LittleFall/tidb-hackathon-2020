//
// Created by littlefall on 2021/1/15.
//

#include "AggFunctionHandleBase.h"

class MaxAggFunctionHandle : public AggFunctionHandleBase{

    std::multiset<Key> accumulator;

    void createAccumulator() override {}

    void accumulate(Key input) override {
        accumulator.insert(input);
    }
    void retract(Key input) override {
        auto it = accumulator.find(input);
        assert(it != accumulator.end());
        accumulator.erase(it);
    }
    std::optional<Key> getValue() override {
        if (accumulator.empty()) {
            return std::nullopt;
        } else {
            return *--accumulator.end();
        }
    }
};