package sql

import (
	"context"
)

type Options struct {
	ctx   context.Context
	sqlDb *SQL
}

type OptionFunc func(*Options)

func WithContext(ctx context.Context) OptionFunc {
	return func(s *Options) {
		s.ctx = ctx
	}
}

func WithConnection(sqlDb *SQL) OptionFunc {
	return func(s *Options) {
		s.sqlDb = sqlDb
	}
}
