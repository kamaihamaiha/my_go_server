## start

### 1. 安装依赖

```shell
GOTOOLCHAIN=local go mod tidy
```

### 2. 启动服务

默认使用：
- HTTP 地址：`:8080`
- SQLite 数据库：`db/laws.db`
- 解析后法律 JSON 目录：`data/law_json`

解析后的法律文件默认按分类放在：

`data/law_json/laws_by_type/type_<lawTypeId>/<versionId>.json`

例如：

`data/law_json/laws_by_type/type_230/2c909fdd678bf17901678bf59da8002d.json`

服务会根据法律本身的 `lawTypeId` 去对应分类目录下读取文件。

也可以通过环境变量覆盖：
- `HTTP_ADDR`
- `LAW_DB_PATH`
- `LAW_DETAIL_JSON_DIR`

启动命令：

```shell
GOTOOLCHAIN=local go run .
```

例如：

```shell
HTTP_ADDR=:9090 LAW_DB_PATH=db/laws.db LAW_DETAIL_JSON_DIR=data/law_json GOTOOLCHAIN=local go run .
```

### 3. 测试接口

健康检查：

```shell
curl http://localhost:8080/healthz
```

查看各具体分类的预览列表（每类最多 20 条）：

```shell
curl http://localhost:8080/api/v1/types/previews
```

分页查看某个分类下的法律列表：

```shell
curl "http://localhost:8080/api/v1/types/230/laws?page=1&pageSize=20"
```

读取某个法律的解析后 JSON：

```shell
curl http://localhost:8080/api/v1/laws/2c909fdd678bf17901678bf59da8002d/parsed
```

如果当前机器还没有准备好解析后的 JSON 文件，接口会返回一个“暂无解析数据”的占位 JSON。

目录说明也保存在：

`data/law_json/README.md`
