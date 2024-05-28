# 👐 welcome to aorb

## 💖 简介

- 用户可以进行“帮我选”和“帮他选”
- “帮我选”中可使用机器自动帮忙选择和社区帮忙选择
- “帮他选”中所有用户都可以帮忙进行选择，也可以进行留言
- 用户可以关注其他用户，当其他用户发布“帮我选”的时候系统会将消息推送给关注者

## 💪 技术栈

- 前端Flutter
- 后端Redis+Gin+MongoDB
- 使用Kubernetes+Promethus+Grafana进行拓展和监控
- 文档注释用swag

## 📋 项目结构
```shell
aorb
├── backend                    # 后端目录
│   ├── api                    # API接口处理模块，负责前后端通信
│   ├── docs                   # 文档目录，包含开发文档和API文档等
│   ├── models                 # 数据模型目录，定义数据结构和数据库交互
│   └── services               # 服务目录，实现业务逻辑的核心代码
├── frontend                   # 前端目录
│   ├── android                # Android平台的前端代码
│   ├── ios                    # iOS平台的前端代码
│   ├── lib                    # Dart/Flutter库文件，通常存放业务逻辑和UI组件
│   ├── linux                  # Linux平台的前端代码
│   ├── macos                  # macOS平台的前端代码
│   ├── test                   # 测试目录，包含前端的测试代码
│   ├── web                    # Web平台的前端代码
│   └── windows                # Windows平台的前端代码
├── k8s                        # Kubernetes部署和配置文件，用于容器化部署
├── monitoring                 # 监控配置目录
│   ├── grafana                # Grafana监控面板配置文件
│   └── prometheus             # Prometheus监控配置文件
└── scripts                    # 脚本目录，包含自动化脚本等
```


## 🚀 快速开始

进入到`frontend/`下，执行以下命令：

1.  获取依赖

```shell
flutter pub get
```

2.  启动前端
```shell
flutter run
```
flutter会自动编译fronted/lib/main.dart文件并运行

figma原型设计共享链接：[Aorb原型设计](https://www.figma.com/design/roDqwgrlbQo29vpSqeCVFw/Aorb?node-id=0-1&t=SOBamnPsEXegjKDF-1)

## 📝 开发文档

[Flutter开发过程用到组件指南](https://bignosecat.yuque.com/org-wiki-bignosecat-an84wk/ww9doa/duovwzv5mrr4nwuk?singleDoc#《Flutter开发用到的组件和相关知识》)

[Git使用指南](http://sirius1y.top/posts/notes/dev/%E6%8C%87%E5%8D%97%E5%9B%A2%E9%98%9Fgit%E5%8D%8F%E4%BD%9C/)