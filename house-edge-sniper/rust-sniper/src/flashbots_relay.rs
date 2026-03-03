//! MEV / Flashbots 交易發送
//! 將套利 bundle 提交至 Flashbots Relay，避免被搶跑並保護隱私

use alloy_primitives::Bytes;
use reqwest::Client;
use serde::{Deserialize, Serialize};
use tracing::info;

const FLASHBOTS_RELAY_ETH: &str = "https://relay.flashbots.net";

/// Flashbots bundle 單筆 tx 表示（簡化）
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BundleTx {
    pub signed_tx: String, // 0x 開頭 hex
}

/// 提交 bundle 請求
#[derive(Debug, Serialize)]
pub struct FlashbotsBundleRequest {
    pub jsonrpc: String,
    pub id: u64,
    pub method: String,
    pub params: Vec<FlashbotsBundleParams>,
}

#[derive(Debug, Serialize)]
pub struct FlashbotsBundleParams {
    pub txs: Vec<String>,
    pub block_number: String,
    pub min_timestamp: Option<u64>,
    pub max_timestamp: Option<u64>,
}

#[derive(Debug, Deserialize)]
pub struct FlashbotsBundleResponse {
    pub jsonrpc: String,
    pub id: u64,
    pub result: Option<BundleResult>,
    pub error: Option<JsonRpcError>,
}

#[derive(Debug, Deserialize)]
pub struct BundleResult {
    pub bundle_hash: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct JsonRpcError {
    pub code: i64,
    pub message: String,
}

/// 將已簽名套利交易以 bundle 形式送交 Flashbots
pub async fn send_bundle(
    relay_url: &str,
    signed_txs: Vec<String>,
    target_block: u64,
) -> Result<Option<String>, Box<dyn std::error::Error + Send + Sync>> {
    let url = format!("{}/", relay_url.trim_end_matches('/'));
    let body = FlashbotsBundleRequest {
        jsonrpc: "2.0".into(),
        id: 1,
        method: "eth_sendBundle".into(),
        params: vec![FlashbotsBundleParams {
            txs: signed_txs,
            block_number: format!("0x{:x}", target_block),
            min_timestamp: None,
            max_timestamp: None,
        }],
    };

    let client = Client::new();
    let res = client.post(&url).json(&body).send().await?;
    let status = res.status();
    let fb: FlashbotsBundleResponse = res.json().await?;

    if let Some(e) = fb.error {
        tracing::error!("Flashbots error: {} - {}", e.code, e.message);
        return Ok(None);
    }

    let bundle_hash = fb.result.and_then(|r| r.bundle_hash);
    if bundle_hash.is_some() {
        info!("Bundle submitted: {:?}", bundle_hash);
    }
    Ok(bundle_hash)
}

/// 取得當前區塊號（需由 caller 從 provider 取得後傳入）
pub fn relay_url_ethereum() -> String {
    std::env::var("FLASHBOTS_RELAY_URL").unwrap_or_else(|_| FLASHBOTS_RELAY_ETH.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn relay_url_default() {
        let u = relay_url_ethereum();
        assert!(u.contains("flashbots"));
    }
}
