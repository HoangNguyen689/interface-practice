package postgresql

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pgxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jackc/pgx.v5"
)

type Client struct {
	*pgxpool.Pool
	logger *log.Logger
}

func NewClient(
	ctx context.Context,
	endpoint,
	dbName,
	username,
	password,
	serviceName string,
) (*Client, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, endpoint, dbName)

	pool, err := pgxtrace.NewPool(ctx, dsn, pgxtrace.WithServiceName(serviceName))
	if err != nil {
		return nil, err
	}

	return &Client{
		Pool:   pool,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

func (c *Client) WaitReady(ctx context.Context, duration, timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err = c.Ping(ctx)
			if err == nil {
				return nil
			}
			c.logger.Println("postgresql db is still not ready, waiting...")

		case <-ctx.Done():
			return
		}
	}
}
