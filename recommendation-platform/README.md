# video-system

本仓库实现三大角色并支持独立启动：

1. **Content Platform**（视频接收与存储）
2. **Recommendation Platform**（推荐算法）
3. **Client**（客户端）

每个角色都可以单独启动，但要得到完整端到端体验，建议一起启动。

---

## 1. Content Platform（video-platform）

环境变量：

- `PLATFORM_ID` 平台ID（例如 `platformA`）
- `PLATFORM_PORT` 端口（默认 8080）
- `INDEX_BASE` 索引服务地址（可空；为空时不写索引）
- `ACCEPT_TAGS` 内容规则（逗号分隔，如 `tech,ai`，为空表示不限制）
- `JWT_SECRET` JWT密钥

启动示例（平台A）：

```
cd video-platform
$env:PLATFORM_ID="platformA"
$env:PLATFORM_PORT="8080"
$env:INDEX_BASE="http://localhost:8083"
$env:ACCEPT_TAGS="tech,ai,technology"
go run .
```

启动示例（平台B，全部视频）：

```
cd video-platform
$env:PLATFORM_ID="platformB"
$env:PLATFORM_PORT="8084"
$env:INDEX_BASE="http://localhost:8083"
$env:ACCEPT_TAGS=""
go run .
```

---

## 2. Recommendation Platform

依赖索引服务（video-index），但可以独立运行。

环境变量：

- `RECOMMEND_PORT` 端口（默认 8082）
- `INDEX_BASE` 索引服务地址（默认 `http://localhost:8083`）
- `JWT_SECRET` JWT密钥

启动：

```
cd recommendation-platform
go run .
```

返回格式：

```
GET /recommend
{
  "videos": [
    "video://platformA/xxxx",
    "video://platformB/yyyy"
  ]
}
```

---

## 3. Client（客户端）

启动静态网页服务：

```
cd client
python -m http.server 5173
```

浏览器访问：

- 首页：`http://localhost:5173/index.html`
- 播放页：`http://localhost:5173/player.html?uri=video://platformA/xxxx`

---

## 4. Index（video-index）

索引服务作为推荐平台的数据层。

```
cd video-index
go run .
```

---

## 5. Streaming Service（跨平台播放）

环境变量：

- `STREAM_PORT`（默认 8081）
- `PLATFORM_MAP` 平台路由映射（逗号分隔 `id=url`）
- `P2P_MAP` P2P网关映射（逗号分隔 `id=url`）

示例：

```
cd streaming-service
$env:PLATFORM_MAP="platformA=http://localhost:8080,platformB=http://localhost:8084"
$env:P2P_MAP="platformA=http://localhost:8090,platformB=http://localhost:8091"
go run .
```

---

## 6. P2P Node（可选）

用于 chunk 拉取加速。

```
cd p2p-node
$env:CHUNK_DIR="..\video-platform\data\video-storage\chunks"
$env:P2P_HTTP_PORT="8090"
go run .
```
