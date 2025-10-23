所做修改

1.从sys_authority_menus表中删除了所有角色(sys_authority_authority_id)对应的sys_base_menu_id=7（官方网站菜单）的记录
即，delete from sys_authority_menus where sys_base_menu_id = 7







