/house-edge-sniper
├── /node_listeners           # Mempool 與區塊監聽 (Producers)
│   ├── /ethereum_wss         # 透過 WebSocket 監聽 ETH 節點
│   ├── /solana_rpc           # Solana 事件監聽
│   └── /abi                  # 目標賭場的合約 ABI
├── /kafka_core               # Kafka 延遲最佳化設定
│   └── server.properties     # 調整 batch.size 與 linger.ms
├── /strategy_engine          # 核心邏輯 (Consumers)
│   ├── /vrf_analyzer         # 隨機數延遲預測模組
│   ├── /pool_monitor         # 資金池流動性計算
│   └── /arb_math             # 期望值 (EV) 計算引擎
├── /execution_bot            # 自動發送交易模組
│   ├── /flashbots_relay      # 繞過公開 Mempool 的防夾擊發送
│   └── /wallet_manager       # 私鑰與 Gas 費用管理
└── /contracts                # 鏈上輔助套利合約 (Solidity)