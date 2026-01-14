package memberlist

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"fmt"
)

// 会员列表

type GetUserApistruct struct {
	Account  string `json:"account"`  // 用户账号
	PageNo   int    `json:"pageNo"`   // 页码
	PageSize int    `json:"pageSize"` // 每页多少数据
	OrderBy  string `json:"orderBy"`  // 排序
	model.BaseStruct
}

// 提取userid
type Useridstruct struct {
	Data struct {
		List []struct {
			UserId int64 `json:"userId"`
		} `json:"list"`
	} `json:"data"`
}

/*
输入用户电话号码，返回用户id
*
*/
func GetUserIdApi(ctx *context.Context, account string) (*model.Response, int64, error) {
	defer func() {
		if r := recover(); r != nil {
			// 处理 panic，记录日志或其他操作
			logger.LogError("GetUserIdApi 函数发生 panic", fmt.Errorf("%v", r))
			logger.Logger.Warn("account 参数值为", account)
		}
	}()
	api := "/api/Users/GetPageList"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &GetUserApistruct{}
	payloadList := []interface{}{account, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList请求失败", err)), -1, err
	} else {
		// 提取用户id
		var response Useridstruct
		if err := json.Unmarshal(respBoy, &response); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList提取用户id解析失败", err)), -1, err
		} else {
			if resp, err := model.ParseResponse(respBoy); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList响应解析失败", err)), -1, err
			} else {
				if len(response.Data.List) == 0 {
					return resp, -1, nil
				} else {
					return resp, response.Data.List[0].UserId, nil
				}
			}
		}

	}
}

// 提取用户邀请码
type UseridInviteCodestruct struct {
	Data struct {
		List []struct {
			InviteCode string `json:"inviteCode"`
		} `json:"list"`
	} `json:"data"`
}

/*
返回用户邀请码
*
*/
func GetUserInviteCodeApi(ctx *context.Context, account string) (*model.Response, string, error) {
	defer func() {
		if r := recover(); r != nil {
			// 处理 panic，记录日志或其他操作
			logger.LogError("GetUserIdApi 函数发生 panic", fmt.Errorf("%v", r))
			logger.Logger.Warn("account 参数值为", account)
		}
	}()
	api := "/api/Users/GetPageList"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &GetUserApistruct{}
	payloadList := []interface{}{account, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList请求失败", err)), "", err
	} else {
		// 提取用户id
		var response UseridInviteCodestruct
		if err := json.Unmarshal(respBoy, &response); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList提取用户id解析失败", err)), "", err
		} else {
			if resp, err := model.ParseResponse(respBoy); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList响应解析失败", err)), "", err
			} else {
				return resp, response.Data.List[0].InviteCode, nil
			}
		}

	}
}

// 后台添加用户的
type addUserInfoStruct struct {
	Account      any `json:"account"` // 添加的用户的账号91号码
	UserType     any `json:"userType"`
	PassWord     any `json:"password"`
	Remark       any `json:"remark"`
	RegisterType any `json:"registerType"`
}

// 添加用户
type AddUserStruct struct {
	AddUserList []addUserInfoStruct `json:"addUserList"`
	Random      any                 `json:"random"`
	Language    any                 `json:"language"`
	Signature   any                 `json:"signature"`
	Timestamp   any                 `json:"timestamp"`
}

