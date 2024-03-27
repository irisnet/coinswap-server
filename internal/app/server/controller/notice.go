package controller

import (
	"context"
	"github.com/irisnet/coinswap-server/internal/app/pkg/kit"
	"github.com/irisnet/coinswap-server/internal/app/server/service"
)

type noticeController struct {
	BaseController
	ls service.NoticeService
}

func NewNoticeController(bc BaseController, ls service.NoticeService) noticeController {
	return noticeController{bc, ls}
}

func (lc noticeController) GetEndpoints() (endpoints []kit.Endpoint) {
	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/notice",
		Method:  "GET",
		Handler: lc.makeHandler(lc.QueryNotice, nil),
	})
	return endpoints
}

func (lc noticeController) QueryNotice(ctx context.Context, _ interface{}) (interface{}, error) {
	return lc.ls.QueryNotice()
}
