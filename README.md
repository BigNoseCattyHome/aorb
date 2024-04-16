# 👐 welcome to aorb

## 💖简介

- 用户可以进行“帮我选”和“帮他选”
- “帮我选”中可使用机器自动帮忙选择和社区帮忙选择
- “帮他选”中所有用户都可以帮忙进行选择，也可以进行留言
- 用户可以关注其他用户，当其他用户发布“帮我选”的时候系统会将消息推送给关注者



## 💪技术栈

- 前端Flutter
- 后端Redis+Gin+MongoDB
- 使用Kubernetes+Promethus+Grafana进行拓展和监控
- 文档注释用swag

## 项目结构
aorb
│
├── /api            # 后端API
│   ├── /handlers   # 请求处理
│   ├── /middleware # 中间件
│   └── /routes     # 路由
│
├── /cmd            # 后端入口
│   └── main.go
│
├── /lib
│   ├── /models     # 数据模型
│   ├── /services   # 服务层
│   ├── /controllers# 控制器
│   └── /views      # 视图层
│
├── /config         # 配置文件
│
├── /pkg
│   ├── /store      # 数据访问层
│   └── /util       # 工具类
│
├── /assets         # 静态资源
│
├── /docs           # 文档
│
├── /tests          # 测试代码
│
└── /deployments    # Kubernetes和监控配置
    ├── /kubernetes
    └── /monitoring