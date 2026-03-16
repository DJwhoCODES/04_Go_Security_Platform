package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/djwhocodes/auth-service/internal/config"
	"github.com/djwhocodes/auth-service/internal/database"
	"github.com/djwhocodes/auth-service/internal/handler"
	"github.com/djwhocodes/auth-service/internal/repository"
	"github.com/djwhocodes/auth-service/internal/security"
	"github.com/djwhocodes/auth-service/internal/service"
	"github.com/djwhocodes/auth-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	http   *http.Server

	db    *pgxpool.Pool
	redis *redis.Client
	cfg   *config.Config
}

func New(cfg *config.Config) *Server {

	router := gin.New()

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		logger.Log.Fatal("postgres connection failed", zap.Error(err))
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Log.Fatal("redis connection failed", zap.Error(err))
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	router.GET("/health", func(c *gin.Context) {

		ctx := c.Request.Context()

		if err := db.Ping(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "db_error",
			})
			return
		}

		if err := redisClient.Ping(ctx).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "redis_error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return &Server{
		engine: router,
		http:   httpServer,
		db:     db,
		redis:  redisClient,
		cfg:    cfg,
	}
}

func (s *Server) Start() {

	logger.Log.Info("starting server",
		zap.String("addr", s.http.Addr),
	)

	baseRepo := repository.NewRepository(s.db)

	userRepo := repository.NewUserRepository(baseRepo)
	tokenRepo := repository.NewTokenRepository(baseRepo)

	jwtManager := security.NewJWTManager(
		s.cfg.JWT.Secret,
		s.cfg.JWT.Issuer,
	)

	authService := service.NewAuthService(
		userRepo,
		tokenRepo,
		jwtManager,
	)

	authHandler := handler.NewAuthHandler(authService)

	auth := s.engine.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("server failed", zap.Error(err))
	}
}

func (s *Server) Shutdown(ctx context.Context) error {

	logger.Log.Info("shutting down server")

	if err := s.redis.Close(); err != nil {
		logger.Log.Error("redis close error", zap.Error(err))
	}

	s.db.Close()

	return s.http.Shutdown(ctx)
}
