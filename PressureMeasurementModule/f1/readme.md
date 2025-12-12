
场景
推荐方式
单机极限 RPS
想保留 f1 复杂逻辑，又要高性能
方式一（f1 里直接用 Vegeta 库）
80k~120k
只需要准备/清理数据，中间纯打流量

方式二（f1 + vegeta 二进制）  50k~80k

要打百万级甚至千万级分布式压测
方式三（f1 导演 + Vegeta 打手）


go run main.go run constant login -r 10/s -d 20s --verbose