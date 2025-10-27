<template>
  <div>
    <div class="sticky top-0.5 z-10">
      <el-input v-model="filterText" class="w-3/5" placeholder="筛选" />
      <el-button class="float-right" type="primary" @click="relation"
        >确 定</el-button
      >
    </div>
    <!-- 菜单树内容区域 -->
    <div class="tree-content clear-both">
      <!-- Element Plus 滚动条组件 -->
      <el-scrollbar>
        <!-- Element Plus 树形组件：用于显示可选择的菜单树 -->
        <!-- 
          ref="menuTree" - 组件引用，用于获取 el-tree 实例调用方法（如 getCheckedNodes）
          :data="menuTreeData" - 绑定菜单树数据源（从第96行 menuTreeData 中取值）
          :default-checked-keys="menuTreeIds" - 默认选中的菜单ID数组（从第97行 menuTreeIds 中取值）
          :props="menuDefaultProps" - 字段映射配置（从第99-107行 menuDefaultProps 中取值）
          default-expand-all - 默认展开所有节点
          highlight-current - 高亮当前选中节点
          node-key="ID" - 节点的唯一标识字段为 "ID"
          show-checkbox - 显示复选框用于选择
          :filter-node-method="filterNode" - 节点过滤方法（从第216行 filterNode 函数中取值）
          @check="nodeChange" - 勾选状态改变时触发（从第139行 nodeChange 函数中取值）
        -->
        <el-tree
          ref="menuTree"
          :data="menuTreeData"
          :default-checked-keys="menuTreeIds"
          :props="menuDefaultProps"
          default-expand-all
          highlight-current
          node-key="ID"
          show-checkbox
          :filter-node-method="filterNode"
          @check="nodeChange"
        >
          <!-- 
            自定义节点插槽：用于定制每个节点的显示内容
            node - Element Plus 提供的节点对象，包含节点的状态信息
              node.label: 节点的显示文本（通过第154-156行的 label 函数映射生成）
              node.checked: 节点是否被选中（boolean）
              node.expanded: 节点是否展开（boolean）
              node.level: 节点的层级（数字）
              node.data: 节点的原始数据（等同于下面的 data）
            data - 当前节点的原始数据对象，来自 menuTreeData 数组
              data.ID: 菜单ID（来自后端数据库）
              data.name: 菜单名称（如 'dashboard', 'superAdmin'）
              data.path: 菜单路由路径
              data.component: 菜单对应的组件路径
              data.meta: 菜单元信息对象
                data.meta.title: 菜单标题（如 '仪表盘', '超级管理员'）
                data.meta.icon: 菜单图标
              data.menuBtn: 菜单对应的按钮权限数组
              data.children: 子菜单数组
            row - 从父组件传入的角色对象（从第85-92行的 props.row 中取值）
              row.authorityId: 角色ID（如 888）
              row.authorityName: 角色名称（如 '普通用户'）
              row.defaultRouter: 默认路由（如 'dashboard'）
              row.parentId: 父角色ID
          -->
          <template #default="{ node, data }">
            <span class="custom-tree-node">
              <!-- node.label - 显示菜单标题（由 label 函数映射，值为 data.meta.title） -->
              <span>{{ node.label }}</span>
              <!-- 
                node.checked - 菜单节点是否被选中
                !data.name?.startsWith('http://') - 排除外部链接菜单（以 http 开头的）
                条件：只有选中的菜单且不是外部链接，才显示"设为首页"按钮
              -->
              <span v-if="node.checked && !data.name?.startsWith('http://') && !data.name?.startsWith('https://')">
                <el-button
                  type="primary"
                  link
                  :style="{
                    color:
                      row.defaultRouter === data.name ? '#E6A23C' : '#85ce61'
                  }"
                  @click.stop="() => setDefault(data)"
                >
                  <!-- 
                    row.defaultRouter === data.name - 判断当前菜单是否是默认路由
                    如果是：显示"首页"；如果不是：显示"设为首页"
                  -->
                  {{ row.defaultRouter === data.name ? '首页' : '设为首页' }}
                </el-button>
              </span>
              <!-- 
                data.menuBtn.length - 菜单是否有关联的按钮权限
                条件：只有当菜单有按钮时，才显示"分配按钮"
              -->
              <span v-if="data.menuBtn.length">
                <el-button type="primary" link @click.stop="() => OpenBtn(data)">
                  分配按钮
                </el-button>
              </span>
            </span>
          </template>
        </el-tree>
      </el-scrollbar>
    </div>
    <el-dialog v-model="btnVisible" title="分配按钮" destroy-on-close>
      <el-table
        ref="btnTableRef"
        :data="btnData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="按钮名称" prop="name" />
        <el-table-column label="按钮备注" prop="desc" />
      </el-table>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import {
    getBaseMenuTree,
    getMenuAuthority,
    addMenuAuthority
  } from '@/api/menu'
  import { updateAuthority } from '@/api/authority'
  import { getAuthorityBtnApi, setAuthorityBtnApi } from '@/api/authorityBtn'
  import { nextTick, ref, watch } from 'vue'
  import { ElMessage } from 'element-plus'

  defineOptions({
    name: 'Menus'
  })

  const props = defineProps({
    row: {
      default: function () {
        return {}
      },
      type: Object
    }
  })

  const emit = defineEmits(['changeRow'])
  const filterText = ref('')
  const menuTreeData = ref([])
  const menuTreeIds = ref([])
  const needConfirm = ref(false)
  const menuDefaultProps = ref({
    children: 'children',
    label: function (data) { 
      return data.meta.title
    },
    disabled: function (data) {
      return props.row.defaultRouter === data.name
    }
  })

  const init = async () => {
    // 获取所有菜单树
    const res = await getBaseMenuTree()
    menuTreeData.value = res.data.menus
    const res1 = await getMenuAuthority({ authorityId: props.row.authorityId })
    const menus = res1.data.menus
    const arr = []
    menus.forEach((item) => {
      // 防止直接选中父级造成全选
      if (!menus.some((same) => same.parentId === item.menuId)) {
        arr.push(Number(item.menuId))
      }
    })
    menuTreeIds.value = arr
  }

  init()

  const setDefault = async (data) => {
    const res = await updateAuthority({
      authorityId: props.row.authorityId,
      AuthorityName: props.row.authorityName,
      parentId: props.row.parentId,
      defaultRouter: data.name
    })
    if (res.code === 0) {
      relation()
      emit('changeRow', 'defaultRouter', res.data.authority.defaultRouter)
    }
  }
  const nodeChange = () => {
    needConfirm.value = true
  }
  // 暴露给外层使用的切换拦截统一方法
  const enterAndNext = () => {
    relation()
  }
  // 关联树 确认方法
  const menuTree = ref(null)
  const relation = async () => {
    const checkArr = menuTree.value.getCheckedNodes(false, true)
    const res = await addMenuAuthority({
      menus: checkArr,
      authorityId: props.row.authorityId
    })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '菜单设置成功!'
      })
    }
  }

  defineExpose({ enterAndNext, needConfirm })

  const btnVisible = ref(false)

  const btnData = ref([])
  const multipleSelection = ref([])
  const btnTableRef = ref()
  let menuID = ''
  const OpenBtn = async (data) => {
    menuID = data.ID
    const res = await getAuthorityBtnApi({
      menuID: menuID,
      authorityId: props.row.authorityId
    })
    if (res.code === 0) {
      openDialog(data)
      await nextTick()
      if (res.data.selected) {
        res.data.selected.forEach((id) => {
          btnData.value.some((item) => {
            if (item.ID === id) {
              btnTableRef.value.toggleRowSelection(item, true)
            }
          })
        })
      }
    }
  }

  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  const openDialog = (data) => {
    btnVisible.value = true
    btnData.value = data.menuBtn
  }

  const closeDialog = () => {
    btnVisible.value = false
  }
  const enterDialog = async () => {
    const selected = multipleSelection.value.map((item) => item.ID)
    const res = await setAuthorityBtnApi({
      menuID,
      selected,
      authorityId: props.row.authorityId
    })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '设置成功' })
      btnVisible.value = false
    }
  }

  const filterNode = (value, data) => {
    if (!value) return true
    // console.log(data.mate.title)
    return data.meta.title.indexOf(value) !== -1
  }

  watch(filterText, (val) => {
    menuTree.value.filter(val)
  })
</script>

<style lang="scss" scoped>
  .custom-tree-node {
    span + span {
      @apply ml-3;
    }
  }
</style>
