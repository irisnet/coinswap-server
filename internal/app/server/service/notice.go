package service

import (
	"github.com/irisnet/coinswap-server/internal/app/server/types"
)

type NoticeService struct {
}

func NewNoticeService() NoticeService {
	return NoticeService{}
}
func (ls *NoticeService) QueryNotice() (types.QueryNotice, error) {
	response, err := notice.FindLatestCreateAtOne()
	if err != nil {
		return types.QueryNotice{}, err
	}

	return types.QueryNotice{Notice: response.Content}, nil
}
