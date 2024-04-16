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

## 📋项目结构
```shell
aorb
├── android               # Android平台特定代码目录，包含原生Android项目文件
│   ├── app               # Android应用的主要项目文件夹，包含源码和资源
│   └── gradle            # Gradle配置文件，用于Android项目的构建管理
├── api                   # 后端API目录，包含所有后端逻辑
│   ├── handlers          # 处理不同API请求的处理器
│   ├── middleware        # 中间件，用于请求处理前的拦截处理
│   └── routes            # 路由配置，定义URL路由到处理器的映射
├── assets                # 静态资源文件夹，存放图像、字体等资源
├── cmd                   # 命令行应用入口，用于执行后端服务
├── config                # 配置文件目录，存放应用的配置文件
├── deployments           # 部署相关文件，如Kubernetes配置
│   ├── kubernetes        # Kubernetes部署配置文件
│   └── monitoring        # 监控配置文件，如Prometheus和Grafana配置
├── docs                  # 文档目录，包含项目文档
├── ios                   # iOS平台特定代码目录
│   ├── Flutter           # Flutter模块和配置
│   ├── Runner            # iOS主要的运行项目
│   ├── RunnerTests       # iOS的测试代码
│   ├── Runner.xcodeproj  # Xcode项目文件
│   └── Runner.xcworkspace # Xcode工作区文件
├── lib                   # Flutter的Dart代码库
│   ├── controllers       # 控制器，负责业务逻辑的实现
│   ├── models            # 数据模型定义
│   ├── services          # 服务层，通常用于网络请求等异步操作
│   └── views             # 视图层，所有的UI组件
├── linux                 # Linux平台特定代码目录
│   └── flutter           # Linux上的Flutter配置和文件
├── macos                 # macOS平台特定代码目录
│   ├── Flutter           # 包含Flutter相关的配置和资源
│   ├── Runner            # macOS应用的主体代码
│   ├── RunnerTests       # macOS应用的测试代码
│   ├── Runner.xcodeproj  # macOS Xcode项目文件
│   └── Runner.xcworkspace # macOS Xcode工作区文件
├── pkg                   # Go语言或其他后端语言的包目录
│   ├── store             # 数据存储层，处理数据持久化
│   └── util              # 工具包，存放通用工具和助手函数
├── test                  # 测试目录，存放所有的测试代码
├── web                   # Web平台特定代码目录
│   └── icons             # Web用图标
└── windows               # Windows平台特定代码目录
    ├── flutter           # 包含Flutter的配置和文件
    └── runner           # Windows应用的主体代码
```

