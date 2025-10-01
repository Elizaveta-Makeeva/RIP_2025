package revenueHandler

import (
	"lab2/internal/app/revenueModel"
)

type RevenueHandler struct {
	RevenueModel *revenueModel.RevenueModel
}

func NewRevenueHandler(r *revenueModel.RevenueModel) *RevenueHandler {
	return &RevenueHandler{
		RevenueModel: r,
	}
}
