// Package risk_manager 控管錢包餘額與最大回撤
package risk_manager

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/house-edge-sniper/go-coordinator/internal/kafka_consumer"
)

const (
	redisKeyWalletBalancePrefix = "house-edge:wallet:"
	redisKeyDrawdown            = "house-edge:drawdown"
	redisKeyMaxDrawdownPct       = "house-edge:max_drawdown_pct"
)

// Manager 風控：檢查餘額、最大回撤後決定是否放行
type Manager struct {
	rdb   *redis.Client
	log   *zap.Logger
	mu    sync.RWMutex
	opts  Options
}

// Options 風控參數
type Options struct {
	MinBalanceWei   string  // 最低錢包餘額（wei 字串）
	MaxDrawdownPct  float64 // 最大回撤百分比，超過則停止下單
	DefaultDrawdown float64 // 當前回撤（可從 Redis 讀取）
}

// New 建立風控經理
func New(rdb *redis.Client, log *zap.Logger, opts Options) *Manager {
	if opts.MaxDrawdownPct <= 0 {
		opts.MaxDrawdownPct = 20.0
	}
	return &Manager{rdb: rdb, log: log, opts: opts}
}

// Allow 檢查是否允許執行此機會（餘額、回撤）
func (m *Manager) Allow(ctx context.Context, opp kafka_consumer.Opportunity) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 從 Redis 讀取當前回撤（若有的話）
	pct, err := m.rdb.HGet(ctx, redisKeyDrawdown, "pct").Float64()
	if err == nil && pct >= m.opts.MaxDrawdownPct {
		m.log.Warn("rejected: max drawdown exceeded", zap.Float64("current_pct", pct))
		return false
	}

	// 可擴展：依 wallet 從 Redis 讀取餘額並與 MinBalanceWei 比較
	return true
}

// RecordResult 記錄一筆執行結果，更新回撤等統計（供 dashboard 與風控用）
func (m *Manager) RecordResult(ctx context.Context, wallet string, profitLossWei string, newDrawdownPct float64) error {
	pipe := m.rdb.Pipeline()
	pipe.HSet(ctx, redisKeyDrawdown, "pct", newDrawdownPct)
	pipe.HSet(ctx, redisKeyWalletBalancePrefix+wallet, "last_pl", profitLossWei)
	_, err := pipe.Exec(ctx)
	return err
}

// GetMaxDrawdownPct 回傳設定的最大回撤比例
func (m *Manager) GetMaxDrawdownPct() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.opts.MaxDrawdownPct
}
