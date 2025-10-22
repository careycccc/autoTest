package invitecode

import (
	financialmanagement "autoTest/API/adminApi/financialManagement"
	"autoTest/API/adminApi/login"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	lotterygameapi "autoTest/API/betApi/LotteryGameApi"
	registerapi "autoTest/API/deskApi/registerApi"
	"autoTest/API/utils"
	"autoTest/store/logger"
	util "autoTest/store/utils"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// User 定义用户结构
type User struct {
	InviteCode   string   // 用户的邀请码
	Subordinates []string // 直接下级的邀请码列表
}

// 定义用户详细信息的全局变量
type UserDetail struct {
	ctx         *context.Context //保存用户的上下文
	userAccount string           // 保存用户的账号（手机号或邮箱）
}

// 全局用户详情映射
var userDetails = make([]*UserDetail, 0)
var userDB = make(map[string]*User)
var dbMutex sync.Mutex

// ProcessNewUser 处理新用户的逻辑,充值和投注
func ProcessNewUser(inviteCode string) error {
	if len(userDetails) == 0 {
		logger.LogError("用户详情列表为空，无法处理新用户后续的请求", nil)
		return errors.New("用户详情列表为空，无法处理新用户后续的请求")
	}
	ctx, err := login.RunAdminSitLogin()
	if err != nil {
		logger.LogError("商户后台登录失败", err)
		return fmt.Errorf("商户后台登录失败: %w", err)
	}
	ctxAdmin := ctx
	// 查找对应的用户详情
	for _, userDetail := range userDetails {
		if _, userId, err := memberlist.GetUserIdApi(ctxAdmin, userDetail.userAccount); err != nil {
			logger.LogError("获取用户ID失败", err)
			continue
		} else {
			// 充值
			moneny, _ := util.GenerateRandomInt(2000, 5000)
			if _, err := financialmanagement.ArtificialRechargeFunc(ctxAdmin, int(userId), moneny, 1); err != nil {
				logger.LogError("绑定层级的后台人工充值失败", err)
				continue
			} else {
				// 上面的充值结束后台，在进行投注
				time.Sleep(time.Second)
				// 先把投注的结果随机出来
				gameCode, betContent, amount, betMultiple := lotterygameapi.GetBetResult()
				if err := lotterygameapi.RunBetFunc(userDetail.ctx, gameCode, betContent, userDetail.userAccount, amount, betMultiple); err != nil {
					continue
				}
			}
		}
	}
	return nil
}

// ========== 🚀 精确版：每层完整绑定 ==========

// BindOneLevel 绑定**完整的一层**（你的精确逻辑）
func BindOneLevel(ctx *context.Context, parentInviteCodes []string, mobileCount int, level int, newUsers *[]string) error {
	if len(parentInviteCodes) == 0 {
		return errors.New("父级邀请码列表不能为空")
	}

	fmt.Printf("🔗 [层级%d] 父级%d人 -> 生成%d个下级...\n", level+1, len(parentInviteCodes), mobileCount)

	// Step 1: **一次性**生成当前层所有手机号
	mobileIds := utils.RandmoUserId(mobileCount)
	if len(mobileIds) != mobileCount {
		return fmt.Errorf("生成手机号数量不匹配: 期望%d, 实际%d", mobileCount, len(mobileIds))
	}
	fmt.Printf("📱 生成手机号: %v...\n", mobileIds[:min(2, len(mobileIds))])

	// Step 2: **同时**注册所有用户（并发）
	var wg sync.WaitGroup
	newInviteCodes := make([]string, 0, mobileCount)
	mu := sync.Mutex{}
	errChan := make(chan error, mobileCount)

	for i, mobile := range mobileIds {
		wg.Add(1)
		func(mobile string, parentIndex int) {
			defer wg.Done()

			// 随机选择一个父级邀请码
			parentInviteCode := parentInviteCodes[rand.Intn(len(parentInviteCodes))]

			// 🔥 **同时**发起注册（你的核心逻辑）
			_, subCtx, err := registerapi.RegisterMobileLoginFunc(mobile, parentInviteCode)
			if err != nil {
				errChan <- fmt.Errorf("注册失败 %s -> %s: %w", mobile, parentInviteCode, err)
				return
			}

			time.Sleep(time.Second) // 避免请求过快
			// 把下级用户的信息保存起来
			userDetail := &UserDetail{
				ctx:         &subCtx,
				userAccount: mobile,
			}
			userDetails = append(userDetails, userDetail)

			// 获取新用户的邀请码
			_, newInviteCode, err := memberlist.GetUserInviteCodeApi(ctx, mobile)
			if err != nil {
				errChan <- fmt.Errorf("获取邀请码失败 %s: %w", mobile, err)
				return
			}

			if newInviteCode == "" {
				errChan <- fmt.Errorf("新用户 %s 邀请码为空", mobile)
				return
			}

			mu.Lock()
			newInviteCodes = append(newInviteCodes, newInviteCode)
			mu.Unlock()

			fmt.Printf("✅ [%d/%d] %s -> %s (邀请码: %s)\n", i+1, mobileCount, mobile, parentInviteCode, newInviteCode)
		}(mobile, i)
	}

	wg.Wait()
	close(errChan)

	// 检查错误
	for err := range errChan {
		return err
	}

	if len(newInviteCodes) != mobileCount {
		return fmt.Errorf("绑定数量不匹配: 期望%d, 实际%d", mobileCount, len(newInviteCodes))
	}

	// Step 3: 更新数据库（批量）
	dbMutex.Lock()
	for _, parentCode := range parentInviteCodes {
		if _, exists := userDB[parentCode]; !exists {
			userDB[parentCode] = &User{InviteCode: parentCode, Subordinates: []string{}}
		}
		// 将新用户分配给父级（按注册时随机分配）
	}

	// 简化：记录所有新用户
	for _, code := range newInviteCodes {
		userDB[code] = &User{InviteCode: code, Subordinates: []string{}}
	}
	dbMutex.Unlock()

	// 记录结果
	*newUsers = append(*newUsers, newInviteCodes...)
	fmt.Printf("🎯 [层级%d] 完成: %d人 -> %v...\n", level+1, len(newInviteCodes), newInviteCodes[:min(2, len(newInviteCodes))])

	return nil
}

// ========== 完整多层级执行 ==========

// RunAAWithBB 精确的多层级绑定（只需执行一次ProcessNewUser）
func RunAAWithBB(ctx *context.Context, rootInviteCode string, subordinates []int) error {
	rand.Seed(time.Now().UnixNano())
	layers := make([][]string, len(subordinates))

	if len(subordinates) == 0 {
		return errors.New("层级人数切片不能为空")
	}

	// 🔥 第1层：绑定到总代
	fmt.Printf("🚀 === 开始第1层绑定到总代 [%s] (%d人) ===\n", rootInviteCode, subordinates[0])
	newUsersForThisCall := []string{}
	if err := BindOneLevel(ctx, []string{rootInviteCode}, subordinates[0], 0, &newUsersForThisCall); err != nil {
		return fmt.Errorf("第一层绑定失败: %w", err)
	}
	layers[0] = newUsersForThisCall
	fmt.Printf("✅ 第1层完成: %d人 -> %v\n", len(layers[0]), layers[0][:min(2, len(layers[0]))])

	// 🔥 后续层级
	currentParentCodes := layers[0]
	for level := 1; level < len(subordinates); level++ {
		mobileCount := subordinates[level]
		fmt.Printf("\n🚀 === 开始第%d层绑定 (%d人) ===\n", level+1, mobileCount)

		newUsersForThisCall := []string{}
		if err := BindOneLevel(ctx, currentParentCodes, mobileCount, level, &newUsersForThisCall); err != nil {
			return fmt.Errorf("第%d层绑定失败: %w", level+1, err)
		}
		layers[level] = newUsersForThisCall
		currentParentCodes = newUsersForThisCall
		fmt.Printf("✅ 第%d层完成: %d人 -> %v\n", level+1, len(layers[level]), layers[level][:min(2, len(layers[level]))])
	}

	// 🔥 **只需执行一次** ProcessNewUser（处理所有新用户）
	newUsers := []string{}
	for _, layer := range layers {
		newUsers = append(newUsers, layer...)
	}

	// 🔥 关键修正：只需调用一次！
	if err := ProcessNewUser(""); err != nil { // 空参数表示处理所有用户
		return fmt.Errorf("充值投注失败: %w", err)
	}

	return nil
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// func printHierarchy(inviteCode string) {
// 	fmt.Println("\n📊 === 最终层级结构 ===")
// 	dbMutex.Lock()
// 	fmt.Printf("👤 总代 %s -> [%d个直属]\n", inviteCode, len(userDB[inviteCode].Subordinates))
// 	for code, user := range userDB {
// 		if code != inviteCode && len(user.Subordinates) > 0 {
// 			fmt.Printf("👤 %s -> [%d个下级]\n", code, len(user.Subordinates))
// 		}
// 	}
// 	dbMutex.Unlock()
// }

// RunInvite 一键执行
func RunInvite() {
	inviteCode := "Z88JQQN" // 2441371
	ctx, err := login.RunAdminSitLogin()
	if err != nil {
		fmt.Println("❌ 登录失败:", err)
		return
	}

	userDB = make(map[string]*User)
	subordinates := []int{8, 5, 7, 3} // 第1层2人，第2层3人
	fmt.Printf("🎯 开始绑定到总代: %s, 层级: %v\n", inviteCode, subordinates)

	err = RunAAWithBB(ctx, inviteCode, subordinates)
	if err != nil {
		fmt.Println("❌ 执行失败:", err)
	} else {
		fmt.Println("\n🎉 执行成功!")
		//printHierarchy(inviteCode)
	}
}
