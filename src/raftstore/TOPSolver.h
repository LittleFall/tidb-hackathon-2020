// TiCDC Open Protocol solver
// Created by littlefall on 2021/1/16.
// 聚合的中间一层，在此处维护时间版本。

#pragma once

#include <cstdint>
#include "../preagg/AggFunctionHandle.h"
#include "TxnSolver.h"

class TOPSolver {



    void createHandle() {

    }

    // 积攒一个 insert event
    void rowChangeInsert(uint64_t column_id, Key key, TimeStamp) {

    }

    // 积攒一个 delete event
    void rowChangeDelete() {

    }

    // resolve，实际处理当前时间点
    void resolve() {

    }

};
