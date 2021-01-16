cmake_minimum_required(VERSION 3.17)
project(dddd)

set(CMAKE_CXX_STANDARD 17)

add_executable(tidb_hackathon_2020 src/main.cpp src/raftstore/TxnSolver.cpp src/raftstore/TxnSolver.h src/TobeImplement.cpp src/preagg/AggFunctionHandle.h src/preagg/AggFunctionHandle.cpp)
