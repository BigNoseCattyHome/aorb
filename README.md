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
├── backend                 # 后端服务代码
│   ├── go-services         # 使用Go编写的服务
│   │   ├── api-gateway     # API网关服务
│   │   ├── auth            # 认证服务
│   │   ├── message         # 消息服务
│   │   ├── poll            # 投票服务
│   │   └── recommendation  # 推荐服务
│   └── java-services       # 使用Java编写的服务
│       └── auth            # 用户服务
├── frontend                # 前端代码
│   ├── android 
│   ├── fonts               # 字体文件
│   ├── images              # 图片资源
│   ├── ios 
│   ├── lib                 # 库文件
│   │   ├── models          # 数据模型
│   │   ├── routes          # 路由
│   │   ├── screens         # 页面
│   │   ├── services        # 服务代码
│   │   └── widgets         # 小部件
│   ├── linux 
│   ├── macos  
│   ├── test 
│   ├── web
│   └── windows 
├── k8s                     # Kubernetes配置文件
├── monitoring              # 监控服务
│   ├── grafana             # Grafana监控配置
│   └── prometheus          # Prometheus监控配置
└── scripts                 # 脚本文件
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

### 后台各个服务的开启
rabbitMQ(MAC):
```shell
brew services start rabbitmq
```
consul(MAC):
```shell
consul agent -dev
```
redis(MAC)
```shell
redis-server
```

### 微服务的启动
执行go-services中每一个模块的main.go文件

## 📝 开发文档

[Flutter开发过程用到组件指南](https://obyi4vacom.feishu.cn/file/E9vdbu0RBocg4yxfV0NcS1kHnwe)

[Git使用指南](http://sirius1y.top/posts/notes/dev/%E6%8C%87%E5%8D%97%E5%9B%A2%E9%98%9Fgit%E5%8D%8F%E4%BD%9C/)

