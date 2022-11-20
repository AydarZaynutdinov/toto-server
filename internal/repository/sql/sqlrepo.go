package sql

import (
	"context"
	"fmt"
	"strings"
	"toto-server/internal/consts"

	"toto-server/internal/entity"
	"toto-server/internal/repository"

	"github.com/jackc/pgx/v4"
)

type SkuConfigSqlRepo struct {
	conn *pgx.Conn
	ctx  context.Context
}

func New(opts ...OptionFunc) repository.ISkuConfig {
	s := &Options{}
	for _, opt := range opts {
		opt(s)
	}

	if s.ctx == nil {
		s.ctx = context.Background()
	}

	return &SkuConfigSqlRepo{
		conn: s.sqlDb.GetConn(),
		ctx:  s.ctx,
	}
}

func (s *SkuConfigSqlRepo) Get(ctx context.Context, params repository.QueryParameters) (entity.SkuConfig, error) {
	filters, args := getFilters(params)
	query := fmt.Sprintf("SELECT id, package, country_code, percentile_min, percentile_max, main_sku "+
		"FROM %s "+
		"WHERE %s "+
		"LIMIT %d;",
		consts.SkuConfigTableName, filters, 1)
	var skuConfig entity.SkuConfig

	err := s.conn.QueryRow(ctx, query, args...).
		Scan(&skuConfig.ID, &skuConfig.Package, &skuConfig.CountryCode, &skuConfig.PercentileMin, &skuConfig.PercentileMax, &skuConfig.MainSku)
	return skuConfig, err
}

func getFilters(params repository.QueryParameters) (string, []interface{}) {
	filters := strings.Builder{}
	filters.WriteString("true")
	var args []interface{}

	ind := 1
	if len(params.Packages) > 0 {
		filters.WriteString(fmt.Sprintf(" AND package LIKE ANY ($%d)", ind))
		args = append(args, params.Packages)
		ind += 1
	}

	if len(params.CountryCode) > 0 {
		filters.WriteString(fmt.Sprintf(" AND country_code LIKE ANY ($%d)", ind))
		args = append(args, params.CountryCode)
		ind += 1
	}

	if params.Percentile > 0 {
		filters.WriteString(fmt.Sprintf(" AND $%d BETWEEN percentile_min + 1 AND percentile_max", ind))
		args = append(args, params.Percentile)
		ind += 1
	}

	return filters.String(), args
}
