#include <iostream>
#include "preagg/TxnSolver.h"

void testTxnSolver()
{
    TxnSolver txn_solver;

    txn_solver.prewrite("k1", 1, "k1++");

    txn_solver.commit("k1", 6);

    txn_solver.prewrite("k2", 3, "k2--");

    txn_solver.commit("k2", 7);
}


int main()
{

    std::cout << "hello the crawl world." << std::endl;

    testTxnSolver();

    std::cout << "indeed." << std::endl;

    return 0;
}
