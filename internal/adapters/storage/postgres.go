package storage

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Util787/test-task/internal/common"
	"github.com/Util787/test-task/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ( // чтобы не загромождать конфиг
	defaultMaxConns        = 10
	defaultConnMaxLifetime = time.Hour
	defaultConnMaxIdleTime = time.Minute * 10
)

type PostgresStorage struct {
	pgxPool *pgxpool.Pool
}

func MustInitPostgres(ctx context.Context, cfg config.PostgresConfig) PostgresStorage {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)

	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(fmt.Errorf("failed to parse postgres connection string: %w", err))
	}

	// Pool configuration
	pgxConfig.MaxConns = defaultMaxConns
	pgxConfig.MaxConnLifetime = defaultConnMaxLifetime
	pgxConfig.MaxConnIdleTime = defaultConnMaxIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create postgres connection pool: %w", err))
	}

	err = pool.Ping(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to ping postgres: %w", err))
	}

	return PostgresStorage{
		pgxPool: pool,
	}
}

func (p *PostgresStorage) Shutdown() {
	p.pgxPool.Close()
}

// учитывая схему, пусть id массива всегда будет = 1

func (p *PostgresStorage) SaveNum(ctx context.Context, num int) error {
	op := common.GetOperationName()

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	upd := qb.Update("items").
		Set("numbers", sq.Expr("array_append(numbers, ?)", num)).
		Where(sq.Eq{"id": 1})

	sqlStr, args, err := upd.ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build update query: %w", op, err)
	}

	res, err := p.pgxPool.Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to execute update query: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		// если не было item с id=1, то создаем
		ins := qb.Insert("items").
			Columns("id", "numbers").
			Values(1, sq.Expr("ARRAY[?]::integer[]", num))

		sqlStr2, args2, _ := ins.ToSql()
		if _, err := p.pgxPool.Exec(ctx, sqlStr2, args2...); err != nil {
			return fmt.Errorf("%s: failed to execute insert query: %w", op, err)
		}
	}

	return nil
}
func (p *PostgresStorage) GetArr(ctx context.Context) ([]int, error) {
	op := common.GetOperationName()

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sel := qb.Select("numbers").
		From("items").
		Where(sq.Eq{"id": 1})

	sqlStr, args, err := sel.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build select query: %w", op, err)
	}

	var arr []int
	err = p.pgxPool.QueryRow(ctx, sqlStr, args...).Scan(&arr)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute select query: %w", op, err)
	}

	return arr, nil
}
