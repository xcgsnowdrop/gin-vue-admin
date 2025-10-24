# Gin-Vue-Admin 权限管理系统详解

## 📋 目录
- [系统概述](#系统概述)
- [权限模型设计](#权限模型设计)
- [角色层次结构](#角色层次结构)
- [权限继承机制](#权限继承机制)
- [关键代码实现](#关键代码实现)
- [配置说明](#配置说明)
- [使用示例](#使用示例)
- [最佳实践](#最佳实践)

## 系统概述

Gin-Vue-Admin采用基于RBAC（Role-Based Access Control）的权限管理模型，支持角色层次结构和权限继承，确保系统安全性和权限管理的灵活性。

### 核心特性
- ✅ 角色层次结构（父子角色关系）
- ✅ 权限继承机制（子角色权限是父角色的子集）
- ✅ 菜单权限控制
- ✅ 数据权限隔离
- ✅ API权限控制（基于Casbin）
- ✅ 按钮级权限控制

## 权限模型设计

### 数据模型结构

#### 1. 角色模型 (SysAuthority)
```go
type SysAuthority struct {
    CreatedAt       time.Time       // 创建时间
    UpdatedAt       time.Time       // 更新时间
    DeletedAt       *time.Time      `sql:"index"`
    AuthorityId     uint            `json:"authorityId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"`
    AuthorityName   string          `json:"authorityName" gorm:"comment:角色名"`
    ParentId        *uint           `json:"parentId" gorm:"comment:父角色ID"`
    DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;"`
    Children        []SysAuthority  `json:"children" gorm:"-"`
    SysBaseMenus    []SysBaseMenu   `json:"menus" gorm:"many2many:sys_authority_menus;"`
    Users           []SysUser       `json:"-" gorm:"many2many:sys_user_authority;"`
    DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"`
}
```

#### 2. 菜单模型 (SysBaseMenu)
```go
type SysBaseMenu struct {
    global.GVA_MODEL
    MenuLevel     uint
    ParentId      uint                                       `json:"parentId" gorm:"comment:父菜单ID"`
    Path          string                                     `json:"path" gorm:"comment:路由path"`
    Name          string                                     `json:"name" gorm:"comment:路由name"`
    Hidden        bool                                       `json:"hidden" gorm:"comment:是否在列表隐藏"`
    Component     string                                     `json:"component" gorm:"comment:对应前端文件路径"`
    Sort          int                                        `json:"sort" gorm:"comment:排序标记"`
    Meta          `json:"meta" gorm:"embedded;comment:附加属性"`
    SysAuthoritys []SysAuthority                             `json:"authoritys" gorm:"many2many:sys_authority_menus;"`
    Children      []SysBaseMenu                              `json:"children" gorm:"-"`
    Parameters    []SysBaseMenuParameter                     `json:"parameters"`
    MenuBtn       []SysBaseMenuBtn                           `json:"menuBtn"`
}
```

#### 3. 用户模型 (SysUser)
```go
type SysUser struct {
    global.GVA_MODEL
    UUID          uuid.UUID      `json:"uuid" gorm:"index;comment:用户UUID"`
    Username      string         `json:"userName" gorm:"index;comment:用户登录名"`
    Password      string         `json:"-"  gorm:"comment:用户登录密码"`
    NickName      string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`
    HeaderImg     string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"`
    AuthorityId   uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`
    Authority     SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
    Authorities   []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
    Phone         string         `json:"phone"  gorm:"comment:用户手机号"`
    Email         string         `json:"email"  gorm:"comment:用户邮箱"`
    Enable        int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`
}
```

## 角色层次结构

### 层次关系设计

gin-vue-admin支持多级角色层次结构，每个角色可以有父角色和子角色：

```
超级管理员 (Super Admin)
├── 系统管理员 (System Admin)
│   ├── 用户管理员 (User Manager)
│   └── 内容管理员 (Content Manager)
├── 部门管理员 (Department Admin)
│   ├── 财务管理员 (Finance Manager)
│   └── 人事管理员 (HR Manager)
└── 普通用户 (Regular User)
    ├── 高级用户 (Premium User)
    └── 基础用户 (Basic User)
```

### 初始化数据示例

```go
// 从 server/source/system/authority.go
entities := []sysModel.SysAuthority{
    {AuthorityId: 888, AuthorityName: "普通用户", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
    {AuthorityId: 9528, AuthorityName: "测试角色", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
    {AuthorityId: 8881, AuthorityName: "普通用户子角色", ParentId: utils.Pointer[uint](888), DefaultRouter: "dashboard"},
}
```

## 权限继承机制

### 核心原则：子角色权限是父角色权限的"子集"

#### 1. 权限继承验证代码

```go
// server/service/system/sys_authority.go
func (authorityService *AuthorityService) CheckAuthorityIDAuth(authorityID, targetID uint) (err error) {
    if !global.GVA_CONFIG.System.UseStrictAuth {
        return nil
    }
    authIDS, err := authorityService.GetStructAuthorityList(authorityID)
    if err != nil {
        return err
    }
    hasAuth := false
    for _, v := range authIDS {
        if v == targetID {
            hasAuth = true
            break
        }
    }
    if !hasAuth {
        return errors.New("您提交的角色ID不合法")
    }
    return nil
}
```

#### 2. 菜单权限限制

```go
// server/service/system/sys_menu.go
func (menuService *MenuService) AddMenuAuthority(menus []system.SysBaseMenu, adminAuthorityID, authorityId uint) (err error) {
    // 当开启了严格的树角色并且父角色不为0时需要进行菜单筛选
    if global.GVA_CONFIG.System.UseStrictAuth && *authority.ParentId != 0 {
        var authorityMenus []system.SysAuthorityMenu
        err = global.GVA_DB.Where("sys_authority_authority_id = ?", adminAuthorityID).Find(&authorityMenus).Error
        if err != nil {
            return err
        }
        for i := range authorityMenus {
            menuIds = append(menuIds, authorityMenus[i].MenuId)
        }

        for i := range menus {
            hasMenu := false
            for j := range menuIds {
                idStr := strconv.Itoa(int(menus[i].ID))
                if idStr == menuIds[j] {
                    hasMenu = true
                }
            }
            if !hasMenu {
                return errors.New("添加失败,请勿跨级操作")
            }
        }
    }

    err = AuthorityServiceApp.SetMenuAuthority(&auth)
    return err
}
```

#### 3. 获取父角色权限范围

```go
// server/service/system/sys_menu.go
func (menuService *MenuService) getBaseMenuTreeMap(authorityID uint) (treeMap map[uint][]system.SysBaseMenu, err error) {
    parentAuthorityID, err := AuthorityServiceApp.GetParentAuthorityID(authorityID)
    if err != nil {
        return nil, err
    }

    var allMenus []system.SysBaseMenu
    treeMap = make(map[uint][]system.SysBaseMenu)
    db := global.GVA_DB.Order("sort").Preload("MenuBtn").Preload("Parameters")

    // 当开启了严格的树角色并且父角色不为0时需要进行菜单筛选
    if global.GVA_CONFIG.System.UseStrictAuth && parentAuthorityID != 0 {
        var authorityMenus []system.SysAuthorityMenu
        err = global.GVA_DB.Where("sys_authority_authority_id = ?", authorityID).Find(&authorityMenus).Error
        if err != nil {
            return nil, err
        }
        var menuIds []string
        for i := range authorityMenus {
            menuIds = append(menuIds, authorityMenus[i].MenuId)
        }
        db = db.Where("id in (?)", menuIds)
    }

    err = db.Find(&allMenus).Error
    for _, v := range allMenus {
        treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
    }
    return treeMap, err
}
```

## 关键代码实现

### 1. 角色创建

```go
// server/service/system/sys_authority.go
func (authorityService *AuthorityService) CreateAuthority(auth system.SysAuthority) (authority system.SysAuthority, err error) {
    if err = global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&system.SysAuthority{}).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
        return auth, ErrRoleExistence
    }

    e := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
        if err = tx.Create(&auth).Error; err != nil {
            return err
        }

        auth.SysBaseMenus = systemReq.DefaultMenu()
        if err = tx.Model(&auth).Association("SysBaseMenus").Replace(&auth.SysBaseMenus); err != nil {
            return err
        }
        casbinInfos := systemReq.DefaultCasbin()
        authorityId := strconv.Itoa(int(auth.AuthorityId))
        rules := [][]string{}
        for _, v := range casbinInfos {
            rules = append(rules, []string{authorityId, v.Path, v.Method})
        }
        return CasbinServiceApp.AddPolicies(tx, rules)
    })

    return auth, e
}
```

### 2. 权限检查中间件

```go
// server/middleware/casbin_rbac.go
func CasbinHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        waitUse, _ := utils.GetClaims(c)
        //获取请求的PATH
        path := c.Request.URL.Path
        obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
        // 获取请求方法
        act := c.Request.Method
        // 获取用户的角色
        sub := strconv.Itoa(int(waitUse.AuthorityId))
        e := utils.GetCasbin() // 判断策略中是否存在
        success, _ := e.Enforce(sub, obj, act)
        if !success {
            response.FailWithDetailed(gin.H{}, "权限不足", c)
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 3. JWT认证中间件

```go
// server/middleware/jwt.go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := utils.GetToken(c)
        if token == "" {
            response.NoAuth("未登录或非法访问，请登录", c)
            c.Abort()
            return
        }
        if isBlacklist(token) {
            response.NoAuth("您的帐户异地登陆或令牌失效", c)
            utils.ClearToken(c)
            c.Abort()
            return
        }
        j := utils.NewJWT()
        claims, err := j.ParseToken(token)
        if err != nil {
            if errors.Is(err, utils.TokenExpired) {
                response.NoAuth("登录已过期，请重新登录", c)
                utils.ClearToken(c)
                c.Abort()
                return
            }
            response.NoAuth(err.Error(), c)
            utils.ClearToken(c)
            c.Abort()
            return
        }

        c.Set("claims", claims)
        // Token刷新逻辑...
        c.Next()
    }
}
```

### 4. 获取用户菜单树

```go
// server/service/system/sys_menu.go
func (menuService *MenuService) getMenuTreeMap(authorityId uint) (treeMap map[uint][]system.SysMenu, err error) {
    var allMenus []system.SysMenu
    var baseMenu []system.SysBaseMenu
    var btns []system.SysAuthorityBtn
    treeMap = make(map[uint][]system.SysMenu)

    var SysAuthorityMenus []system.SysAuthorityMenu
    err = global.GVA_DB.Where("sys_authority_authority_id = ?", authorityId).Find(&SysAuthorityMenus).Error
    if err != nil {
        return
    }

    var MenuIds []string
    for i := range SysAuthorityMenus {
        MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
    }

    err = global.GVA_DB.Where("id in (?)", MenuIds).Order("sort").Preload("Parameters").Find(&baseMenu).Error
    if err != nil {
        return
    }

    for i := range baseMenu {
        allMenus = append(allMenus, system.SysMenu{
            SysBaseMenu: baseMenu[i],
            AuthorityId: authorityId,
            MenuId:      baseMenu[i].ID,
            Parameters:  baseMenu[i].Parameters,
        })
    }

    err = global.GVA_DB.Where("authority_id = ?", authorityId).Preload("SysBaseMenuBtn").Find(&btns).Error
    if err != nil {
        return
    }
    var btnMap = make(map[uint]map[string]uint)
    for _, v := range btns {
        if btnMap[v.SysMenuID] == nil {
            btnMap[v.SysMenuID] = make(map[string]uint)
        }
        btnMap[v.SysMenuID][v.SysBaseMenuBtn.Name] = authorityId
    }
    for _, v := range allMenus {
        v.Btns = btnMap[v.SysBaseMenu.ID]
        treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
    }
    return treeMap, err
}
```

## 配置说明

### 1. 系统配置

```yaml
# server/config.yaml
system:
  use-strict-auth: true  # 启用严格权限控制
  router-prefix: /api    # API路由前缀
  addr: 8888            # 服务端口
```

### 2. JWT配置

```yaml
jwt:
  signing-key: e3ad29b7-c241-4bd3-a894-ef867165178d
  expires-time: 7d
  buffer-time: 1d
  issuer: qmPlus
```

### 3. 数据库配置

```yaml
mysql:
  prefix: ""
  port: "3306"
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: gva
  username: root
  password: toor
  path: 127.0.0.1
```

## 使用示例

### 1. 创建角色

```go
// 创建父角色
parentRole := system.SysAuthority{
    AuthorityId:   1001,
    AuthorityName: "部门管理员",
    ParentId:      utils.Pointer[uint](0),
    DefaultRouter: "dashboard",
}

// 创建子角色
childRole := system.SysAuthority{
    AuthorityId:   1002,
    AuthorityName: "部门员工",
    ParentId:      utils.Pointer[uint](1001),
    DefaultRouter: "dashboard",
}
```

### 2. 分配菜单权限

```go
// 为角色分配菜单权限
func assignMenuToRole(authorityId uint, menuIds []uint) error {
    var menus []system.SysBaseMenu
    err := global.GVA_DB.Where("id in (?)", menuIds).Find(&menus).Error
    if err != nil {
        return err
    }
    
    var auth system.SysAuthority
    auth.AuthorityId = authorityId
    auth.SysBaseMenus = menus
    
    return AuthorityServiceApp.SetMenuAuthority(&auth)
}
```

### 3. 权限检查

```go
// 检查用户是否有权限访问某个API
func checkPermission(userId uint, path string, method string) bool {
    var user system.SysUser
    err := global.GVA_DB.Preload("Authority").Where("id = ?", userId).First(&user).Error
    if err != nil {
        return false
    }
    
    e := utils.GetCasbin()
    authorityId := strconv.Itoa(int(user.AuthorityId))
    success, _ := e.Enforce(authorityId, path, method)
    return success
}
```

## 最佳实践

### 1. 角色设计原则

- **层次清晰**：设计清晰的角色层次结构
- **权限最小化**：遵循最小权限原则
- **职责分离**：不同角色承担不同职责

### 2. 权限分配策略

- **继承优先**：优先使用权限继承，减少重复配置
- **细粒度控制**：支持菜单、按钮、数据多级权限控制
- **动态调整**：支持运行时权限调整

### 3. 安全建议

- **启用严格模式**：设置 `use-strict-auth: true`
- **定期审计**：定期检查权限分配是否合理
- **最小权限**：用户只分配必要的权限
- **权限回收**：及时回收离职用户权限

### 4. 性能优化

- **缓存策略**：合理使用权限缓存
- **批量操作**：减少数据库查询次数
- **索引优化**：为权限相关字段建立索引

## 总结

Gin-Vue-Admin的权限管理系统通过角色层次结构、权限继承机制和严格的安全控制，为企业级应用提供了完善的权限管理解决方案。系统设计遵循"子角色权限是父角色权限的子集"的原则，确保了权限管理的安全性和可控性。

通过合理配置和使用，可以实现从简单的小型应用到复杂的企业级应用的各种权限管理需求。
