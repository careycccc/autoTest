package sixearn

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

type MasterAgentQueryStruct struct {
	UserId                 int  `json:"userId"`
	IsAll                  bool `json:"isAll"`                  // 查询返佣代理线
	IsIncludeSelfAndParent bool `json:"isIncludeSelfAndParent"` // 列出自身及所有上级
	model.QueryPayloadStruct
}

// User 子用户/下级信息
type User struct {
	UserID          int64   `json:"userId"`          // 用户ID
	ParentID        int64   `json:"parentId"`        // 上级ID
	GeneralAgentID  int64   `json:"generalAgentId"`  // 总代ID（可能用于追踪根代理）
	Hierarchy       int     `json:"hierarchy"`       // 层级
	FirstChildCount int     `json:"firstChildCount"` // 直推人数
	ChildCount      int     `json:"childCount"`      // 团队总人数
	Balance         float64 `json:"balance"`         // 余额
	RegisterTime    int64   `json:"registerTime"`    // 注册时间戳（毫秒）
	LastLoginTime   int64   `json:"lastLoginTime"`   // 最后登录时间戳（毫秒）
	RebateState     int     `json:"rebateState"`     // 返佣状态
	RebateMode      int     `json:"rebateMode"`      // 返佣模式
	RebateLevel     int     `json:"rebateLevel"`     // 返佣等级
	RebateSetTime   int64   `json:"rebateSetTime"`   // 返佣设置时间
}

// PageData 分页数据
type PageData struct {
	List       []User `json:"list"`       // 用户列表
	PageNo     int    `json:"pageNo"`     // 当前页码
	TotalPage  int    `json:"totalPage"`  // 总页数
	TotalCount int    `json:"totalCount"` // 总记录数
}

// Response 整体响应结构
type Response struct {
	Data          PageData    `json:"data"`          // 业务数据
	MsgParameters any `json:"msgParameters"` // 通常为 null，可用 any 或 any
	Code          int         `json:"code"`          // 状态码
	Msg           string      `json:"msg"`           // 消息
	MsgCode       int         `json:"msgCode"`       // 消息码（有些系统区分 code 和 msgCode）
}

// 输入总代进行总代及其下级的团队人员查询
//
//	UserId int `json:"userId"`
//	IsAll bool `json:"isAll"` //查询返佣代理线
//	IsIncludeSelfAndParent bool `json:"isIncludeSelfAndParent"` // 列出自身及所有上级
func MasterAgentQuery(ctx *context.Context, userId int, isAll, isIncludeSelfAndParent bool) (*model.Response, []User, error) {
	api := "/api/Agent/GetPageListAgentList"
	payload := &MasterAgentQueryStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []any{userId, isAll, isIncludeSelfAndParent, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payload, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("api/Agent/GetPageListAgentList请求失败", err)), nil, err
	} else {
		if string(respBoy) == "null" {
			return model.HandlerErrorRes(model.ErrorLoggerType("api/Agent/GetPageListAgentList查询失败", nil)), nil, nil
		} else {
			var resposne Response
			if err := json.Unmarshal(respBoy, &resposne); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("api/Agent/GetPageListAgentList解析失败", err)), nil, err
			} else {
				if resp, err := model.ParseResponse(respBoy); err != nil {
					return model.HandlerErrorRes(model.ErrorLoggerType("api/Agent/GetPageListAgentList解析失败", err)), nil, err
				} else {
					if len(resposne.Data.List) == 0 {
						return resp, nil, nil
					} else {
						return resp, resposne.Data.List, nil
					}
				}
			}
		}
	}
}
