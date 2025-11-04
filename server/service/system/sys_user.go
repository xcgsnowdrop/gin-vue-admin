package system

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"gmserver/model/common"
	systemReq "gmserver/model/system/request"

	"gmserver/global"
	"gmserver/model/system"
	"gmserver/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

type UserService struct{}

var UserServiceApp = new(UserService)

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.New()
	err = global.GVA_DB.Create(&u).Error
	return u, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.GVA_DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	err = global.GVA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
		MenuServiceApp.UserAuthorityDefaultRouter(&user)
	}
	return &user, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.SysUser, newPassword string
//@return: err error

func (userService *UserService) ChangePassword(u *system.SysUser, newPassword string) (err error) {
	var user system.SysUser
	err = global.GVA_DB.Select("id, password").Where("id = ?", u.ID).First(&user).Error
	if err != nil {
		return err
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return errors.New("原密码错误")
	}
	pwd := utils.BcryptHash(newPassword)
	err = global.GVA_DB.Model(&user).Update("password", pwd).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserInfoList(info systemReq.GetUserList, currentUserAuthorityId uint) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysUser{})
	var userList []system.SysUser

	if info.NickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+info.NickName+"%")
	}
	if info.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+info.Phone+"%")
	}
	if info.Username != "" {
		db = db.Where("username LIKE ?", "%"+info.Username+"%")
	}
	if info.Email != "" {
		db = db.Where("email LIKE ?", "%"+info.Email+"%")
	}

	// 权限过滤：只返回当前用户可以管理的用户（角色等级低于当前用户的用户）
	// 获取当前用户可以管理的所有子角色列表
	manageableAuthIDs, err := AuthorityServiceApp.GetStructAuthorityList(currentUserAuthorityId)
	if err != nil {
		return nil, 0, errors.New("获取角色权限列表失败")
	}

	// 如果可管理的角色列表为空，则不能查看任何其他用户（只能查看自己）
	if len(manageableAuthIDs) == 0 {
		// 只允许查看自己（通过其他方式过滤，这里不返回任何用户）
		db = db.Where("1 = 0") // 不返回任何结果
	} else {
		// 只返回角色在可管理范围内的用户
		db = db.Where("authority_id IN ?", manageableAuthIDs)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func (userService *UserService) SetUserAuthority(id uint, authorityId uint) (err error) {

	assignErr := global.GVA_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}

	var authority system.SysAuthority
	err = global.GVA_DB.Where("authority_id = ?", authorityId).First(&authority).Error
	if err != nil {
		return err
	}
	var authorityMenu []system.SysAuthorityMenu
	var authorityMenuIDs []string
	err = global.GVA_DB.Where("sys_authority_authority_id = ?", authorityId).Find(&authorityMenu).Error
	if err != nil {
		return err
	}

	for i := range authorityMenu {
		authorityMenuIDs = append(authorityMenuIDs, authorityMenu[i].MenuId)
	}

	var authorityMenus []system.SysBaseMenu
	err = global.GVA_DB.Preload("Parameters").Where("id in (?)", authorityMenuIDs).Find(&authorityMenus).Error
	if err != nil {
		return err
	}
	hasMenu := false
	for i := range authorityMenus {
		if authorityMenus[i].Name == authority.DefaultRouter {
			hasMenu = true
			break
		}
	}
	if !hasMenu {
		return errors.New("找不到默认路由,无法切换本角色")
	}

	err = global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", id).Update("authority_id", authorityId).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserAuthorities
//@description: 设置一个用户的权限
//@param: id uint, authorityIds []string
//@return: err error

func (userService *UserService) SetUserAuthorities(adminAuthorityID, id uint, authorityIds []uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var user system.SysUser
		TxErr := tx.Where("id = ?", id).First(&user).Error
		if TxErr != nil {
			global.GVA_LOG.Debug(TxErr.Error())
			return errors.New("查询用户数据失败")
		}

		// 权限检查1: 检查目标用户当前的角色是否是操作者角色的子角色（或同级/父级则不允许修改）
		// 只有目标用户的角色是操作者角色的子角色时，才允许修改
		// 规则：不能修改角色等级比自己高的用户，不能修改角色等级和自己同级的用户，只能修改角色等级比自己低的用户
		if user.AuthorityId != 0 {
			// 如果目标用户的角色就是操作者的角色，不允许修改（同级）
			if user.AuthorityId == adminAuthorityID {
				return errors.New("无权修改该用户的角色（目标用户角色与您的角色相同）")
			}

			// 获取操作者可以管理的所有子角色列表（递归获取所有子角色）
			manageableAuthIDs, err := AuthorityServiceApp.GetStructAuthorityList(adminAuthorityID)
			if err != nil {
				return errors.New("获取角色权限列表失败")
			}

			// 检查目标用户当前的角色是否在操作者可管理的范围内（即是否是子角色）
			// 如果不在，说明目标用户的角色等级 >= 操作者角色等级，不允许修改
			if !slices.Contains(manageableAuthIDs, user.AuthorityId) {
				return errors.New("无权修改该用户的角色（目标用户角色等级高于或等于您的角色等级）")
			}
		}

		// 权限检查2: 检查要设置的新角色是否在操作者角色的管理范围内
		if len(authorityIds) == 0 {
			return errors.New("至少需要设置一个角色")
		}

		// 检查每个要设置的角色
		for _, v := range authorityIds {
			e := AuthorityServiceApp.CheckAuthorityIDAuth(adminAuthorityID, v)
			if e != nil {
				return e
			}
		}

		TxErr = tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []system.SysUserAuthority
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, system.SysUserAuthority{
				SysUserId: id, SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = tx.Model(&user).Update("authority_id", authorityIds[0]).Error
		if TxErr != nil {
			return TxErr
		}
		// 返回 nil 提交事务
		return nil
	})
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteUser
//@description: 删除用户
//@param: id float64
//@return: err error

func (userService *UserService) DeleteUser(adminAuthorityID uint, targetUserID int) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var targetUser system.SysUser
		if err := tx.Where("id = ?", targetUserID).First(&targetUser).Error; err != nil {
			return errors.New("用户不存在")
		}

		// 权限检查：只能删除角色等级比自己低的用户
		if targetUser.AuthorityId != 0 {
			// 如果目标用户的角色就是操作者的角色，不允许删除（同级）
			if targetUser.AuthorityId == adminAuthorityID {
				return errors.New("无权删除该用户（目标用户角色与您的角色相同）")
			}

			// 获取操作者可以管理的所有子角色列表
			manageableAuthIDs, err := AuthorityServiceApp.GetStructAuthorityList(adminAuthorityID)
			if err != nil {
				return errors.New("获取角色权限列表失败")
			}

			// 检查目标用户的角色是否在操作者可管理的范围内
			if !slices.Contains(manageableAuthIDs, targetUser.AuthorityId) {
				return errors.New("无权删除该用户（目标用户角色等级高于或等于您的角色等级）")
			}
		}

		if err := tx.Where("id = ?", targetUserID).Delete(&system.SysUser{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", targetUserID).Error; err != nil {
			return err
		}
		return nil
	})
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser model.SysUser
//@return: err error, user model.SysUser

