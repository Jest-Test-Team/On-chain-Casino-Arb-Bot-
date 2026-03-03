//! House Edge Sniper - 前線監聽與執行
//! 監聽 Mempool → 計算 EV → 可選 Flashbots 發送 → 機會打入 Kafka

mod ev_calculator;
mod flashbots_relay;
mod kafka_producer;
mod mempool_listener;

use std::sync::Arc;
use tracing::info;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    dotenvy::dotenv().ok();
    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::from_default_env().add_directive("house_edge_sniper=info".parse()?),
        )
        .init();

    info!("House Edge Sniper starting");

    let kafka_brokers = std::env::var("KAFKA_BROKERS").unwrap_or_else(|_| "localhost:9092".into());
    let producer = Arc::new(kafka_producer::KafkaProducer::new(&kafka_brokers)?);

    // 可選：Ethereum WS 端點，用於 Mempool 監聽
    let eth_ws = std::env::var("ETH_WS_URL").ok();
    if let Some(ws_url) = eth_ws {
        let producer_clone = Arc::clone(&producer);
        tokio::spawn(async move {
            if let Err(e) = mempool_listener::run_ethereum_listener(ws_url, producer_clone).await {
                tracing::error!("Ethereum mempool listener error: {}", e);
            }
        });
    }

    // 保持主程式運行；實際可改為從 Kafka 消費「執行指令」並呼叫 flashbots_relay
    tokio::signal::ctrl_c().await?;
    info!("Shutting down");
    Ok(())
}
