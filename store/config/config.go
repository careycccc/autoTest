package config

import "time"

// 一些配置信息  sit环境
const (
	ADMIN_SYSTEM_URL      = "https://sit-tenantadmin-3003.mggametransit.com" // 后台地址,包括 domain refer
	REGISTER_URL          = "https://sit-3003-register.mggametransit.com"    // 前台注册地址
	SIT_WEB_API           = "https://sit-webapi.mggametransit.com"           // h5端的请求地址
	PLANT_H5              = "https://sit-plath5-y1.mggametransit.com"        // y1前台地址 包括domain refer
	WMG_H5                = "https://h5.wmgametransit.com"                   // y1彩票投注相关的请求头地址
	LOTTERY_H5            = "https://sit-lotteryh5.wmgametransit.com"        // y1彩票投注相关的请求体地址
	GoodsDeposit_URL      = "https://mggametransit.com/"                     // 充值相关的returenUrl链接
	TENANTID              = "3003"                                           // h5的商户id
	Log_Level             = "INFO"                                           // 设置日志等级
	MAXWaitTIME           = time.Second * 2                                  // 最大等待时间
	MAXRtryNUMBER         = 3                                                // 最大重试次数
	FIXEDTIME             = time.Second * 3                                  // 固定等待时间
	LANGUAGE              = "en"                                             // 系统语言
	ADMIN_UERNAME         = "carey3003"                                      // 后台商户账号
	ADMIN_PWD             = "qwer1234"                                       // 后台商户密码
	SUB_PWD               = "qwer1234"                                       // 后台修改的密码
	MIN_MONENY            = 5000                                             // 充值金额的最大值
	MAX_MONENY            = 10000                                            // 充值金额的最小值
	SUB_MINNUMBER         = 3                                                // 下级邀请人数的最小值
	SUB_MAXMUMBER         = 5                                                // 下级邀请人数的最大值
	SUB_CONCURRENT        = 3                                                // 邀请下级的并发数
	GeneralAgentNumber    = 2                                                // 邀请转盘的总代数量
	SUBUSERYAML           = "./subUser.yaml"                                 // 下级用户写入的目录
	WHEELNUMBER           = 1                                                // 活动的轮数，比如邀请转盘的就是1轮
	SiteMessageNumber     = 5                                                // 站内信发送个数
	SiteMessageConcurrent = 3                                                // 站内信并发
	WithdrawPassword      = "123456"                                         // 提现密码
	Remark                = "carey4399"                                      // 备注信息
)

// 一些配置信息  uat环境

// const (
// 	ADMIN_SYSTEM_URL      = "https://3101-tenantadmin.arplatsaasuat.com"   // 后台地址,包括domain refer
// 	REGISTER_URL          = "https://3101-register-uat.arplatsaasuat.com/" // 前台注册地址,包括domain refer
// 	SIT_WEB_API           = "https://api.arplatsaasuat.com"                // h5端的请求地址
// 	PLANT_H5              = "https://3101h5.arplatsaasuat.com"             // y1前台地址 包括domain refer
// 	WMG_H5                = "https://h5.wmgametransit.com"                 // y1彩票投注相关的请求头地址
// 	LOTTERY_H5            = "https://sit-lotteryh5.wmgametransit.com"
// 	Log_Level             = "INFO" // 设置日志登记
// 	TENANTID              = "3101"
// 	MAXWaitTIME           = time.Second * 3 // 最大等待时间
// 	MAXRtryNUMBER         = 3               // 最大重试次数
// 	FIXEDTIME             = time.Second * 3 // 固定等待时间
// 	LANGUAGE              = "en"
// 	ADMIN_UERNAME         = "carey_3101"     // 后台商户账号
// 	ADMIN_PWD             = "qwer1234"       // 后台商户密码
// 	SUB_PWD               = "qwer1234"       // 后台修改的密码
// 	MIN_MONENY            = 1000             // 充值金额的最大值
// 	MAX_MONENY            = 10000            // 充值金额的最小值
// 	SUB_MINNUMBER         = 1                // 下级邀请人数的最小值
// 	SUB_MAXMUMBER         = 2                // 下级邀请人数的最大值
// 	SUB_CONCURRENT        = 3                // 邀请下级的并发数
// 	GeneralAgentNumber    = 2                // 邀请转盘的总代数量
// 	SUBUSERYAML           = "./subUser.yaml" // 下级用户写入的目录
// 	WHEELNUMBER           = 1                // 活动的轮数，比如邀请转盘的就是1轮
// 	SiteMessageNumber     = 5                // 站内信发送个数
// 	SiteMessageConcurrent = 3                // 站内信并发
// 	WithdrawPassword      = "123456"         // 提现密码
// 	Remark                = "carey4399"      // 备注信息
// )
