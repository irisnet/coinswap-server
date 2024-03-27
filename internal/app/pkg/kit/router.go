package kit

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// TimeLayout time format
const TimeLayout = "2006-01-02"

type (
	//Router define a router for http Handler
	Endpoint struct {
		URI     string
		Method  string
		Handler http.Handler
	}

	Handler            = endpoint.Endpoint
	Server             = httptransport.Server
	RequestFunc        = httptransport.RequestFunc
	ServerResponseFunc = httptransport.ServerResponseFunc

	//IController define a interface for all http Controller
	IController interface {
		GetEndpoints() []Endpoint
	}

	Application interface {
		IController
		Initialize()
		Stop()
	}

	Controller struct {
		validate *validator.Validate
	}
)

func NewController() Controller {
	return Controller{validator.New()}
}

// makeHandler create a http hander for request
func (c Controller) MakeHandler(handler endpoint.Endpoint, request interface{},
	before []httptransport.RequestFunc,
	mid []httptransport.ServerOption,
	after []httptransport.ServerResponseFunc,
) *httptransport.Server {
	return httptransport.NewServer(
		handler,
		c.decodeRequest(request),
		c.encodeResponse,
		c.serverOptions(before, mid, after)...,
	)
}

func (c Controller) GetIntValue(ctx context.Context, key string) (int, error) {
	value := ctx.Value(key)
	if value == nil {
		return 0, errors.Errorf("Not found key: %s in Context", key)
	}

	v, err := strconv.ParseInt(value.(string), 10, 64)
	if err != nil {
		logger.Error("Invalid key, must be int type")
		return 0, errors.Errorf("Value: %s is not int type", value)
	}
	return int(v), nil
}

func (c Controller) GetStringValue(ctx context.Context, key string) (string, error) {
	value := ctx.Value(key)
	if value == nil {
		return "", errors.Errorf("Not found key: %s in Context", key)
	}

	v, ok := value.(string)
	if !ok {
		logger.Error("Invalid key, must be string type")
		return "", errors.Errorf("Value: %s is not string type", value)
	}
	return v, nil
}

func (c Controller) GetDateValue(ctx context.Context, key string) (*time.Time, error) {
	value := ctx.Value(key)
	if value == nil {
		return nil, errors.Errorf("Not found key: %s in Context", key)
	}

	tim, err := time.Parse(TimeLayout, value.(string))
	if err != nil {
		logger.Error("Invalid key, must be string type")
		return nil, errors.Errorf("Value: %s is not string type", value)
	}
	return &tim, nil
}

func (c Controller) GetPagation(ctx context.Context) (int, int) {
	page, err := c.GetIntValue(ctx, "page")
	if err != nil {
		page = 1
	}

	size, err := c.GetIntValue(ctx, "size")
	if err != nil {
		size = 10
	}
	return page, size
}

// decodeRequest decode request(http.request -> model.request)
func (c Controller) decodeRequest(req interface{}) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		logger.Debug("Execute decode request", logger.String("method", "decodeRequest"))
		if req == nil {
			return nil, err
		}

		reqNew := New(req)
		if err := json.NewDecoder(r.Body).Decode(&reqNew); err != nil {
			logger.Error("Execute decode request failed", logger.String("error", err.Error()))
			return nil, err
		}

		//validate request
		if err := c.validate.Struct(req); err != nil {
			return nil, err
		}

		return reqNew, nil
	}
}

// encodeResponse encode the  response(model.response -> http.response)
func (c Controller) encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	logger.Debug("Execute encode response", logger.String("method", "encodeResponse"))
	response := Response{
		Code:    OK,
		Message: "success",
		Data:    resp,
	}
	return httptransport.EncodeJSONResponse(ctx, w, response)
}

func (c Controller) serverOptions(
	before []httptransport.RequestFunc,
	mid []httptransport.ServerOption,
	after []httptransport.ServerResponseFunc,
) []httptransport.ServerOption {
	//copy params from Form,PostForm to Context
	copyParams := func(ctx context.Context, request *http.Request) context.Context {
		logger.Debug("Merge request params to Context", logger.String("method", "serverBefore"))

		if err := request.ParseForm(); err != nil {
			logger.Error("Parse form failed", logger.String("error", err.Error()))
			return ctx
		}

		improveValue := func(vs []string) interface{} {
			if len(vs) == 1 {
				return vs[0]
			}
			return vs
		}

		for k, v := range request.Form {
			ctx = context.WithValue(ctx, k, improveValue(v))
		}

		for k, v := range request.PostForm {
			ctx = context.WithValue(ctx, k, improveValue(v))
		}

		for k, v := range mux.Vars(request) {
			ctx = context.WithValue(ctx, k, v)
		}

		return ctx
	}

	//format error
	errorEncoderOption := func(ctx context.Context, err error, w http.ResponseWriter) {
		response := Response{
			Code:    Failed,
			Message: err.Error(),
		}
		bz, _ := json.Marshal(response)
		_, _ = w.Write(bz)
	}

	var options []httptransport.ServerOption
	before = append(
		[]httptransport.RequestFunc{httptransport.PopulateRequestContext, copyParams},
		before...,
	)
	options = append(options, httptransport.ServerBefore(before...))
	options = append(options, append(mid, httptransport.ServerErrorEncoder(errorEncoderOption))...)
	options = append(options, httptransport.ServerAfter(after...))
	return options
}
