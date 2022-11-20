package sql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"toto-server/config"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	conn *pgx.Conn
	once sync.Once
)

type SQL struct {
	conn *pgx.Conn
	ctx  context.Context
}

func ProvideSQLConnection(cfg *config.DBConfig, ctx context.Context) (*SQL, error) {
	if cfg == nil {
		return nil, errors.New("cannot establish db connection due to empty config")
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	once.Do(func() {
		url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		connection, err := pgx.Connect(ctx, url)
		if err != nil {
			log.Fatalln(err)
		}
		conn = connection
	})

	s := &SQL{
		ctx:  ctx,
		conn: conn,
	}

	select {
	case <-ctxWithTimeout.Done():
		return nil, errors.New("exiting due to timeout trying to connect to db")
	default:
		go s.checkContext()
		return s, nil
	}
}

func (s *SQL) GetConn() *pgx.Conn {
	return s.conn
}

func (s *SQL) checkContext() {
	for {
		select {
		case <-s.ctx.Done():
			_ = s.conn.Close(context.Background())
			return
		default:
		}
	}
}
