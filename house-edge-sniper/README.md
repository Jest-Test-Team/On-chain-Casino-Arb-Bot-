# House Edge Sniper

依 `structure.md` 實作之 On-chain Casino 套利監聽與協調架構。

## 架構

| 目錄 | 說明 |
|------|------|
| **rust-sniper** | 前線：Mempool 監聽、EV 計算、Flashbots 發送、Kafka 生產者 |
| **go-coordinator** | 中樞：Kafka 消費、風控、WebSocket 推送、REST API |
| **next-dashboard** | 戰情面板：即時圖表、設定頁 |

## 快速啟動

### 1. 基礎設施（Kafka、Zookeeper、Redis）

```bash
cd house-edge-sniper
docker-compose up -d
```

### 2. Go 協調器

```bash
cd go-coordinator
go build -o server ./cmd/server
./server
# 預設 HTTP :8080，WebSocket /ws，健康檢查 /health
```

環境變數（選填）：`KAFKA_BROKERS=localhost:9092`、`REDIS_ADDR=localhost:6379`

### 3. Next 戰情面板

```bash
cd next-dashboard
npm install && npm run dev
# 開啟 http://localhost:3000，根路徑會導向 /dashboard
```

設定 `NEXT_PUBLIC_API_URL=http://localhost:8080`、`NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws`（預設即為此）

### 4. Rust Sniper（可選）

- 需 **cmake**（或系統安裝 librdkafka 後改 Cargo.toml 為系統連結）
- 需 Kafka、Redis 已啟動

```bash
cd rust-sniper
cargo build --release
# 設定 ETH_WS_URL（Ethereum WebSocket）、KAFKA_BROKERS、可選 CASINO_CONTRACTS
./target/release/house-edge-sniper
```

## API 摘要

- `GET /health` — 健康檢查
- `GET /api/risk` — 最大回撤等風控參數
- `POST /api/settings` — 寫入錢包、最低餘額、最大回撤（存 Redis）
- `GET /ws` — WebSocket，推送 `{ type: "opportunity", payload: {...} }`

## Kafka Topic

- `house-edge.opportunities` — Rust Sniper 寫入套利機會，Go Coordinator 消費

## 實作狀態

- Docker Compose、Go Coordinator（Kafka 消費、風控、WS、API）、Next 戰情與設定頁已可運行。
- Rust Sniper 程式結構與依賴已就緒，需本機安裝 cmake（或 librdkafka）後編譯；Alloy 訂閱 pending tx 依版本可能需微調 API。
