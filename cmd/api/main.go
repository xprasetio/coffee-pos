package main

import (
	"fmt"
	"os"

	"github.com/xprasetio/coffee-pos/config"
	"github.com/xprasetio/coffee-pos/internal/handler"
	"github.com/xprasetio/coffee-pos/pkg/database"
	"github.com/xprasetio/coffee-pos/pkg/redis"
	"github.com/xprasetio/coffee-pos/pkg/validator"
)

func main() {
	fmt.Println("Coffee Shop POS starting...")

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize MySQL connection
	db, err := database.NewMySQL(cfg.MysqlDSN())
	if err != nil {
		fmt.Printf("Error connecting to MySQL: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println("MySQL connected.")

	// Initialize Redis connection
	rdb, err := redis.NewRedis(cfg.RedisAddr(), cfg.RedisPassword)
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		os.Exit(1)
	}
	defer rdb.Close()
	fmt.Println("Redis connected.")

	// Initialize validator
	v := validator.New()

	// Initialize router
	router := handler.NewRouter(db, cfg, v)

	// Start server
	if err := router.Run(":" + cfg.AppPort); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}
}
