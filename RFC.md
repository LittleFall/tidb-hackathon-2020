<!--
This is a template for TiDB's change proposal process, documented [here](./README.md).
-->

# Proposal: 基于 raft log 实现 TiDB 的物化视图

- Author(s):     齐智，耿立琪，胡志峰，刘继聪<!-- Author Name, Co-Author Name, with the link(s) of the GitHub profile page -->
- Last updated:  2021.01.11<!-- Date -->
- Discussion at: this repo.

## Abstract

本项目的主要目标是实现基于 raft log 实现 TiDB 的物化视图。

<!--
A short summary of the proposal:
- What is the issue that the proposal aims to solve?
- What needs to be done in this proposal?
- What is the impact of this proposal?
-->

## Background

实际业务场景中经常会频繁执行一些相似的查询，比如对一张大表每隔几秒进行一次聚合，以此来生成报表，这会给 TiDB 造成较大的压力。

两次执行之间表中的数据往往变化不大，如果我们能将聚合结果保存下来，再在每次有新数据到来时只对需要的结果进行更新，就可以免去很多重复的计算，这就是预聚合的概念。

利用预聚合的方式，可以极大的减轻类似场景下的查询压力。

物化视图是查询结果的本地存储，预聚合就是物化视图的功能之一，除了预聚合之外，物化视图还包括增量 Join 等功能，都可以用来提升查询的效率。

<!--
An introduction of the necessary background and the problem being solved by the proposed change:
- The drawback of the current feature and the corresponding use case
- The expected outcome of this proposal.
-->

## Proposal

本项目的主要目标是实现基于 raft log 实现 TiDB 的物化视图。

通过实现类似于 TiFlash Proxy 的 raft learner，将导出的 raft log 进行预聚合等处理，以此来在 TiDB 获得中类似于物化视图的效果。

这里的预聚合结果带有 schema，并且可以满足事务（多版本控制）以及最终一致性。

除了预聚合之外，还可以通过接入 Flink 等方式来实现更多的功能，如流式 Join。


<!--
A precise statement of the proposed change:
- The new named concepts and a set of metrics to be collected in this proposal (if applicable)
- The overview of the design.
- How it works?
- What needs to be changed to implement this design?
- What may be positively influenced by the proposed change?
- What may be negatively impacted by the proposed change?
-->

## Rationale

考虑在预聚合的过程中时间窗口的存在，为了实现最终一致性，可以考虑实现类似 stale read 的方式，每次读取一个固定时间点（如 10 秒前）的数据。

<!--
A discussion of alternate approaches and the trade-offs, advantages, and disadvantages of the specified approach:
- How other systems solve the same issue?
- What other designs have been considered and what are their disadvantages?
- What is the advantage of this design compared with other designs?
- What is the disadvantage of this design?
- What is the impact of not doing this?
-->

## Compatibility and Migration Plan

不会导致兼容性问题。
<!--
A discussion of the change with regard to the compatibility issues:
- Does this proposal make TiDB not compatible with the old versions?
- Does this proposal make TiDB not compatible with TiDB tools?
    + [BR](https://github.com/pingcap/br)
    + [DM](https://github.com/pingcap/dm)
    + [Dumpling](https://github.com/pingcap/dumpling)
    + [TiCDC](https://github.com/pingcap/ticdc)
    + [TiDB Binlog](https://github.com/pingcap/tidb-binlog)
    + [TiDB Lightning](https://github.com/pingcap/tidb-lightning)
- If the existing behavior will be changed, how will we phase out the older behavior?
- Does this proposal make TiDB more compatible with MySQL?
- What is the impact(if any) on the data migration:
    + from MySQL to TiDB
    + from TiDB to MySQL
    + from old TiDB cluster to new TiDB cluster
-->

## Implementation

1. 实现一个 raft learner，可以接受 raft log 并按方便处理的形式导出到一个模块中（称为“状态机”）。
2. 让状态机能够接收 TiDB 的 schema.
3. 在状态机中实现预聚合相关的逻辑。
4. (optional) 在状态机中实现增量 Join 等其它物化视图的相关逻辑。

<!--
A detailed description for each step in the implementation:
- Does any former steps block this step?
- Who will do it?
- When to do it?
- How long it takes to accomplish it?
-->

## Testing Plan

以一组标准查询（如 TPC-C, TPC-H 等），按使用物化视图与不使用物化视图的方式分别查询多次，比对结果正确性以及效率。

<!--
A brief description on how the implementation will be tested. Both integration test and unit test should consider the following things:
- How to ensure that the implementation works as expected?
- How will we know nothing broke?
-->

## Open issues (if applicable)

暂无。

<!--
A discussion of issues relating to this proposal for which the author does not know the solution. This section may be omitted if there are none.
-->
