# 🎰 House Edge Sniper (On-chain Casino Arb-Bot)

## 專案簡介
針對 Web3 去中心化賭場（DEX Casinos）的套利與監控機器人。透過直接監聽區塊鏈節點的 Mempool 與 Event Logs，利用 Kafka 的高吞吐量特性，尋找賠率錯誤、莊家資金池枯竭或隨機數（VRF）延遲的漏洞，並自動發起搶跑（Front-running）交易。

## 核心技術棧
* **區塊鏈互動:** Node.js (Ethers.js), Rust (Alloy/Ethers-rs) -> 追求極致速度推薦用 Rust
* **訊息串流:** Apache Kafka (低延遲調優)
* **記憶體緩存:** Redis (快速比對資金池狀態)
* **智能合約:** Solidity (執行套利的 Flashbots/MEV 合約)