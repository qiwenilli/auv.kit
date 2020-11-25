# auv.kit
一个微服务框架工具集合，包括微服务框架、服务注册发现、限流自动熔断、路径跟踪


## 功能列表
1. 服务注册、及发现 (仅支持etcd)
2. 集群场景下自动熔断
3. 单实例限流、支持配置ip白名单及可通配符规则（白名单不限制）、支持设置多组url rule配置限流（可用于服务降级）
4. 生成api swagger 文档
5. 支持分布式acid
6. protobuf 定义接口，并生成go server接口定义、及多语言client调用代码
7. 支持自定义 middleware
8. 支持server exit hook
9. 友好的log（logrus）配置，及response callback time
10. 集成 traceId 生成，输出日志自动携带，用于链路跟踪
11. 配置中心化
12. 服务降级，自动降级

## 辅助功能
1. 支持runtime参数 cross domain开关，友好支持前端调试
2. 支持runtime参数 go pprof 是否开启pprof
3. 支持runtime参数 server web监控
4. 支持runtime参数 是否打swagger ui


