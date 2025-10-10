package messagemanagement

const WorkOrderNumber = 22 // 工单的数量

// 工单的枚举
type workOrderNumber int8

const (
	Externallinks workOrderNumber = iota + 1 // 外部链接
	OneToone
)

// 工单系统的list
var WordOrderList = []string{"外部链接", "一对一客服", "其他问题", "存款未到账自动化", "取款未到账", "修改银行信息", "修改真实姓名半自动", "修改登录密码半自动", "忘记会员账号",
	"会员账号解冻半自动", "修改IFSC自动化", "修改银行名称自动化", "删除USDT半自动", "删除银行卡半自动", "删除PIX自动化", "删除电子钱包半自动", "新增USDT半自动",
	"删除银行卡自动化", "删除USDT自动化", "删除电子钱包自动化", "修改提现密码自动化", "修改提现密码半自动化"}
