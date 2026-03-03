//! 期望值與合約邏輯分析
//! 輸入：合約位址、calldata、鏈上狀態 → 輸出：是否為正 EV、建議 gas

use serde::{Deserialize, Serialize};

/// 單一套利機會描述，供 Kafka 與 Go 協調器消費
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ArbOpportunity {
    pub chain: String,
    pub tx_hash: String,
    pub contract_address: Option<String>,
    /// 期望獲利 (wei)，若為 None 表示尚未計算完
    pub expected_value_wei: Option<String>,
    pub timestamp_ms: i64,
}

/// 根據合約與 calldata 計算期望值（簡化版：僅結構，實際需接鏈上 call）
pub fn calculate_ev(
    _contract: &str,
    _calldata: &[u8],
    _pool_state_json: Option<&str>,
) -> Option<ArbOpportunity> {
    // TODO: 呼叫合約 view 函數取得賠率、池子餘額，計算 EV
    // 若 EV > 閾值則回傳 Some(ArbOpportunity { ... })
    None
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn arb_opportunity_serializes() {
        let o = ArbOpportunity {
            chain: "ethereum".into(),
            tx_hash: "0xabc".into(),
            contract_address: Some("0x123".into()),
            expected_value_wei: Some("1000000000000000000".into()),
            timestamp_ms: 1234567890,
        };
        let j = serde_json::to_string(&o).unwrap();
        assert!(j.contains("ethereum"));
    }
}
