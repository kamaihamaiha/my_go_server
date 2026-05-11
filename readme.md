## law 服务端

总结方案清单：
语言：Go 1.22+

Web 框架：Gin

ORM：GORM

关系型数据库：PostgreSQL 15+

搜索引擎：Meilisearch (强烈推荐)

缓存：Redis

部署：Docker + Docker Compose (一键部署所有组件，方便迁移服务器)

## 当前基础框架

当前仓库已经先落地了一个基于 SQLite 的只读 HTTP 服务基础版本，用来承接最初版法律数据。

已提供接口：
- `GET /healthz`
- `GET /api/v1/types/previews`
- `GET /api/v1/types/:typeId/laws?page=1&pageSize=20`
- `GET /api/v1/laws/:versionId/parsed`

相关文档：
- 启动说明：`docs/start.md`
- 实现计划：`docs/http-service-plan.md`


---

### 当前实现功能
现在 main.go 会启动一个基于 Gin 的只读法律服务，配置、SQLite 连接、仓储/服务/路由都拆开了，核心入口在 main.go ，接口实现主要在 internal/handler/law_handler.go、internal/service/law_service.go 和 internal/server/router.go。

已经提供这 3 个能力：

GET /api/v1/types/previews：只返回有直接法律数据的具体分类，每类最多 20 条。
GET /api/v1/types/:typeId/laws：支持 page 和 pageSize 分页，默认按 effectDate DESC，空生效日期排后，再按 publishDate DESC。
GET /api/v1/laws/:versionId/parsed：预留本地 JSON 文件读取；文件不存在时返回“暂无解析数据”的占位 JSON。
文档和测试也补上了：计划文件在 docs/http-service-plan.md，启动说明在 docs/start.md，集成测试在 internal/server/router_test.go。我已经跑通了 GOTOOLCHAIN=local GOPROXY=https://goproxy.cn,direct go test
