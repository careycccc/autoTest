# AutoTest - 游戏平台自动化测试框架

基于 Go 语言开发的游戏平台自动化测试框架，支持多租户环境的接口测试。

## 项目简介

本项目是一个全面的自动化测试解决方案，涵盖了游戏平台的前台和后台功能测试，包括用户注册、登录、充值、提现、游戏投注、活动管理等核心业务流程。

## 技术栈

- Go 1.25.1
- Zap (日志框架)
- Allure (测试报告)
- Excelize (Excel 处理)
- WebSocket (实时通信)

## 项目结构

```
.
├── API/                          # API 接口封装
│   ├── adminApi/                 # 后台管理 API
│   │   ├── activeManagement/     # 活动管理
│   │   ├── financialManagement/  # 财务管理
│   │   ├── GameManagement/       # 游戏管理
│   │   ├── memberList/           # 会员列表
│   │   ├── operationsManagement/ # 运营管理
│   │   └── reportManagement/     # 报表管理
│   ├── deskApi/                  # 前台用户 API
│   │   ├── active/               # 活动相关
│   │   ├── loginApi/             # 登录注册
│   │   ├── registerApi/          # 用户注册
│   │   ├── topUp/                # 充值
│   │   ├── WithdrawCash/         # 提现
│   │   └── invitationCarousel/   # 邀请转盘
│   ├── betApi/                   # 投注 API
│   │   ├── LotteryGameApi/       # 彩票游戏
│   │   └── threeGameApi/         # 三方游戏
│   └── ReportProcessingCenter/   # 报表处理中心
├── store/                        # 核心工具库
│   ├── config/                   # 配置管理
│   ├── logger/                   # 日志管理
│   ├── model/                    # 数据模型
│   ├── request/                  # 请求封装
│   └── utils/                    # 工具函数
├── PressureMeasurementModule/    # 性能测试模块
│   ├── boomer/                   # Boomer 压测
│   ├── f1/                       # F1 压测
│   └── accounts/                 # 测试账号
├── k6_load/                      # K6 性能测试
├── testcase/                     # 测试用例
├── requstModle/                  # 请求模型
└── main.go                       # 程序入口
```

## 功能模块

### 前台功能 (deskApi)
- 用户注册（邮箱注册、邀请码注册）
- 用户登录
- 充值管理
- 提现管理
- 活动参与（邀请转盘、充值转盘、每日签到）
- VIP 系统
- 用户信息管理

### 后台功能 (adminApi)
- 会员管理（会员列表、黑名单、状态管理）
- 财务管理（人工充值、充值订单、提现订单）
- 活动管理（签到活动、提现超时赔付）
- 游戏管理（投注订单查询）
- 运营管理（站内信、工单管理）
- 报表管理（平台报表、会员报表、资金流水）

### 游戏功能 (betApi)
- 彩票游戏投注
- 三方游戏接入
- 鸡路小游戏


## 环境配置

项目支持多环境配置，在 `store/config/config.go` 中切换：

- SIT 3001 环境
- SIT 3002 环境
- SIT 3003 环境
- SIT 3004 环境（当前）
- UAT 3101 环境

### 配置项说明

```go
ADMIN_SYSTEM_URL      // 后台管理地址
REGISTER_URL          // 前台注册地址
SIT_WEB_API           // H5 API 地址
TENANTID              // 租户 ID
ADMIN_UERNAME         // 后台账号
ADMIN_PWD             // 后台密码
MIN_MONENY / MAX_MONENY  // 充值金额范围
SUB_MINNUMBER / SUB_MAXMUMBER  // 下级邀请人数范围
```

## 快速开始

### 安装依赖

```bash
go mod download
```

### 运行测试

```bash
# 编译
go build -o main main.go

# 运行
./main
```

### 主要测试场景

在 `main.go` 中取消注释相应的函数即可运行：

```go
// 邀请转盘测试
invitationcarousel.RunSpinInvitedWheelWork()

// 邮箱注册测试
registerapi.RunEmailregeister()

// 提现测试
withdrawcash.RunWithDrawCase()

// 充值测试
topup.RunRechargeGoods()

// 游戏投注订单查询
GameBetOrders.RunGameBetOrders()

// 每日签到活动
dailycheckin.RunDailyCheckInActivity()

// 人工充值
financialmanagement.RunArtificialRechargeFunc()
```

## 日志管理

项目使用 Zap 日志框架，支持：
- 控制台输出
- 文件输出（可选）
- 日志级别控制
- 日志轮转

初始化日志：
```go
logger.InitLogger()  // 控制台输出
logger.InitLogger2() // 文件输出
```

## 测试报告

支持 Allure 测试报告生成，详见 `allure报告的生成.txt`



## 注意事项

1. 运行前需要在 `store/config/config.go` 中配置正确的环境
2. 确保测试环境网络可达
3. 部分功能需要后台管理员权限
4. 性能测试前请确认测试账号充足

## 开发规范

- 结构体字段顺序必须与 payloadList 切片顺序一致
- 所有 API 请求需要包含签名验证
- 错误处理使用统一的 logger.LogError
- 测试数据使用随机生成，避免硬编码

## 贡献指南

1. Fork 本仓库
2. 创建特性分支
3. 提交变更
4. 推送到分支
5. 创建 Pull Request

## 许可证

内部项目，仅供团队使用
