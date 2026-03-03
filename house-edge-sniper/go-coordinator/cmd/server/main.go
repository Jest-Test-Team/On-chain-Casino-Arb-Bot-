// 訊息中樞與 API：消費 Kafka、風控、WebSocket 推送
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/house-edge-sniper/go-coordinator/internal/kafka_consumer"
	"github.com/house-edge-sniper/go-coordinator/internal/risk_manager"
	"github.com/house-edge-sniper/go-coordinator/internal/ws_hub"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	kafkaBrokers := envList("KAFKA_BROKERS", "localhost:9092")
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("redis ping", zap.Error(err))
	}
	defer rdb.Close()

	rm := risk_manager.New(rdb, log, risk_manager.Options{
		MaxDrawdownPct: 20.0,
	})

	hub := ws_hub.New(log)
	go hub.Run()

	// HTTP：WebSocket 與 API
	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		upgrader := ws_hub.Upgrader()
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		hub.Register(conn)
	})
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/api/risk", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"max_drawdown_pct": rm.GetMaxDrawdownPct()})
	})

	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http server", zap.Error(err))
		}
	}()

	// Kafka consumer：收到機會後經風控，再廣播給前端
	handler := func(ctx context.Context, opp kafka_consumer.Opportunity) error {
		if !rm.Allow(ctx, opp) {
			return nil
		}
		hub.BroadcastJSON(gin.H{"type": "opportunity", "payload": opp})
		return nil
	}
	consumer := kafka_consumer.New(kafkaBrokers, "go-coordinator", log, handler)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := consumer.Run(ctx); err != nil && ctx.Err() == nil {
			log.Error("kafka consumer", zap.Error(err))
		}
		_ = consumer.Close()
	}()

	log.Info("go-coordinator listening", zap.String("http", "8080"))
	<-ctx.Done()
	_ = srv.Shutdown(context.Background())
	log.Info("shutdown complete")
}

func envList(key, defaultVal string) []string {
	s := os.Getenv(key)
	if s == "" {
		s = defaultVal
	}
	return strings.Split(s, ",")
}
