package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"
	"toto-server/internal/consts"
	"toto-server/internal/repository"
	"toto-server/internal/response"

	"github.com/go-redis/redis/v9"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Service struct {
	log         *zap.Logger
	repo        repository.ISkuConfig
	redisClient *redis.Client
	random      *rand.Rand
}

func New(opts ...OptionFunc) IService {
	s := &Options{}
	for _, opt := range opts {
		opt(s)
	}

	return &Service{
		log:         s.Log,
		repo:        s.SkuConfigRepo,
		redisClient: s.RedisClient,
		random:      rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (s *Service) Get(ctx context.Context) fasthttp.RequestHandler {
	return func(requestCtx *fasthttp.RequestCtx) {
		params, err := s.getRequestParameters(requestCtx, ctx)
		if err != nil {
			return
		}

		skuConfig, err := s.repo.Get(ctx, params)
		if err != nil {
			if err == pgx.ErrNoRows {
				s.log.Sugar().Info(map[string]interface{}{
					consts.Module: consts.ServiceModule,
					consts.Action: "get sku_config from db",
					consts.Params: fmt.Sprintf("params: %v", params),
					consts.Error:  "there is now sku_config by params",
				})
				s.prepareErrorResponse(requestCtx, fasthttp.StatusNotFound, "there is no sku config with received parameters")
				return
			}
			s.log.Sugar().Error(map[string]interface{}{
				consts.Module: consts.ServiceModule,
				consts.Action: "get sku_config from db",
				consts.Params: fmt.Sprintf("params: %v", params),
				consts.Error:  err.Error(),
			})
			requestCtx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		body := response.MainSkuResponse{MainSku: skuConfig.MainSku}
		payload, err := json.Marshal(body)
		if err != nil {
			s.log.Sugar().Error(map[string]interface{}{
				consts.Module: consts.ServiceModule,
				consts.Action: "marshal response's body",
				consts.Params: fmt.Sprintf("body: %v", body),
				consts.Error:  err.Error(),
			})
			requestCtx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		requestCtx.SetStatusCode(fasthttp.StatusOK)
		requestCtx.SetBody(payload)
	}
}

func (s *Service) getRequestParameters(requestCtx *fasthttp.RequestCtx, ctx context.Context) (repository.QueryParameters, error) {
	// get package from requestCtx
	pack := string(requestCtx.QueryArgs().Peek("package"))
	if pack == "" {
		s.log.Sugar().Error(map[string]interface{}{
			consts.Module: consts.ServiceModule,
			consts.Action: "parse request",
			consts.Params: fmt.Sprintf("pack: %s", pack),
		})
		s.prepareErrorResponse(requestCtx, fasthttp.StatusBadRequest, "empty 'package'")
		return repository.QueryParameters{}, fmt.Errorf("empty 'package'")
	}

	// get country_code by request IP
	requestIP := requestCtx.RemoteIP()
	countryCode := s.getCountryCode(ctx, requestIP)

	params := repository.QueryParameters{
		Packages:    []string{pack},
		CountryCode: []string{countryCode},
		Percentile:  s.random.Intn(consts.MaxPercentile) + 1,
	}
	return params, nil
}

func (s *Service) prepareErrorResponse(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	ctx.SetStatusCode(statusCode)
	body := response.ErrorResponse{message}
	payload, err := json.Marshal(body)
	if err != nil {
		s.log.Sugar().Error(map[string]interface{}{
			consts.Module: consts.ServiceModule,
			consts.Action: "marshal error's body",
			consts.Params: fmt.Sprintf("body: %v", body),
			consts.Error:  err.Error(),
		})
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(payload)
}

func (s *Service) getCountryCode(ctx context.Context, ip net.IP) string {
	key := fmt.Sprintf("%s%s", consts.RedisIPKey, ip.String())
	redisVal := s.redisClient.Get(ctx, key)

	if redisVal.Err() == nil {
		var countryCode string
		err := redisVal.Scan(&countryCode)
		if err == nil {
			return countryCode
		}
	}

	countryCode, err := s.callGeoService(ip)
	if err != nil {
		return consts.DefaultCountryCode
	}

	s.redisClient.Set(ctx, key, countryCode, -1)
	return countryCode
}

func (s *Service) callGeoService(ip net.IP) (string, error) {
	url := fmt.Sprintf("%s%s", consts.GeoServiceUrl, ip.String())
	resp, err := http.Get(url)
	if err != nil {
		s.log.Sugar().Error(map[string]interface{}{
			consts.Module: consts.ServiceModule,
			consts.Action: "call geo service",
			consts.Params: fmt.Sprintf("url: %v", url),
			consts.Error:  err.Error(),
		})
		return "", err
	}
	if resp.StatusCode != 200 {
		s.log.Sugar().Error(map[string]interface{}{
			consts.Module: consts.ServiceModule,
			consts.Action: "call geo service",
			consts.Params: fmt.Sprintf("url: %v", url),
			consts.Status: resp.StatusCode,
		})
		return "", fmt.Errorf("bad response status")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Sugar().Error(map[string]interface{}{
			consts.Module: consts.ServiceModule,
			consts.Action: "read response's body",
			consts.Params: fmt.Sprintf("ip: %v", ip),
			consts.Error:  err.Error(),
		})
		return "", err
	}

	geoResp := &response.GeoResponse{}
	if err = json.Unmarshal(data, geoResp); err != nil {
		s.log.Sugar().Error(map[string]interface{}{
			consts.Module: consts.ServiceModule,
			consts.Action: "unmarshal response's body",
			consts.Params: fmt.Sprintf("body: %s", data),
			consts.Error:  err.Error(),
		})
		return "", err
	}

	return geoResp.Country.Iso, nil
}
