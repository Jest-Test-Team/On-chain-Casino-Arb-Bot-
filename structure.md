/house-edge-sniper
├── /rust-sniper              # 前線監聽與執行 (Producer & Executor)
│   ├── Cargo.toml
│   ├── /src
│   │   ├── mempool_listener.rs # WebSocket 監聽 Ethereum/Solana
│   │   ├── ev_calculator.rs    # 期望值與合約邏輯分析
│   │   ├── flashbots_relay.rs  # MEV 交易發送
│   │   └── kafka_producer.rs   # 將發現的機會打入 Kafka
├── /go-coordinator           # 訊息中樞與 API (Consumer & Backend)
│   ├── go.mod
│   ├── /cmd/server/main.go
│   ├── /internal
│   │   ├── /kafka_consumer     # 訂閱套利訊號與執行結果
│   │   ├── /risk_manager       # 控管錢包餘額與最大回撤
│   │   └── /ws_hub             # 推送實時數據給前端
├── /next-dashboard           # 戰情面板 (Frontend)
│   ├── package.json
│   ├── /app
│   │   ├── /dashboard          # 實時圖表與數據
│   │   └── /settings           # 錢包與策略參數設定
│   └── /components/charts    
└── docker-compose.yml        # 一鍵啟動 Kafka, Zookeeper, Redis