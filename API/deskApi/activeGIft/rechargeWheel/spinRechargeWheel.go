package rechargewheel

import "autoTest/store/model"

// 旋转充值转盘
type SpinRechargeWheel struct {
	RechargeWheelType any `json:"rechargeWheelType"` // 充值转盘类型
	model.BaseStruct
}

// func SpinRechargeWheelApi(ctx *context.Context) (*model.Response, error) {}
