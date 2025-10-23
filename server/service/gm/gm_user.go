package gm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gm"
	gmReq "github.com/flipped-aurora/gin-vue-admin/server/model/gm/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

type UserService struct{}

// GetUser 获取游戏用户详情
func (userService *UserService) GetUser(id string) (user gm.GMUser, err error) {
	var userId uint64
	userId, err = strconv.ParseUint(id, 10, 32)
	if err != nil {
		return user, errors.New("无效的用户ID")
	}

	err = global.GVA_DB.Where("id = ?", uint(userId)).First(&user).Error
	return user, err
}

// CreateUser 创建游戏用户
func (userService *UserService) CreateUser(user gmReq.CreateUserRequest) (err error) {
	// 检查游戏用户ID是否已存在
	var existingUser gm.GMUser
	err = global.GVA_DB.Where("game_user_id = ?", user.GameUserId).First(&existingUser).Error
	if err == nil {
		return errors.New("游戏用户ID已存在")
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// 检查用户名是否已存在
	err = global.GVA_DB.Where("user_name = ?", user.UserName).First(&existingUser).Error
	if err == nil {
		return errors.New("用户名已存在")
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// 创建新用户
	newUser := gm.GMUser{
		GameUserId:  user.GameUserId,
		UserName:    user.UserName,
		NickName:    user.NickName,
		Password:    user.Password, // 注意：实际项目中应该加密密码
		Phone:       user.Phone,
		Email:       user.Email,
		HeaderImg:   user.HeaderImg,
		Enable:      user.Enable,
		AuthorityId: user.AuthorityId,
	}

	err = global.GVA_DB.Create(&newUser).Error
	return err
}

// UpdateUser 更新游戏用户信息
func (userService *UserService) UpdateUser(user gmReq.UpdateUserRequest) (err error) {
	var existingUser gm.GMUser
	err = global.GVA_DB.Where("id = ?", user.ID).First(&existingUser).Error
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查游戏用户ID是否被其他用户使用
	if user.GameUserId != "" && user.GameUserId != existingUser.GameUserId {
		var conflictUser gm.GMUser
		err = global.GVA_DB.Where("game_user_id = ? AND id != ?", user.GameUserId, user.ID).First(&conflictUser).Error
		if err == nil {
			return errors.New("游戏用户ID已被其他用户使用")
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	// 检查用户名是否被其他用户使用
	if user.UserName != "" && user.UserName != existingUser.UserName {
		var conflictUser gm.GMUser
		err = global.GVA_DB.Where("user_name = ? AND id != ?", user.UserName, user.ID).First(&conflictUser).Error
		if err == nil {
			return errors.New("用户名已被其他用户使用")
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	// 更新用户信息
	updateData := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if user.GameUserId != "" {
		updateData["game_user_id"] = user.GameUserId
	}
	if user.UserName != "" {
		updateData["user_name"] = user.UserName
	}
	if user.NickName != "" {
		updateData["nick_name"] = user.NickName
	}
	if user.Phone != "" {
		updateData["phone"] = user.Phone
	}
	if user.Email != "" {
		updateData["email"] = user.Email
	}
	if user.HeaderImg != "" {
		updateData["header_img"] = user.HeaderImg
	}
	if user.Enable != 0 {
		updateData["enable"] = user.Enable
	}
	if user.AuthorityId != "" {
		updateData["authority_id"] = user.AuthorityId
	}

	err = global.GVA_DB.Model(&existingUser).Updates(updateData).Error
	return err
}

// DeleteUser 删除游戏用户
func (userService *UserService) DeleteUser(id uint) (err error) {
	var user gm.GMUser
	err = global.GVA_DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return errors.New("用户不存在")
	}

	err = global.GVA_DB.Delete(&user).Error
	return err
}

// ResetPassword 重置游戏用户密码
func (userService *UserService) ResetPassword(id uint, password string) (err error) {
	var user gm.GMUser
	err = global.GVA_DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return errors.New("用户不存在")
	}

	// 注意：实际项目中应该加密密码
	err = global.GVA_DB.Model(&user).Update("password", password).Error
	return err
}

// SetUserStatus 设置游戏用户状态
func (userService *UserService) SetUserStatus(id uint, status int) (err error) {
	var user gm.GMUser
	err = global.GVA_DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return errors.New("用户不存在")
	}

	err = global.GVA_DB.Model(&user).Update("enable", status).Error
	return err
}

// BatchOperate 批量操作游戏用户
func (userService *UserService) BatchOperate(operation string, userIds []uint) (err error) {
	switch operation {
	case "enable":
		err = global.GVA_DB.Model(&gm.GMUser{}).Where("id IN ?", userIds).Update("enable", 1).Error
	case "disable":
		err = global.GVA_DB.Model(&gm.GMUser{}).Where("id IN ?", userIds).Update("enable", 2).Error
	case "delete":
		err = global.GVA_DB.Where("id IN ?", userIds).Delete(&gm.GMUser{}).Error
	default:
		return errors.New("不支持的操作类型")
	}
	return err
}

// ExportUsers 导出游戏用户数据
func (userService *UserService) ExportUsers(info gmReq.SearchUserParams) (filePath string, err error) {
	// 这里应该实现Excel导出逻辑
	// 暂时返回模拟文件路径
	filePath = fmt.Sprintf("exports/gm_users_%d.xlsx", time.Now().Unix())
	return filePath, nil
}

// GetUserStats 获取游戏用户统计信息
func (userService *UserService) GetUserStats() (stats map[string]interface{}, err error) {
	var totalUsers, activeUsers, newUsersToday int64

	// 总用户数
	err = global.GVA_DB.Model(&gm.GMUser{}).Count(&totalUsers).Error
	if err != nil {
		return nil, err
	}

	// 活跃用户数（启用状态）
	err = global.GVA_DB.Model(&gm.GMUser{}).Where("enable = ?", 1).Count(&activeUsers).Error
	if err != nil {
		return nil, err
	}

	// 今日新增用户数
	today := time.Now().Format("2006-01-02")
	err = global.GVA_DB.Model(&gm.GMUser{}).Where("DATE(created_at) = ?", today).Count(&newUsersToday).Error
	if err != nil {
		return nil, err
	}

	stats = map[string]interface{}{
		"totalUsers":    totalUsers,
		"activeUsers":   activeUsers,
		"newUsersToday": newUsersToday,
		"onlineUsers":   0, // 在线用户数需要根据实际业务逻辑计算
	}

	return stats, nil
}

// GetPlayerListFromMongo 从 MongoDB 获取玩家列表（示例方法）
func (userService *UserService) GetPlayerListFromMongo(info gmReq.SearchUserParams) (list []*gm.PlayerInfoWithStatus, total int64, err error) {
	ctx := context.Background()

	// 构建查询条件
	filter := bson.M{}
	if info.UserId != "" {
		filter["user_id"] = info.UserId
	}
	if info.PlayerId != "" {
		filter["player_id"] = info.PlayerId
	}
	if info.UniqueId != "" {
		filter["unique_id"] = info.UniqueId
	}
	if info.NickName != "" {
		filter["nickname"] = info.NickName
	}

	// 构建分页参数
	page := pagination.PageParam{
		PageNum:  int64(info.Page),
		PageSize: int64(info.PageSize),
	}

	// 构建查询选项
	options := &pagination.FindOptions{
		Sort: bson.D{}, // 排序条件
		Projection: bson.M{
			"_id":           0, // 排除 _id 字段
			"user_id":       1,
			"player_id":     1,
			"nickname":      1,
			"lv":            1,
			"power":         1,
			"area_id":       1,
			"unique_id":     1,
			"login_time":    1,
			"register_time": 1,
		},
	}

	// 执行分页查询（使用新的方法）
	pageRes, err := pagination.FindWithPageOptions(ctx, global.GVA_MONGO_PLAYER_INFO, filter, page, func() *gm.PlayerInfo { return &gm.PlayerInfo{} }, options)
	if err != nil {
		return nil, 0, err
	}

	// 批量获取状态（性能优化）
	playerIds := make([]string, 0, len(pageRes.List))
	for _, item := range pageRes.List {
		playerIds = append(playerIds, item.PlayerId)
	}

	onlineMap, banLoginMap, banChatMap := userService.getPlayersStatusBatch(playerIds)

	// 转换为带状态的结构体并设置状态
	resultList := make([]*gm.PlayerInfoWithStatus, 0, len(pageRes.List))
	for _, item := range pageRes.List {
		playerWithStatus := item.ToPlayerInfoWithStatus()

		// 从批量获取的map中设置状态
		playerWithStatus.Online = onlineMap[item.PlayerId]
		playerWithStatus.BanLogin = banLoginMap[item.PlayerId]
		playerWithStatus.BanChat = banChatMap[item.PlayerId]

		resultList = append(resultList, playerWithStatus)
	}

	return resultList, pageRes.Total, nil
}

// getPlayerOnlineStatus 获取玩家在线状态
func (userService *UserService) getPlayerOnlineStatus(playerId string) bool {
	// TODO: 从Redis或在线用户列表获取在线状态
	// 示例实现：
	// return global.GVA_REDIS.SIsMember(ctx, "online_players", playerId).Val()
	return false // 默认离线
}

// getPlayerBanLoginStatus 获取玩家封号状态
func (userService *UserService) getPlayerBanLoginStatus(playerId string) bool {
	// TODO: 从封号列表获取封号状态
	// 示例实现：
	// return global.GVA_REDIS.SIsMember(ctx, "ban_login_players", playerId).Val()
	return false // 默认未封号
}

// getPlayerBanChatStatus 获取玩家禁言状态
func (userService *UserService) getPlayerBanChatStatus(playerId string) bool {
	// TODO: 从禁言列表获取禁言状态
	// 示例实现：
	// return global.GVA_REDIS.SIsMember(ctx, "ban_chat_players", playerId).Val()
	return false // 默认未禁言
}

// getPlayersStatusBatch 批量获取玩家状态（性能优化）
func (userService *UserService) getPlayersStatusBatch(playerIds []string) (onlineMap, banLoginMap, banChatMap map[string]bool) {

	// 初始化结果map
	onlineMap = make(map[string]bool)
	banLoginMap = make(map[string]bool)
	banChatMap = make(map[string]bool)

	// TODO: 批量从Redis获取状态
	// 示例实现：
	// onlinePlayers := global.GVA_REDIS.SMembers(ctx, "online_players").Val()
	// banLoginPlayers := global.GVA_REDIS.SMembers(ctx, "ban_login_players").Val()
	// banChatPlayers := global.GVA_REDIS.SMembers(ctx, "ban_chat_players").Val()

	// 构建状态map
	// for _, playerId := range onlinePlayers {
	//     onlineMap[playerId] = true
	// }
	// for _, playerId := range banLoginPlayers {
	//     banLoginMap[playerId] = true
	// }
	// for _, playerId := range banChatPlayers {
	//     banChatMap[playerId] = true
	// }

	// 设置默认值
	for _, playerId := range playerIds {
		onlineMap[playerId] = false
		banLoginMap[playerId] = false
		banChatMap[playerId] = false
	}

	return onlineMap, banLoginMap, banChatMap
}

// CreatePlayerInMongo 在 MongoDB 中创建玩家（示例方法）
func (userService *UserService) CreatePlayerInMongo(player gm.GMUser) error {
	ctx := context.Background()

	// 检查游戏用户ID是否已存在
	filter := bson.M{"game_user_id": player.GameUserId}
	count, err := global.GVA_MONGO_PLAYER_INFO.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("游戏用户ID已存在")
	}

	// 检查用户名是否已存在
	filter = bson.M{"user_name": player.UserName}
	count, err = global.GVA_MONGO_PLAYER_INFO.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 创建新玩家
	_, err = global.GVA_MONGO_PLAYER_INFO.InsertOne(ctx, player)
	return err
}
