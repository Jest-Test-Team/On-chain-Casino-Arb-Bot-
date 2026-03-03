//! WebSocket 監聽 Ethereum (或可擴展 Solana) Mempool / 事件
//! 篩選與賭場合約相關的 pending 交易，供 EV 計算與套利判斷

use crate::ev_calculator::ArbOpportunity;
use crate::kafka_producer::KafkaProducer;
use alloy_primitives::TxHash;
use alloy_provider::ProviderBuilder;
use alloy_rpc_types::Filter;
use futures_util::{SinkExt, StreamExt};
use std::sync::Arc;
use tracing::{debug, info};

/// 已知賭場合約位址（可從環境變數或設定檔載入）
fn casino_contract_addresses() -> Vec<alloy_primitives::Address> {
    std::env::var("CASINO_CONTRACTS")
        .ok()
        .map(|s| {
            s.split(',')
                .filter_map(|a| a.trim().parse().ok())
                .collect()
        })
        .unwrap_or_default()
}

/// 執行 Ethereum Mempool / 新區塊事件監聽，發現機會時寫入 Kafka
pub async fn run_ethereum_listener(
    ws_url: String,
    producer: Arc<KafkaProducer>,
) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let provider = ProviderBuilder::new().on_ws(ws_url.parse()?).await?;
    let casinos = casino_contract_addresses();
    if casinos.is_empty() {
        info!("No CASINO_CONTRACTS set; listening to all pending txs (demo mode)");
    }

    // 訂閱 new pending transactions
    let mut sub = provider.subscribe_pending_txs().await?;
    info!("Subscribed to pending transactions");

    while let Some(hash) = sub.next().await {
        let tx_hash = TxHash::from(hash);
        debug!("Pending tx: {:?}", tx_hash);

        // 簡化：不在此解碼 calldata，僅示範流程；實際應取得 tx 詳情並篩選 to in casinos
        if casinos.is_empty() {
            // Demo: 對每筆 pending 做一次「模擬機會」寫入 Kafka
            let opp = ArbOpportunity {
                chain: "ethereum".into(),
                tx_hash: format!("{:?}", tx_hash),
                contract_address: None,
                expected_value_wei: None,
                timestamp_ms: std::time::SystemTime::now()
                    .duration_since(std::time::UNIX_EPOCH)
                    .unwrap()
                    .as_millis() as i64,
            };
            if let Err(e) = producer.send_opportunity(&opp).await {
                tracing::error!("Kafka send error: {}", e);
            }
        }
    }

    Ok(())
}
