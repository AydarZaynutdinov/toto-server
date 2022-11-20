package service

import (
	"context"

	"github.com/valyala/fasthttp"
)

type IService interface {
	Get(ctx context.Context) fasthttp.RequestHandler
}
