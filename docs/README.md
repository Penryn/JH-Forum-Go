# 精弘论坛

## 项目说明

本论坛使用了开源论坛框架[paopao-ce](https://github.com/rocboss/paopao-ce)。
在此基础上修改为供zjut在校生使用的校园论坛。

## 🏗 快速开始

### 环境要求

* Go (1.20+)
* MySQL (8.0)
* Redis 
* MinIO 
* MeiliSearch (1.4)
* ...

以上环境仅供参考，其他版本的环境未进行充分测试

### 配置说明

`config.yaml.example` 是一份完整的配置文件模版，JH-Forum启动时会读取`./config.yaml`配置文件。

```sh
cp config.yaml.example config.yaml
vim config.yaml # 修改参数
```

配置文件中的 `Features` 小节是声明JH-Forum运行时开启哪些功能项:

```yaml
...

Features:
   Default: [ "MySQL", "Option", "LocalOSS", "LoggerFile" ]

...
```


这里 `Default`套件 代表的意思是： 开启了`MySQL/LocalOSS/LoggerFile` 3项功能。



### 使用说明

克隆代码库

   ```sh
   git clone https://github.com/zjutjh/JH-Forum-Go.git

   ```

#### 项目目录介绍
```go
JH-Forum-Go
├── cmd                // 存放主程序入口文件
├── custom             // 日志存放地址
├── docs               // 文档相关文件
│   ├── database
│   │   └── README.md  // 数据库相关文档
│   └── README.md      // 项目文档
├── internal           // 内部实现代码
│   ├── conf           // 配置相关代码
│   ├── core           // 核心功能代码
│   ├── dao            // 数据访问层（Data Access Object）
│   ├── events         // 事件处理相关代码
│   ├── metrics        // 指标监控相关代码
│   ├── model          // 模型定义
│   ├── servants       // 服务层相关代码
│   ├── service        // 业务逻辑层
│   └── internal.go    // 内部包的入口文件
├── mirc               // Mir 相关代码
│   ├── auto           // 自动生成的代码
│   ├── web            // Web 相关代码
│   └── README.md      // Mir 相关文档
├── pkg                // 第三方库和工具
│   ├── app            // 应用相关代码
│   ├── convert        // 数据转换相关代码
│   ├── debug          // 调试相关代码
│   ├── http           // HTTP 相关工具
│   ├── json           // JSON 处理工具
│   ├── naming         // 命名相关工具
│   ├── obx            // 其他业务扩展（Obx）
│   ├── Tool           // 工具集合
│   ├── sensitive      // 文字敏感信息处理工具
│   ├── types          // 类型定义和处理工具
│   ├── utils          // 通用工具
│   ├── version        // 显示版本相关工具
│   └── xerror         // 错误处理工具
├── release            // 打包文件
├── script             // 脚本文件
│   └── forum-mysql.sql // 初始化数据库的 SQL 脚本
├── .gitignore         // Git 忽略文件
├── .golangci.yml      // GolangCI-Lint 配置文件
├── build-image.sh     // 构建 Docker 镜像的脚本
├── build-release.sh   // 构建发布版本的脚本
├── config.yaml.example// 配置文件示例
├── go.mod             // Go 依赖管理文件
├── go.sum             // Go 依赖版本锁定文件
├── LICENSE            // 项目许可证
├── main.go            // 项目入口文件
├── Makefile           // Makefile 文件
├── README.md          // 项目简介和说明文档

```



#### 后端打包和运行方法

1. 导入项目根目录下的 `scripts/forum-mysql.sql` 文件至MySQL数据库
2. 拷贝项目根目录下 `config.yaml.example` 文件至 `config.yaml`，按照注释完成配置`mysql、redis、meilisearch、minIO`
   
3. 直接运行后端，作为后端测试使用（运行）    
   运行api服务:
```sh
   make run
   ## 或者
   go run main.go
 ```

4. 编译后端获得可执行文件（打包）    
   编译api服务:
    ```sh
   ## 在linux环境下
   make build

   ## 在Windows环境(cmd)下
   SET GOOS=linux
   SET GOARCH=amd64
   make linux-amd64
    ```
   编译后在`release`目录可以找到对应可执行文件。
    ```sh
   ## 在linux环境
   release/JH-Forum
   ## 在Windows环境
   release/linux-amd64/JH-Forum-GO/JH-Forum
    ``` 
5. 项目部署方法
* 本地部署
后端服务使用 systemctl或supervisor 守护进程，并通过 nginx 反向代理后，提供API给前端服务调用。
```
systemctl start jh-forum
```
* docker部署
可以通过脚本打包成镜像再运行该镜像
```sh
## 为脚本添加可执行权限
chmod +x ./build-image.sh
## 构建镜像
./build-image.sh
## 运行容器，直接docker部署至少要更改对应的mysql，redis，minio，meili的配置
mkdir custom &&docker run -d   -p 8008:8008   -v ${PWD}/custom:/app/data/custom -v ${PWD}/config.yaml.example:/app/paopao-ce/config.yaml   --name jh-forum   zjutjh/jh-forum
```
* docker-compose 部署
默认是使用config.yaml.example的配置，如果需要自定义配置，请拷贝默认配置文件(比如config.yaml)，修改后再同步配置到docker-compose.yaml
```
### 给予minio权限，不然minio容器会报错
chmod -R 777 ./custom/data/minio/data
### 启动
docker compose up -d
### 停止并删除容器
docker compose down
```


#### 后端编写接口
1. 本项目采用go-mir自动生成框架，因此先在`mirc/web`下编写对应的接口
2. 然后在命令行下输入下面这串命令
```sh
make generate
```
3. go-mir会在`mirc/auto`自动生成对应的api
4. 然后再到`internal/servants/web`把执行函数编写出来


#### 后端使用工具
1. pprof
   * 介绍
pprof 是 Go 语言的性能分析工具，用于分析和诊断 Go 程序的性能瓶颈。它可以帮助开发者了解程序在运行时的 CPU 和内存使用情况，从而找出性能问题并优化程序。
   * 使用
在配置文件下，将Pprof功能加入进去，可以启动Pprof服务，接着在新建一个终端输入以下命令，可以进入 pprof 的交互式模式进行分析
```sh
## 启动应用后，可以通过访问如下 URL 收集 30 秒的 CPU profile：
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
### 如果你有多个 profile 文件，可以使用 pprof 工具将它们合并成一个 default.pgo 文件：
go tool pprof -proto profile1.pprof profile2.pprof > default.pgo
### 使用 PGO 进行编译,将生成的 default.pgo 文件放在你的主包目录中，然后使用以下命令进行编译：
go build -pgo=auto

```
2. golangci-lint
   * 介绍
golangci-lint 是一个 Go 语言的静态代码分析工具，用于集成多个 Go 语言的代码检查工具，帮助开发者在开发过程中保持代码质量和风格一致性。它汇集了许多流行的 linter（代码检查器），使得代码的静态分析变得更加高效和全面。
   * 使用
需要现在本地安装golangci-lint环境，然后可以调整.golangci.yml配置文件修改检查的内容，再终端输入以下命令执行
```sh
golangci-lint run
```  
3. Pyroscope
   * 介绍
Pyroscope 是一个开源的连续性能分析平台，用于监控和分析应用程序的性能。它提供了实时的 CPU 和内存使用情况剖析，帮助开发者识别和解决性能瓶颈。Pyroscope 支持多种数据采集方式，并且具有友好的 Web 界面来可视化性能数据。
   * 使用
在配置文件的Features中加入Pyroscope，然后在下面的配置文件输入相关配置信息，即可在浏览器输入该url查看该服务的运行情况
```
      http://ip:4040
```