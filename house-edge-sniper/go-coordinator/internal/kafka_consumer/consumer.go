// Package kafka_consumer 訂閱套利訊號與執行結果
package kafka_consumer

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	TopicOpportunities = "house-edge.opportunities"
	TopicExecution     = "house-edge.execution"
)

// Opportunity 與 Rust ArbOpportunity 對應
type Opportunity struct {
	Chain             string  `json:"chain"`
	TxHash            string  `json:"tx_hash"`
	ContractAddress   *string `json:"contract_address,omitempty"`
	ExpectedValueWei  *string `json:"expected_value_wei,omitempty"`
	TimestampMs       int64   `json:"timestamp_ms"`
}

// Handler 收到機會時的回呼（可由 risk_manager 過濾後再執行）
type Handler func(ctx context.Context, opp Opportunity) error

// Consumer 訂閱 Kafka 並將訊息轉發給 Handler
type Consumer struct {
	reader *kafka.Reader
	log    *zap.Logger
	handle Handler
	once   sync.Once
}

// New 建立 consumer，訂閱 house-edge.opportunities
func New(brokers []string, groupID string, log *zap.Logger, handle Handler) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    TopicOpportunities,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 1e6,
	})
	return &Consumer{reader: r, log: log, handle: handle}
}

// Run 在背景消費，遇到錯誤或 ctx 取消時結束
func (c *Consumer) Run(ctx context.Context) error {
	c.log.Info("Kafka consumer started", zap.String("topic", TopicOpportunities))
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return err
			}
			var opp Opportunity
			if err := json.Unmarshal(msg.Value, &opp); err != nil {
				c.log.Warn("invalid opportunity payload", zap.Error(err), zap.ByteString("value", msg.Value))
				continue
			}
			if c.handle != nil {
				if err := c.handle(ctx, opp); err != nil {
					c.log.Warn("handler error", zap.Error(err), zap.String("tx_hash", opp.TxHash))
				}
			}
		}
	}
}

// Close 關閉 reader
func (c *Consumer) Close() error {
	var err error
	c.once.Do(func() { err = c.reader.Close() })
	return err
}
