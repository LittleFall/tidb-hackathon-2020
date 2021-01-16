cmake_minimum_required(VERSION 3.17)
project(dddd)

set(CMAKE_CXX_STANDARD 17)

add_executable(tidb_hackathon_2020 src/main.cpp src/raftstore/TxnSolver.cpp src/raftstore/TxnSolver.h src/TobeImplement.cpp src/preagg/AggFunctionHandleBase.h src/preagg/CountAggFunctionHandle.cpp src/preagg/SumAggFunctionHandle.cpp src/preagg/MaxAggFunctionHandle.cpp src/preagg/MinAggFucntionHandle.cpp src/preagg/AvgAggFunctionHanle.cpp)
