<template>
  <div>
    <warning-bar title="注：GM管理 - 游戏用户管理" />
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="UserID">
          <el-input v-model="searchInfo.userId" placeholder="用户UserID" />
        </el-form-item>
        <el-form-item label="PlayerID">
          <el-input v-model="searchInfo.playerId" placeholder="玩家PlayerID" />
        </el-form-item>
        <el-form-item label="UniqueID">
          <el-input v-model="searchInfo.uniqueId" placeholder="玩家UniqueID" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchInfo.nickname" placeholder="玩家昵称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" @click="onReset"> 重置 </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addUser"
          >新增游戏用户</el-button
        >
      </div>
      <el-table :data="tableData" row-key="user_id">
        <el-table-column align="left" label="UserID" min-width="180" prop="user_id" />
        <el-table-column
          align="left"
          label="PlayerID"
          min-width="120"
          prop="player_id"
        />
        <el-table-column
          align="left"
          label="UniqueID"
          min-width="120"
          prop="unique_id"
        />
        <el-table-column
          align="left"
          label="昵称"
          min-width="150"
          prop="nickname"
        />
        <el-table-column
          align="left"
          label="等级"
          min-width="80"
          prop="lv"
        />
        <el-table-column
          align="left"
          label="战力"
          min-width="180"
          prop="power"
        />
        <el-table-column
          align="left"
          label="区服"
          min-width="80"
          prop="area_id"
        />
        <el-table-column align="left" label="注册时间" min-width="180" prop="register_time_formatted" />
        <el-table-column align="left" label="最后登录时间" min-width="180" prop="login_time_formatted" />
        <el-table-column label="操作" :min-width="appStore.operateMinWith" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteUserFunc(scope.row)"
              >删除</el-button
            >
            <el-button
              type="primary"
              link
              icon="edit"
              @click="openEdit(scope.row)"
              >编辑</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    
    <el-drawer
      v-model="addUserDialog"
      :size="appStore.drawerSize"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">用户</span>
          <div>
            <el-button @click="closeAddUserDialog">取 消</el-button>
            <el-button type="primary" @click="enterAddUserDialog"
              >确 定</el-button
            >
          </div>
        </div>
      </template>

      <el-form
        ref="userForm"
        :rules="rules"
        :model="userInfo"
        label-width="80px"
      >
        <el-form-item
          v-if="dialogFlag === 'add'"
          label="用户名"
          prop="userName"
        >
          <el-input v-model="userInfo.userName" />
        </el-form-item>
        <el-form-item v-if="dialogFlag === 'add'" label="密码" prop="password">
          <el-input v-model="userInfo.password" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickName">
          <el-input v-model="userInfo.nickName" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="userInfo.phone" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userInfo.email" />
        </el-form-item>
        <el-form-item label="用户角色" prop="authorityId">
          <el-cascader
            v-model="userInfo.authorityIds"
            style="width: 100%"
            :options="authOptions"
            :show-all-levels="false"
            :props="{
              multiple: true,
              checkStrictly: true,
              label: 'authorityName',
              value: 'authorityId',
              disabled: 'disabled',
              emitPath: false
            }"
            :clearable="false"
          />
        </el-form-item>
        <el-form-item label="启用" prop="disabled">
          <el-switch
            v-model="userInfo.enable"
            inline-prompt
            :active-value="1"
            :inactive-value="2"
          />
        </el-form-item>
        <el-form-item label="头像" label-width="80px">
          <SelectImage v-model="userInfo.headerImg" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import { useGMUserStore } from '@/pinia/gm/user'
  import { getAuthorityList } from '@/api/authority'
  import WarningBar from '@/components/warningBar/warningBar.vue'

  import { ref, watch, onMounted } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import SelectImage from '@/components/selectImage/selectImage.vue'
  import { useAppStore } from "@/pinia";

  defineOptions({
    name: 'GMUser'
  })

  const appStore = useAppStore()
  const gmUserStore = useGMUserStore()

  // 使用store中的状态
  const { 
    userList: tableData, 
    total, 
    page, 
    pageSize, 
    searchInfo,
    fetchUserList,
    createUser,
    updateUser,
    removeUser,
    resetSearchInfo,
    setPage,
    setPageSize
  } = gmUserStore

  // 调试信息
  watch(
    () => tableData.value,
    (newData) => {
      console.log('🔍 Table data changed:', newData)
      console.log('🔍 Table data length:', newData?.length)
      if (newData && newData.length > 0) {
        console.log('🔍 First table item:', newData[0])
        console.log('🔍 First item keys:', Object.keys(newData[0]))
      }
    },
    { immediate: true, deep: true }
  )

  const onSubmit = () => {
    setPage(1)
    fetchUserList()
  }

  const onReset = () => {
    resetSearchInfo()
    setPage(1)
    fetchUserList()
  }
  // 初始化相关
  const setAuthorityOptions = (AuthorityData, optionsData) => {
    AuthorityData &&
      AuthorityData.forEach((item) => {
        if (item.children && item.children.length) {
          const option = {
            authorityId: item.authorityId,
            authorityName: item.authorityName,
            children: []
          }
          setAuthorityOptions(item.children, option.children)
          optionsData.push(option)
        } else {
          const option = {
            authorityId: item.authorityId,
            authorityName: item.authorityName
          }
          optionsData.push(option)
        }
      })
  }

  // 分页
  const handleSizeChange = (val) => {
    setPageSize(val)
    fetchUserList()
  }

  const handleCurrentChange = (val) => {
    setPage(val)
    fetchUserList()
  }

  watch(
    () => tableData.value,
    () => {
      setAuthorityIds()
    }
  )

  const initPage = async () => {
    await fetchUserList()
    const res = await getAuthorityList()
    setOptions(res.data)
  }

  onMounted(() => {
    initPage()
  })

  const setAuthorityIds = () => {
    tableData.value &&
      tableData.value.forEach((user) => {
        user.authorityIds =
          user.authorities &&
          user.authorities.map((i) => {
            return i.authorityId
          })
      })
  }

  const authOptions = ref([])
  const setOptions = (authData) => {
    authOptions.value = []
    setAuthorityOptions(authData, authOptions.value)
  }

  const deleteUserFunc = async (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await removeUser(row.user_id)
        ElMessage.success('删除成功')
      } catch (error) {
        ElMessage.error(error.message || '删除失败')
      }
    })
  }

  // 弹窗相关
  const userInfo = ref({
    userName: '',
    password: '',
    nickName: '',
    headerImg: '',
    authorityId: '',
    authorityIds: [],
    enable: 1
  })

  const rules = ref({
    userName: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 5, message: '最低5位字符', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入用户密码', trigger: 'blur' },
      { min: 6, message: '最低6位字符', trigger: 'blur' }
    ],
    nickName: [{ required: true, message: '请输入用户昵称', trigger: 'blur' }],
    phone: [
      {
        pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/,
        message: '请输入合法手机号',
        trigger: 'blur'
      }
    ],
    email: [
      {
        pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g,
        message: '请输入正确的邮箱',
        trigger: 'blur'
      }
    ],
    authorityId: [
      { required: true, message: '请选择用户角色', trigger: 'blur' }
    ]
  })
  const userForm = ref(null)
  const enterAddUserDialog = async () => {
    userInfo.value.authorityId = userInfo.value.authorityIds[0]
    userForm.value.validate(async (valid) => {
      if (valid) {
        try {
          if (dialogFlag.value === 'add') {
            await createUser(userInfo.value)
            ElMessage({ type: 'success', message: '创建成功' })
          }
          if (dialogFlag.value === 'edit') {
            await updateUser(userInfo.value)
            ElMessage({ type: 'success', message: '编辑成功' })
          }
          closeAddUserDialog()
        } catch (error) {
          ElMessage.error(error.message || '操作失败')
        }
      }
    })
  }

  const addUserDialog = ref(false)
  const closeAddUserDialog = () => {
    userForm.value.resetFields()
    userInfo.value.headerImg = ''
    userInfo.value.authorityIds = []
    addUserDialog.value = false
  }

  const dialogFlag = ref('add')

  const addUser = () => {
    dialogFlag.value = 'add'
    addUserDialog.value = true
  }


  const openEdit = (row) => {
    dialogFlag.value = 'edit'
    userInfo.value = JSON.parse(JSON.stringify(row))
    addUserDialog.value = true
  }

</script>

<style lang="scss">
  .header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
  }
</style>
