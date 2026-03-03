# House Edge Sniper 實作計畫（已實作）

依據專案根目錄 `structure.md` 與 README 之架構，完成以下項目：

## 已完成

1. **docker-compose.yml**  
   - 一鍵啟動 Kafka、Zookeeper、Redis（`house-edge-sniper/docker-compose.yml`）

2. **rust-sniper**（前線監聽與執行）  
   - `Cargo.toml`：依賴 alloy、rdkafka、tokio、reqwest、serde 等  
   - `src/mempool_listener.rs`：WebSocket 監聽 Ethereum pending tx，機會寫入 Kafka  
   - `src/ev_calculator.rs`：ArbOpportunity 結構與 EV 計算介面  
   - `src/flashbots_relay.rs`：MEV bundle 提交 Flashbots Relay  
   - `src/kafka_producer.rs`：寫入 `house-edge.opportunities`  
   - `src/main.rs`：組裝並啟動 listener（需 `ETH_WS_URL`、`KAFKA_BROKERS`）  
   - 編譯需 cmake（或系統 librdkafka）

3. **go-coordinator**（訊息中樞與 API）  
   - `go.mod`：Gin、kafka-go、redis、gorilla/websocket、zap  
   - `internal/kafka_consumer`：訂閱 `house-edge.opportunities`，回呼 Handler  
   - `internal/risk_manager`：餘額/最大回撤檢查，Redis 讀寫  
   - `internal/ws_hub`：WebSocket Hub 廣播  
   - `cmd/server/main.go`：Gin HTTP、/ws、/health、/api/risk、POST /api/settings，Kafka 消費後經風控再廣播

4. **next-dashboard**（戰情面板）  
   - Next.js 14 App Router、TypeScript  
   - `app/dashboard`：即時機會表、機會數量趨勢圖、WebSocket 連線狀態  
   - `app/settings`：錢包、最低餘額、最大回撤設定表單  
   - `components/charts/OpportunityChart.tsx`：簡易長條趨勢圖  
   - 深色主題、導覽列

## 使用方式

見 `house-edge-sniper/README.md`：先 `docker-compose up -d`，再啟動 go-coordinator 與 next-dashboard；Rust Sniper 可選，需 cmake 或系統 librdkafka。
