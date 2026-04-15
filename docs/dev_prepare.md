## 开发环境准备(mac)

### 1.安装 Go 语言：
``brew install go``


### 2.数据库与中间件安装

#### 安装 PostgreSQL：
```shell
brew install postgresql@15

# 验证是否安装成功 方式一
psql --version
# 验证是否安装成功 方式二（通过 brew安装的）
brew list | grep postgresql

#启动
brew services start postgresql@15 

# 查看服务启动状态
brew services list
```
#### 安装 Redis（用于法条缓存）：

```shell
brew install redis
brew services start redis
```

#### 安装 Meilisearch（推荐的搜索引擎，非常轻量）：
```shell
brew install meilisearch
brew services start meilisearch
```

### 3.PostgreSQL 数据库初始化
安装完后，你需要创建一个供开发使用的数据库：

#### 尝试连接数据库
终端输入：``psql postgres``

在 SQL 命令行中输入：
```sql
SQL
CREATE DATABASE law_db;
CREATE USER zhangkx WITH PASSWORD '你的密码';
GRANT ALL PRIVILEGES ON DATABASE law_db TO zhangkx;
\q
```

### 4.GoLand 项目初始化
### go.mod
1. 在 Go 项目中，项目根目录下的 go.mod 文件是通过 go mod init 命令生成的。 // go mod init <模块名称>
它是 Go Modules（Go 的包管理工具）的核心，用于定义模块路径和管理依赖版本。


2. 如何自动填充依赖？
   go mod init 只是创建了文件，并不会自动把代码里 import 的第三方库加进去。你需要接着运行：

```shell
go mod tidy
```
- go mod tidy 的作用：
    - 扫描：它会扫描你项目中所有的 .go 源文件。
    - 添加：发现代码中引用了但 go.mod 里没有的库，会自动下载并写入 go.mod。
    - 移除：发现 go.mod 里有但代码中已经不再使用的库，会自动删除。

生成 Lock 文件：同时会生成或更新一个 go.sum 文件，用于记录依赖库的哈希值，确保版本安全性。

3. 为什么需要这个文件？
   摆脱 GOPATH：在有 go.mod 的情况下，你的项目可以放在电脑的任何位置，而不需要非得放在 ~/go/src 下。

版本控制：它精确记录了项目依赖的库版本，确保在其他机器（如服务器或同事的电脑）上编译出的结果是一致的。


### 5.安装gin
go get -u github.com/gin-gonic/gin
