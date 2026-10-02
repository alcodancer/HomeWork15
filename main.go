// main.go (обновленный)
package main

import (
	"chess/internal/model/player"
	"chess/internal/repository"
	"chess/internal/service"
	"chess/internal/logger"
	"context"
	"fmt"
	"bufio"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"chess/internal/service"
	postgresstore "chess/internal/storage/postgres"
	"chess/internal/storage/postgres/sqlc"
	redisstore "chess/internal/storage/redis"
	httptransport "chess/internal/transport/http"
)

func main() {
	ctx := context.Background()

	const dsn = "postgres://app:secret@localhost:5432/appdb?sslmode=disable"

	// --- миграции ---
	if err := postgresstore.RunMigrations(dsn, "migrations"); err != nil {
		log.Fatal("migrate:", err)
	}

	// --- пул соединений ---
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("pgx pool:", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("pg ping:", err)
	}

	// --- sqlc ---
	q := sqlc.New(pool)

	// --- репозитории ---
	userRepo := postgresstore.NewUserRepo(q)
	orderRepo := postgresstore.NewOrderRepo(q)

	// --- redis ---
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("redis ping:", err)
	}
	historyRepo := redisstore.NewHistoryRepo(rdb, 7*24*time.Hour)

	// --- транзакции ---
	txRunner := postgresstore.NewTxRunner(pool)

	// --- сервис ---
	svc := service.NewUserService(userRepo, orderRepo, historyRepo, txRunner)

	// --- HTTP ---
	h := httptransport.NewUserHandler(svc)
	mux := http.NewServeMux()
	mux.HandleFunc("/users/create", h.Create)
	mux.HandleFunc("/users/with-order", h.CreateWithOrder)
	mux.HandleFunc("/users/get", h.Get)
	mux.HandleFunc("/users/update", h.Update)
	mux.HandleFunc("/users/delete", h.Delete)
	mux.HandleFunc("/users/list", h.List)

	log.Println("HTTP on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}