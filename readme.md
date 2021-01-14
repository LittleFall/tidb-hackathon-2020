# DDDD

## 编程

语言：C++ 17

环境：Clion

语言指导：[Cpp reference](https://zh.cppreference.com/w/%E9%A6%96%E9%A1%B5)

代码规范：[ClickHouse 编码建议](https://clickhouse.tech/docs/zh/development/style/)

编码流程：在[主仓库](https://github.com/LittleFall/tidb-hackathon-2020)上创建 branch 或 fork，通过 pr 合并，至少需 1 人 review。请勿直接提交到主分支（main）。 

## 项目结构

- src/ 源码目录
    - preagg/ 与预聚合有关的实现
        - TxnSolver 解决单 Region 的事务问题
        
## Todo queue

产品侧
- 实现 apply 逻辑，包括常用的聚合函数
- 实现预聚合的多版本控制
- 实现 learner
- 实现 schema 相关的逻辑
- 完成 tidb 的查询导引
- 启动时做 snapshot

项目结构侧
- 整理 cmake，实现 include 尖括号
- 添加一些基础支持，如日志，单元测试，format，CI


## 实现原理

### 事务处理

处理单 region 的事务。

数据结构：
- 一个堆，存放所有锁的 key 和 ts，堆顶的 ts 最小
- 一个哈希表，存放所有的 commit 事件，key 为键，ts 为值.

算法：
- 每当一个 prewrite 事件到达时，将锁放入堆中。
- 每当一个 commit 事件到达时，将其放入哈希表中。
    然后从堆顶开始循环尝试，如果一个锁已经被解开了，就 apply 它的事件。
    直到碰到第一个未解开的锁，或者所有锁都被解开了。