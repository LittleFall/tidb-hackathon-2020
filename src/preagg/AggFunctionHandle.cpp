//
// Created by littlefall on 2021/1/16.
//

#include "AggFunctionHandle.h"

// Todo: too ugly, need repair

void AggFunctionHandle::accumulate(Key input) {
    sum += input;
    ++count;
    ++values[input];

    if(values[input]==1) {
        distinct_sum += input;
    }
}

void AggFunctionHandle::retract(Key input) {
    sum -= input;
    --count;
    --values[input];

    if(values[input]==0) {
        distinct_sum -= input;
        values.erase(input);
    }
}

std::optional<Key> AggFunctionHandle::getSum() {
    if (count) {
        return std::nullopt;
    } else {
        return sum;
    }
}

std::optional<Key> AggFunctionHandle::getCount() {
    if (count) {
        return std::nullopt;
    } else {
        return count;
    }
}

std::optional<Key> AggFunctionHandle::getMin() {
    if (count) {
        return std::nullopt;
    } else {
        return values.begin()->first;
    }
}

std::optional<Key> AggFunctionHandle::getMax() {
    if (count) {
        return std::nullopt;
    } else {
        return (--values.end())->first;
    }
}

std::optional<double> AggFunctionHandle::getAvg() {
    if (count) {
        return std::nullopt;
    } else {
        return double(sum)/count;
    }
}

std::optional<Key> AggFunctionHandle::getSumDistinct() {
    if (count) {
        return std::nullopt;
    } else {
        return distinct_sum;
    }
}

std::optional<Key> AggFunctionHandle::getCountDistinct() {
    if (count) {
        return std::nullopt;
    } else {
        return values.size();
    }
}

std::optional<Key> AggFunctionHandle::getMinDistinct() {
    if (count) {
        return std::nullopt;
    } else {
        return values.begin()->first;
    }
}

std::optional<Key> AggFunctionHandle::getMaxDistinct() {
    if (count) {
        return std::nullopt;
    } else {
        return (--values.end())->first;
    }
}
std::optional<double> AggFunctionHandle::getAvgDistinct() {
    if (count) {
        return std::nullopt;
    } else {
        return double(distinct_sum)/count;
    }
}