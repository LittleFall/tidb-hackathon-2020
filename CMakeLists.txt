cmake_minimum_required(VERSION 3.17)
project(dddd)

set(CMAKE_CXX_STANDARD 17)

add_executable(tidb_hackathon_2020 src/main.cpp src/preagg/TxnSolver.cpp src/preagg/TxnSolver.h src/TobeImplement.cpp)
