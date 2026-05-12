# Parsed Law JSON Directory

将解析后的法律 JSON 文件放到当前目录下。

命名规则：

`laws_by_type/type_<lawTypeId>/<versionId>.json`

例如：

`laws_by_type/type_230/2c909fdd678bf17901678bf59da8002d.json`

服务会先按法律的 `lawTypeId` 到对应分类目录查找文件。
为了兼容旧数据，也仍然支持历史平铺路径：

`<versionId>.json`

默认情况下，服务会从 `LAW_DETAIL_JSON_DIR` 读取文件；
如果未显式配置，该目录默认就是 `data/law_json`。
