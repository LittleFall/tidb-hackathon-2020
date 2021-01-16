//
// Created by littlefall on 2021/1/15.
//

#pragma once
#include <string>
#include <map>
#include <set>

using Key = uint64_t;

class AggFunctionHandle {

public:

    // Todo: use template parameter to get fast.
    AggFunctionHandle() = default;

    void accumulate(Key input);
    void retract(Key input);

    std::optional<Key> getSum();
    std::optional<Key> getCount();
    std::optional<Key> getMin();
    std::optional<Key> getMax();
    std::optional<double> getAvg();

    std::optional<Key> getSumDistinct();
    std::optional<Key> getCountDistinct();
    std::optional<Key> getMinDistinct();
    std::optional<Key> getMaxDistinct();
    std::optional<double> getAvgDistinct();

private:
    Key sum{}; // Todo: 考虑溢出的情况
    Key distinct_sum{};
    uint64_t count{};

    std::map<Key, uint64_t> values;
};

