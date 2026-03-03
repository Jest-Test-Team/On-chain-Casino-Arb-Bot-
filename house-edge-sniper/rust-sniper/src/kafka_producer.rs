//! 將發現的套利機會寫入 Kafka，供 Go 協調器消費與風控

use crate::ev_calculator::ArbOpportunity;
use rdkafka::config::ClientConfig;
use rdkafka::producer::{FutureProducer, FutureRecord};
use serde::Serialize;
use std::time::Duration;
use tracing::debug;

const TOPIC_OPPORTUNITIES: &str = "house-edge.opportunities";

pub struct KafkaProducer {
    producer: FutureProducer,
}

impl KafkaProducer {
    pub fn new(brokers: &str) -> Result<Self, rdkafka::error::KafkaError> {
        let producer: FutureProducer = ClientConfig::new()
            .set("bootstrap.servers", brokers)
            .set("message.timeout.ms", "5000")
            .set("queue.buffering.max.messages", "100000")
            .create()?;
        Ok(Self { producer })
    }

    /// 非同步發送一筆套利機會
    pub async fn send_opportunity(&self, opp: &ArbOpportunity) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        let payload = serde_json::to_string(opp)?;
        let key = format!("{}:{}", opp.chain, opp.tx_hash);
        let record = FutureRecord::to(TOPIC_OPPORTUNITIES)
            .key(&key)
            .payload(&payload)
            .headers(rdkafka::message::OwnedHeaders::new().insert(rdkafka::message::Header {
                key: "source",
                value: Some("rust-sniper"),
            }));

        self.producer
            .send(record, Duration::from_secs(0))
            .await
            .map_err(|(e, _)| e)?;
        debug!("Sent opportunity to {}: {}", TOPIC_OPPORTUNITIES, key);
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::ev_calculator::ArbOpportunity;

    #[test]
    fn topic_constant() {
        assert_eq!(TOPIC_OPPORTUNITIES, "house-edge.opportunities");
    }
}
