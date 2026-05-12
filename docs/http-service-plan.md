# Go 法律 HTTP 服务基础框架

## 摘要
- 将当前只有 `/ping` 的 Gin 程序扩成一个可启动的只读 HTTP 服务，首版只接现有 SQLite，不接 PostgreSQL、Redis、Meilisearch。
- 数据源保持 `db/laws.db` 原样，不做迁移或 `AutoMigrate`；应用以只读方式连接 SQLite，避免误写和锁冲突。
- `go.mod` 调整为 `go 1.26.1`，与当前机器本地工具链一致，先解决“项目能直接跑起来”的问题。

## 对外接口
- `GET /healthz`
  - 返回服务存活状态。
- `GET /api/v1/types/previews`
  - 返回“具体分类”预览列表，只包含 `types` 中直接有法律数据的分类。
  - 每个分类最多 20 条，分类顺序按 `types.id ASC`。
  - 每条法律返回 `versionId`、`title`、`lawTypeId`、`lawType`、`publishDate`、`effectDate`、`effectiveStatus`、`authorityName`。
- `GET /api/v1/types/:typeId/laws?page=1&pageSize=20`
  - `typeId` 使用 `types.id`。
  - 返回 `page`、`pageSize`、`total`、`totalPages`、`items`。
  - `page` 默认 1，`pageSize` 默认 20，最大 100。
- `GET /api/v1/laws/:versionId/parsed`
  - 从本地目录读取解析后的 JSON 文件，目录由 `LAW_DETAIL_JSON_DIR` 指定，默认 `data/law_json`。
  - 默认按分类目录读取：`laws_by_type/type_<lawTypeId>/<versionId>.json`。
  - 兼容旧路径：`<versionId>.json`。
  - 命中时返回统一响应包，`data.available=true`，`data.content` 为原始解析 JSON（`json.RawMessage`）。
  - 文件缺失或目录不存在时返回 HTTP 200，`data.available=false`，`data.content=null`，`message` 为“暂无解析数据”。
  - `versionId` 在 `laws_list` 中不存在时返回 HTTP 404。

## 实现变更
- 启动与配置
  - 将启动逻辑拆成配置加载、数据库初始化、路由注册三段。
  - 支持 `HTTP_ADDR`、`LAW_DB_PATH`、`LAW_DETAIL_JSON_DIR`；默认分别为 `:8080`、`db/laws.db`、`data/law_json`。
- 数据访问与分层
  - 引入 `GORM + sqlite driver`，仅建读模型 `Type` 与 `LawList`，不执行自动建表或结构变更。
  - 使用 `handler -> service -> repository` 三层，便于后续把 SQLite 替换为 PostgreSQL。
  - SQLite 连接使用只读 DSN，并设置 busy timeout；列表查询统一按 `effectDate DESC` 排序，`publishDate DESC` 作为次级排序，空日期放最后。
- 业务规则
  - 接口 1 只返回具体分类，不返回 `101/102/222` 这类汇总节点。
  - 接口 1 先查出有直接数据的分类，再按分类逐个取前 20 条；当前分类数不多，首版接受这种简单实现，不引入复杂窗口函数。
  - 接口 2 在分页前先验证 `typeId` 是否存在。
  - 接口 3 不读取 `laws_list.detailJson`，只按本地文件目录预留能力。
- 响应与文档
  - 统一响应结构：`code`、`message`、`data`。
  - 更新启动文档，补充运行方式、环境变量和 3 个 `curl` 示例。

## 测试
- HTTP 集成测试使用 `httptest`。
- 为列表接口准备小型临时 SQLite 测试库，覆盖：
  - 预览接口只返回具体分类，且每类最多 20 条。
  - 指定分类分页正确，`page/pageSize/total/totalPages` 正确。
  - 排序按 `effectDate DESC`，缺失 `effectDate` 时排后，并由 `publishDate DESC` 决定先后。
- 为详情接口准备临时目录，覆盖：
  - 文件存在时返回 `available=true` 和原始 JSON。
  - 文件缺失或目录不存在时返回占位 JSON。
  - `versionId` 不存在时返回 404。
- 最终验收以 `go test ./...` 和本地启动后 `curl` 三个接口为准。

## 假设与默认值
- 首版目标是“先把 HTTP 服务跑起来并接上现有 SQLite 数据”，不实现 README 里提到的 PostgreSQL、Redis、Meilisearch。
- 分类参数统一使用 `types.id`，不支持按分类名检索。
- 不修改 `db/laws.db` 表结构，也不新增索引；如果后续查询性能不够，再单独评估索引或迁移。
- 解析后 JSON 文件的目录本轮只预留配置和读取逻辑，仓库里可以暂时没有实际文件。
