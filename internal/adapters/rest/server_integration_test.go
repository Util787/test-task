package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/Util787/test-task/internal/adapters/storage"
	"github.com/Util787/test-task/internal/config"
	"github.com/Util787/test-task/internal/logger/slogpretty"
	"github.com/Util787/test-task/internal/usecase"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	defer postgresC.Terminate(ctx)

	host, _ := postgresC.Host(ctx)
	port, _ := postgresC.MappedPort(ctx, "5432")

	cfg := config.Config{
		PostgresConfig: config.PostgresConfig{
			Host:     host,
			Port:     port.Int(),
			User:     "test_user",
			Password: "test_password",
			DbName:   "test_db",
		},
		HTTPServerConfig: config.HTTPServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}

	logger := slogpretty.NewPrettyLogger(os.Stdout, slog.LevelDebug)

	var postgreStorage storage.PostgresStorage
	for i := 0; i < 5; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Postgres not ready yet, retrying...")
				}
			}()
			postgreStorage = storage.MustInitPostgres(ctx, cfg.PostgresConfig)
		}()
		time.Sleep(time.Second)
	}

	defer postgreStorage.Shutdown()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresConfig.User,
		cfg.PostgresConfig.Password,
		cfg.PostgresConfig.Host,
		cfg.PostgresConfig.Port,
		cfg.PostgresConfig.DbName,
	)

	migrationPath, err := filepath.Abs("../../../migrations/postgres")
	if err != nil {
		log.Fatalf("Failed to get absolute path for migrations: %v", err)
	}

	mi, err := migrate.New(
		"file://"+filepath.ToSlash(migrationPath),
		dsn,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := mi.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	sortUsecase := usecase.NewSortUsecase(&postgreStorage)

	server := NewRestServer(logger, cfg.HTTPServerConfig, sortUsecase)
	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	defer server.Shutdown(ctx)

	time.Sleep(2 * time.Second)

	code := m.Run()

	os.Exit(code)
}

func TestSortEndpoint(t *testing.T) {
	requests := []saveNumRequest{
		{Num: 1024},
		{Num: 65536},
		{Num: 99999},
		{Num: 1024},
		{Num: 5},
	}

	var allNums []int

	for i, reqData := range requests {
		allNums = append(allNums, reqData.Num)

		body, err := json.Marshal(reqData)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		resp, err := http.Post("http://localhost:8080/api/v1/sort", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to make POST request: %v", err)
		}

		var result []int
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", resp.StatusCode)
		}

		if len(result) != i+1 {
			t.Errorf("Expected %d numbers in response, got %d", i+1, len(result))
		}

		if i == len(requests)-1 {
			sort.Ints(allNums)
			expected := allNums

			if len(result) != len(expected) {
				t.Errorf("Expected %d numbers in response, got %d", len(expected), len(result))
			}

			for i := range result {
				if result[i] != expected[i] {
					t.Errorf("Expected sorted result %v, got %v", expected, result)
					break
				}
			}

		}
	}
}
