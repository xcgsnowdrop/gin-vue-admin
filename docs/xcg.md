所做修改

1.从sys_authority_menus表中删除了所有角色(sys_authority_authority_id)对应的sys_base_menu_id=7（官方网站菜单）的记录
即，delete from sys_authority_menus where sys_base_menu_id = 7


# 数据库初始化逻辑
1. 每个表的数据的初始化在server/soruce/system/目录下各个文件的InitializeData方法中
例如，
api表，在server/source/system/casbin.go的InitializeData方法
菜单表，在server/source/system/menu.go的InitializeData方法中



# 菜单说明
1.除最高权限用户外可以不用的菜单
官方网站

# 权限说明
最高权限--系统管理员，仅1个，可以管理所有