func (userService *UserService) SetUserInfo(adminAuthorityID uint, req system.SysUser) error {
	// 权限检查：只能修改角色等级比自己低的用户
	var targetUser system.SysUser
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&targetUser).Error; err != nil {
		return errors.New("用户不存在")
	}

	if targetUser.AuthorityId != 0 {
		// 如果目标用户的角色就是操作者的角色，不允许修改（同级）
		if targetUser.AuthorityId == adminAuthorityID {
			return errors.New("无权修改该用户（目标用户角色与您的角色相同）")
		}

		// 获取操作者可以管理的所有子角色列表
		manageableAuthIDs, err := AuthorityServiceApp.GetStructAuthorityList(adminAuthorityID)
		if err != nil {
			return errors.New("获取角色权限列表失败")
		}

		// 检查目标用户的角色是否在操作者可管理的范围内
		if !slices.Contains(manageableAuthIDs, targetUser.AuthorityId) {
			return errors.New("无权修改该用户（目标用户角色等级高于或等于您的角色等级）")
		}
	}

	return global.GVA_DB.Model(&system.SysUser{}).
		Select("updated_at", "nick_name", "header_img", "phone", "email", "enable").
		Where("id=?", req.ID).
		Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"nick_name":  req.NickName,
			"header_img": req.HeaderImg,
			"phone":      req.Phone,
			"email":      req.Email,
			"enable":     req.Enable,
		}).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetSelfInfo
//@description: 设置用户信息
//@param: reqUser model.SysUser
//@return: err error, user model.SysUser

func (userService *UserService) SetSelfInfo(req system.SysUser) error {
	return global.GVA_DB.Model(&system.SysUser{}).
		Where("id=?", req.ID).
		Updates(req).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetSelfSetting
//@description: 设置用户配置
//@param: req datatypes.JSON, uid uint
//@return: err error

func (userService *UserService) SetSelfSetting(req common.JSONMap, uid uint) error {
	return global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", uid).Update("origin_setting", req).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user system.SysUser

func (userService *UserService) GetUserInfo(uuid uuid.UUID) (user system.SysUser, err error) {
	var reqUser system.SysUser
	err = global.GVA_DB.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	MenuServiceApp.UserAuthorityDefaultRouter(&reqUser)
	return reqUser, err
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserById(id int) (user *system.SysUser, err error) {
	var u system.SysUser
	err = global.GVA_DB.Where("id = ?", id).First(&u).Error
	return &u, err
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserByUuid(uuid string) (user *system.SysUser, err error) {
	var u system.SysUser
	if err = global.GVA_DB.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("用户不存在")
	}
	return &u, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ResetPassword
//@description: 修改用户密码
//@param: ID uint
//@return: err error

func (userService *UserService) ResetPassword(adminAuthorityID uint, targetUserID uint, password string) (err error) {
	// 权限检查：只能重置角色等级比自己低的用户的密码
	var targetUser system.SysUser
	if err := global.GVA_DB.Where("id = ?", targetUserID).First(&targetUser).Error; err != nil {
		return errors.New("用户不存在")
	}

	if targetUser.AuthorityId != 0 {
		// 如果目标用户的角色就是操作者的角色，不允许重置（同级）
		if targetUser.AuthorityId == adminAuthorityID {
			return errors.New("无权重置该用户的密码（目标用户角色与您的角色相同）")
		}

		// 获取操作者可以管理的所有子角色列表
		manageableAuthIDs, err := AuthorityServiceApp.GetStructAuthorityList(adminAuthorityID)
		if err != nil {
			return errors.New("获取角色权限列表失败")
		}

		// 检查目标用户的角色是否在操作者可管理的范围内
		if !slices.Contains(manageableAuthIDs, targetUser.AuthorityId) {
			return errors.New("无权重置该用户的密码（目标用户角色等级高于或等于您的角色等级）")
		}
	}

	err = global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", targetUserID).Update("password", utils.BcryptHash(password)).Error
	return err
}
