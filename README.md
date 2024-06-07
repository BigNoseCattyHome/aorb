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

推荐使用vscode进行开发，安装flutter插件，以及dart插件

### 将项目克隆到本地

```shell
git clone https://github.com/BigNoseCattyHome/aorb.git
```

### 前端开发 
开发和测试flutter应用

```shell
make run_frontend
```

flutter会自动编译fronted/lib/main.dart文件并运行，选择一个合适的平台进行查看就好，不同平台需要满足特定的工具包。


figma原型设计共享链接：[Aorb原型设计](https://www.figma.com/design/roDqwgrlbQo29vpSqeCVFw/Aorb?node-id=0-1&t=SOBamnPsEXegjKDF-1)

### 数据库初始化

这里是一篇[MongoDB安装和简单上手](https://obyi4vacom.feishu.cn/file/DTTWb1DMjoGynkxmgOBc0qgInWd)文档，可以参考一下

确保在本地安装好MongoDB后，进行数据库初始化：

```shell    
mongosh
```

进入到mongodb shell之后输入命令：
```shell
load("scripts/init_db.js")
```

## 📝 开发文档

[Flutter开发过程用到组件指南](https://obyi4vacom.feishu.cn/file/E9vdbu0RBocg4yxfV0NcS1kHnwe)

[Git使用指南](http://sirius1y.top/posts/notes/dev/%E6%8C%87%E5%8D%97%E5%9B%A2%E9%98%9Fgit%E5%8D%8F%E4%BD%9C/)