// 需要传入用户的电话号码
func AddUsers(ctx *context.Context, userAmount string) (*model.Response, error) {
	api := "/api/Users/AddUsers"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &AddUserStruct{}
	payloadList := []interface{}{userAmount, 0, config.SUB_PWD, "", 1, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/AddUsers请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/AddUsers解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

type GetUserListStruct struct {
	UserType any `json:"userType"` // 用户类型
	model.QueryPayloadStruct
}

type GetUserVipListStruct struct {
	VipLevel any `json:"vipLevel"` // vip等级
	UserType any `json:"userType"` // 用户类型
	model.QueryPayloadStruct
}

type UserInfo struct {
	UserId           int64   `json:"userId"`
	Account          string  `json:"account"`
	NickName         string  `json:"nickName"`
	VipLevel         int32   `json:"vipLevel"`
	Balance          float64 `json:"balance"`
	ParentId         int64   `json:"parentId"`
	GeneralAgentId   int64   `json:"generalAgentId"`
	State            int32   `json:"state"`
	IsFuncation      bool    `json:"isFuncation"`
	RegisterTime     int64   `json:"registerTime"`
	RegisterSource   int32   `json:"registerSource"`
	LastLoginTime    int64   `json:"lastLoginTime"`
	Remark           string  `json:"remark"`
	UserType         int32   `json:"userType"`
	InviteCode       string  `json:"inviteCode"`
	ChannelId        int32   `json:"channelId"`
	RegisterIp       string  `json:"registerIp"`
	LastLoginIp      string  `json:"lastLoginIp"`
	IsBlackListIPTag bool    `json:"isBlackListIPTag"`
	PackageId        int32   `json:"packageId"`
	PackageName      string  `json:"packageName"`
}

type GetUserListResponseStruct struct {
	Data struct {
		List []UserInfo `json:"list"`
	} `json:"data"`
}

/*
userNumber 传入一个数字，返回对应数量的用户列表的详细信息
userType 0 正式账号 1 测试账号 2 游客账号
*
*/
func GetUserListApi(ctx *context.Context, pageSize, userNumber int, userType int8) (*model.Response, []*UserInfo, error) {
	api := "/api/Users/GetPageList"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &GetUserListStruct{}
	payloadList := []interface{}{userType, pageSize, userNumber, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList请求失败", err)), nil, err
	} else {

		// 提取用户id
		var response GetUserListResponseStruct
		if err := json.Unmarshal(respBoy, &response); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList提取用户列表解析失败", err)), nil, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList解析失败", err)), nil, err
		} else {

			userList := make([]*UserInfo, userNumber)
			for i := 0; i < userNumber; i++ {
				userList[i] = &response.Data.List[i]
			}
			return resp, userList, nil
		}
	}
}

/*
userNumber 传入一个数字，返回对应数量的用户列表的详细信息
userType 0 正式账号 1 测试账号 2 游客账号
vipLevel vip等级
*
*/
func GetUserVipListApi(ctx *context.Context, pageSize, userNumber int, userType int8, vipLevel int) (*model.Response, []*UserInfo, error) {
	api := "/api/Users/GetPageList"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &GetUserVipListStruct{}
	payloadList := []interface{}{vipLevel, userType, pageSize, userNumber, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList请求失败", err)), nil, err
	} else {
		// 提取用户id
		var response GetUserListResponseStruct
		if err := json.Unmarshal(respBoy, &response); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList提取用户列表解析失败", err)), nil, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList解析失败", err)), nil, err
		} else {
			if len(response.Data.List) == 0 || response.Data.List == nil {
				logger.Logger.Warn("获取用户列表为空")
				return resp, nil, nil
			}
			userList := make([]*UserInfo, 0, len(response.Data.List))
			for i := 0; i < len(response.Data.List); i++ {
				userList = append(userList, &response.Data.List[i])
			}
			return resp, userList, nil
		}
	}
}

type UserAmount struct {
	UserId int `json:"userId"`
	model.BaseStruct
}

type UserAmountResponse struct {
	Data          string `json:"data"`
	MsgParameters string `json:"msgParameters"`
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
	MsgCode       int    `json:"msgCode"`
}

// 用户id返回用户账号
func GetUserAmount(ctx *context.Context, userid int) (*model.Response, string, error) {
	api := "/api/Users/GetUserAccount"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &UserAmount{}
	payloadList := []interface{}{userid, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetUserAccount请求失败", err)), "", err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetUserAccount解析失败", err)), "", err
		} else {
			var res UserAmountResponse
			if err := json.Unmarshal(respBoy, &res); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetUserAccount[UserAmountResponse]解析失败", err)), "", err
			} else {
				// logger.Logger.Info("用户的账号", res.Data)
				return resp, res.Data, nil
			}
		}
	}
}

type UserinfoStruct struct {
	UserId int `json:"userId"`
	model.QueryPayloadStruct
}

// 根据id会员获得用户的详情  用户的类型
func GetUserDetail(ctx *context.Context, userid int) (*model.Response, *UserInfo, error) {
	api := "/api/Users/GetPageList"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &UserinfoStruct{}
	payloadList := []interface{}{userid, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList请求失败", err)), nil, err
	} else {
		// 提取用户id
		var response GetUserListResponseStruct
		if err := json.Unmarshal(respBoy, &response); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList提取用户列表解析失败[GetUserListResponseStruct]", err)), nil, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList解析失败", err)), nil, err
		} else {
			if len(response.Data.List) == 0 || response.Data.List == nil {
				logger.Logger.Warn("获取用户列表为空")
				return resp, nil, nil
			}
			return resp, &response.Data.List[0], nil
		}
	}
}
