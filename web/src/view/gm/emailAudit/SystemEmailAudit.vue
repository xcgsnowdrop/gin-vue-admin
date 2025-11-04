<template>
    <div>
      <warning-bar title="注：GM管理 - 系统邮件列表" />
      <div class="gva-search-box">
        <el-form ref="searchForm" :inline="true" :model="searchInfo">
          <el-form-item label="PlayerId">
            <el-input v-model="searchInfo.player_id" placeholder="PlayerId" />
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
          <el-button type="primary" icon="plus" @click="openDialog">
            新增
          </el-button>
        </div>
        <el-table :data="tableData" row-key="id" v-loading="loading">
          <el-table-column align="left" label="ID" min-width="80" prop="id" />
          <el-table-column align="left" label="申请人" min-width="80">
            <template #default="scope">
              {{ scope.row.applicant?.nickName || scope.row.applicantId || '-' }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="申请时间" min-width="180" prop="applicantTime">
            <template #default="scope">
              {{ scope.row.applicantTime ? formatTimestamp(new Date(scope.row.applicantTime)) : '-' }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="审核人" min-width="80">
            <template #default="scope">
              {{ scope.row.auditor?.nickName || scope.row.auditorId || '-' }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="审核说明" min-width="180" prop="auditComment">
            <template #default="scope">
              {{ scope.row.auditComment || '-' }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="状态" min-width="100" prop="status">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="邮件类型"
            min-width="100"
            prop="emailType"
          />
          <el-table-column
            align="left"
            label="发送区服"
            min-width="150"
            prop="areaIds"
          />
          <el-table-column align="left" label="开始生效时间" min-width="180">
            <template #default="scope">
              {{ scope.row.startTime ? formatTimestamp(scope.row.startTime) : '-' }}
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="附件"
            min-width="300"
          >
            <template #default="scope">
              <div v-if="scope.row.emailAttachments && scope.row.emailAttachments.length > 0" class="attachments-list">
                <el-tag
                  v-for="(attachment, index) in scope.row.emailAttachments"
                  :key="index"
                  size="small"
                  type="success"
                  class="attachment-tag"
                >
                  {{ formatAttachment(attachment) }}
                </el-tag>
              </div>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="标题"
            min-width="200"
          >
            <template #default="scope">
              <MultilingualCell :value="scope.row.emailTitle" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="内容"
            min-width="250"
          >
            <template #default="scope">
              <MultilingualCell :value="scope.row.emailContent" :max-length="30" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="备注"
            min-width="150"
            prop="emailRemark"
          />
          
          <el-table-column align="left" label="最大注册时间" min-width="180">
            <template #default="scope">
              {{ scope.row.maxRegTime ? formatTimestamp(scope.row.maxRegTime) : '-' }}
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="操作"
            fixed="right"
            min-width="350"
          >
            <template #default="scope">
              <el-button
                v-if="canEdit(scope.row)"
                type="primary"
                link
                icon="edit"
                class="table-button"
                @click="updateRow(scope.row)"
              >
                编辑
              </el-button>
              <el-button
                v-if="canWithdraw(scope.row)"
                type="warning"
                link
                icon="delete"
                @click="withdrawRow(scope.row)"
              >
                撤回
              </el-button>
              <el-button
                v-if="canReview(scope.row)"
                type="success"
                link
                icon="check"
                @click="openReviewDialog(scope.row)"
              >
                审核
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="gva-pagination">
          <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
          />
        </div>
      </div>
  
      <!-- 新增/修改弹窗 -->
      <el-drawer
        v-model="dialogFormVisible"
        destroy-on-close
        size="800"
        :show-close="false"
        :before-close="closeDialog"
      >
        <template #header>
          <div class="flex justify-between items-center">
            <span class="text-lg">{{ type === 'create' ? '添加' : '修改' }}</span>
            <div>
              <el-button type="primary" @click="enterDialog"> 确 定 </el-button>
              <el-button @click="closeDialog"> 取 消 </el-button>
            </div>
          </div>
        </template>
  
        <el-form
          ref="elFormRef"
          :model="formData"
          label-position="top"
          :rules="rule"
          label-width="80px"
        >
          <el-form-item label="邮件类型:" prop="emailType" required>
            <el-select
              v-model="formData.emailType"
              placeholder="请选择邮件类型"
              style="width: 100%"
              :clearable="true"
            >
              <el-option label="类型1" :value="1" />
              <el-option label="类型2" :value="2" />
              <el-option label="类型3" :value="3" />
            </el-select>
          </el-form-item>
          <!-- <MultilingualInput
            label="邮件发件人"
            prop="senderI18n"
            :required="true"
            v-model="formData.senderI18n"
            v-model:active-tab="activeSenderTab"
          /> -->
          <MultilingualInput
            label="邮件标题"
            prop="emailTitle"
            :required="true"
            v-model="formData.emailTitle"
            v-model:active-tab="activeTitleTab"
          />
          <MultilingualRichEdit
            label="邮件内容"
            prop="emailContent"
            :required="true"
            v-model="formData.emailContent"
            v-model:active-tab="activeContentTab"
          />
          <el-form-item label="邮件附件:" prop="emailAttachments">
            <div class="attachments-form">
              <div
                v-for="(attachment, index) in formData.emailAttachments"
                :key="index"
                class="attachment-item"
              >
                <el-select
                  v-model="attachment.type"
                  placeholder="选择资源类型"
                  style="width: 30%"
                  clearable
                  @change="handleAttachmentTypeChange(index, $event)"
                >
                  <el-option
                    v-for="resType in resourceTypes"
                    :key="resType.type"
                    :label="resType.name"
                    :value="resType.type"
                  />
                </el-select>
                <el-select
                  v-model="attachment.id"
                  placeholder="选择资源"
                  style="width: 30%; margin-left: 10px"
                  clearable
                  :disabled="!attachment.type"
                >
                  <el-option
                    v-for="resource in getResourceListByType(attachment.type)"
                    :key="resource.id"
                    :label="resource.name"
                    :value="resource.id"
                  />
                </el-select>
                <el-input-number
                  v-model="attachment.num"
                  :min="1"
                  placeholder="数量"
                  style="width: 25%; margin-left: 10px"
                  controls-position="right"
                />
                <el-button
                  type="danger"
                  icon="delete"
                  circle
                  style="margin-left: 10px"
                  @click="removeAttachment(index)"
                />
              </div>
              <el-button
                type="primary"
                icon="plus"
                style="width: 100%; margin-top: 10px"
                @click="addAttachment"
              >
                添加附件
              </el-button>
            </div>
          </el-form-item>
          <el-form-item label="邮件备注:" prop="emailRemark">
            <el-input v-model="formData.emailRemark" placeholder="请输入邮件备注" />
          </el-form-item>
          <el-form-item label="开始生效时间:" prop="startTime" required>
            <el-date-picker
              v-model="formData.startTime"
              type="datetime"
              placeholder="请选择开始时间"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="生效区服列表:" prop="areaIds">
            <el-input v-model="formData.areaIds" placeholder="请输入生效区服列表,用英文逗号分隔，为空表示全区服" />
          </el-form-item>
          <el-form-item label="最大注册时间:" prop="maxRegTime" required>
            <el-date-picker
              v-model="formData.maxRegTime"
              type="datetime"
              placeholder="请选择最大注册时间"
              style="width: 100%"
            />
          </el-form-item>
        </el-form>
      </el-drawer>

      <!-- 审核弹窗 -->
      <el-dialog
        v-model="reviewDialogVisible"
        title="审核邮件申请"
        width="600px"
        :before-close="closeReviewDialog"
      >
        <el-form
          ref="reviewFormRef"
          :model="reviewForm"
          label-position="top"
          :rules="reviewRules"
          label-width="80px"
        >
          <el-form-item label="审核结果" prop="status" required>
            <el-radio-group v-model="reviewForm.status">
              <el-radio :label="2">通过</el-radio>
              <el-radio :label="3">拒绝</el-radio>
              <el-radio :label="4">待修改</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="审核说明" prop="auditComment">
            <el-input
              v-model="reviewForm.auditComment"
              type="textarea"
              :rows="4"
              placeholder="请输入审核说明（可选）"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="closeReviewDialog">取 消</el-button>
            <el-button type="primary" @click="submitReview">确 定</el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script setup>
  import { onMounted, watch, ref, reactive } from 'vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import MultilingualCell from '@/components/multilingual/MultilingualCell.vue'
  import MultilingualInput from '@/components/multilingual/MultilingualInput.vue'
  import MultilingualRichEdit from '@/components/multilingual/MultilingualRichEdit.vue'
  import { useGMSystemEmailAuditStore } from '@/pinia/gm/systemEmailAudit'
  import { useUserStore } from '@/pinia/modules/user'
  import { storeToRefs } from 'pinia'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { languageOptions, initMultilingualData, useMultilingual } from '@/composables/useMultilingual'
  import { timestampToDate, formatTimestamp } from '@/utils/timestamp'

  defineOptions({
    name: 'GmSystemEmailAudit'
  })
  
  const gmSystemEmailAuditStore = useGMSystemEmailAuditStore()
  const userStore = useUserStore()
  
  // 使用store中的状态
  const { 
    systemEmailAuditList: tableData, 
    loading, 
    total, 
    page, 
    pageSize, 
    searchInfo,
    resourceTypes,
    resourceList,
    resourceMap
  } = storeToRefs(gmSystemEmailAuditStore)
  
  const { userInfo } = storeToRefs(userStore)
  
  const {
    fetchSystemEmailAuditList,
    sendSystemEmailAudit,
    fetchResourceTypes,
    fetchResourceList,
    getResourceTypeName,
    getResourceName,
    resetSearchInfo,
    setPage,
    setPageSize,
    deleteSystemEmailAudit,
    updateSystemEmailAudit,
    reviewSystemEmailAudit,
  } = gmSystemEmailAuditStore

  // 获取当前用户ID
  const getCurrentUserId = () => {
    return userInfo.value?.ID || userInfo.value?.id || 0
  }

  // 状态格式化
  const getStatusText = (status) => {
    const statusMap = {
      1: '待审核',
      2: '通过',
      3: '拒绝',
      4: '待修改'
    }
    return statusMap[status] || '未知'
  }

  const getStatusType = (status) => {
    const typeMap = {
      1: 'warning',
      2: 'success',
      3: 'danger',
      4: 'info'
    }
    return typeMap[status] || ''
  }

  // 权限判断：是否可以编辑
  const canEdit = (row) => {
    const currentUserId = getCurrentUserId()
    // 只能编辑自己的申请，且状态为待审核或待修改
    return row.applicantId === currentUserId && (row.status === 1 || row.status === 4)
  }

  // 权限判断：是否可以撤回
  const canWithdraw = (row) => {
    const currentUserId = getCurrentUserId()
    // 只能撤回自己的申请，且状态为待审核
    return row.applicantId === currentUserId && row.status === 1
  }

  // 权限判断：是否可以审核
  // 注意：这里的权限判断应该调用后端API来确认，前端只是简单判断状态
  // 实际权限验证在后端完成
  const canReview = (row) => {
    const currentUserId = getCurrentUserId()
    // 不能审核自己的申请
    if (row.applicantId === currentUserId) {
      return false
    }
    // 只能审核待审核状态的申请
    return row.status === 1
  }
  
  // 格式化附件显示
  const formatAttachment = (attachment) => {
    if (!attachment) return ''
    const typeName = getResourceTypeName(attachment.type)
    const resourceName = getResourceName(attachment.type, attachment.id)
    return `${typeName} - ${resourceName} × ${attachment.num}`
  }

  
  // 查询数据
  const onSubmit = () => {
    setPage(1)
    fetchSystemEmailAuditList()
  }
  
  // 重置搜索
  const onReset = () => {
    resetSearchInfo()
    setPage(1)
    fetchSystemEmailAuditList()
  }
  
  // 分页处理
  const handleCurrentChange = (val) => {
    setPage(val)
    fetchSystemEmailAuditList()
  }
  
  const handleSizeChange = (val) => {
    setPageSize(val)
    fetchSystemEmailAuditList()
  }

  
    const elFormRef = ref()

    // 行为控制标记（弹窗内部需要增还是改）
    const type = ref('')

    // 弹窗控制标记
    const dialogFormVisible = ref(false)

    // 使用多语言 Composable
    const { activeTitleTab, activeContentTab, resetActiveTabs, setActiveTabsFromData } = useMultilingual()

    // 附件资源列表缓存（每个类型对应一个资源列表）
    const attachmentResourceLists = ref({})

    // 验证规则
    const rule = reactive({
      emailType: [
        { required: true, message: '请选择邮件类型', trigger: 'change' }
      ],
      emailTitle: [
        { 
          validator: (rule, value, callback) => {
            if (!value || typeof value !== 'object') {
              callback(new Error('请输入邮件标题'))
              return
            }
            // 检查是否至少有一个语言版本有值
            const hasValue = Object.values(value).some(v => v && v.trim())
            if (!hasValue) {
              callback(new Error('请输入至少一种语言的标题'))
              return
            }
            callback()
          },
          trigger: 'blur'
        }
      ],
      emailContent: [
        { 
          validator: (rule, value, callback) => {
            if (!value || typeof value !== 'object') {
              callback(new Error('请输入邮件内容'))
              return
            }
            // 检查是否至少有一个语言版本有值
            const hasValue = Object.values(value).some(v => v && v.trim())
            if (!hasValue) {
              callback(new Error('请输入至少一种语言的内容'))
              return
            }
            callback()
          },
          trigger: 'blur'
        }
      ],
      startTime: [
        { required: true, message: '请选择开始生效时间', trigger: 'change' }
      ],
    })

    // 自动化生成的字典（可能为空）以及字段
    const formData = ref({
      id: null, // 更新时需要
      emailType: 1,  // 默认选中类型1
      emailTitle: initMultilingualData(),
      emailContent: initMultilingualData(),
      emailAttachments: [],
      emailRemark: '',
      startTime: null,
      areaIds: '', // 逗号分隔的区服ID列表，为空表示全区服
      maxRegTime: null, // 最大注册时间
    })

    // 更新行
    const updateRow = async (row) => {
        type.value = 'update'
        // 确保多语言数据格式正确
        const data = { ...row }
        // 确保ID字段存在
        if (!data.id && row.id) {
          data.id = row.id
        }
      
        // 确保所有语言字段都存在
        languageOptions.forEach(lang => {
          if (!data.emailTitle[lang.code]) {
            data.emailTitle[lang.code] = ''
          }
          if (!data.emailContent[lang.code]) {
            data.emailContent[lang.code] = ''
          }
        })

        // 转换时间戳为 Date 对象（用于日期选择器）
        if (data.startTime) {
          data.startTime = timestampToDate(data.startTime)
        }
        if (data.maxRegTime) {
          data.maxRegTime = timestampToDate(data.maxRegTime)
        }

        // 处理区服ID列表
        if (Array.isArray(data.areaIds)) {
          data.areaIds = data.areaIds.join(',')
        } else if (!data.areaIds) {
          data.areaIds = ''
        }

        formData.value = data
        
        // 设置默认活动标签页为第一个有内容的语言
        setActiveTabsFromData(data)
        
        dialogFormVisible.value = true
    }

    // 撤回行
    const withdrawRow = (row) => {
      ElMessageBox.confirm('确定要撤回此申请吗？撤回后无法恢复。', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        withdrawRowFunc(row)
      }).catch(() => {
        // 用户取消
      })
    }

    // 撤回行
    const withdrawRowFunc = async (row) => {
      try {
        await deleteSystemEmailAudit(row.id)
        ElMessage({
          type: 'success',
          message: '撤回成功'
        })
        // 如果撤回后当前页没有数据且不是第一页，则回到上一页
        if (tableData.value.length === 1 && page.value > 1) {
          setPage(page.value - 1)
        }
      } catch (error) {
        ElMessage({
          type: 'error',
          message: error.message || '撤回失败'
        })
      }
    }

    // 审核相关
    const reviewDialogVisible = ref(false)
    const reviewFormRef = ref()
    const reviewForm = ref({
      id: null,
      status: 2, // 默认通过
      auditComment: ''
    })
    const currentReviewRow = ref(null)

    const reviewRules = reactive({
      status: [
        { required: true, message: '请选择审核结果', trigger: 'change' }
      ]
    })

    // 打开审核对话框
    const openReviewDialog = (row) => {
      currentReviewRow.value = row
      reviewForm.value = {
        id: row.id,
        status: 2, // 默认通过
        auditComment: ''
      }
      reviewDialogVisible.value = true
    }

    // 关闭审核对话框
    const closeReviewDialog = () => {
      reviewDialogVisible.value = false
      reviewFormRef.value?.resetFields()
      reviewForm.value = {
        id: null,
        status: 2,
        auditComment: ''
      }
      currentReviewRow.value = null
    }

    // 提交审核
    const submitReview = async () => {
      reviewFormRef.value?.validate(async (valid) => {
        if (!valid) return
        
        try {
          await reviewSystemEmailAudit(reviewForm.value)
          ElMessage({
            type: 'success',
            message: '审核成功'
          })
          closeReviewDialog()
        } catch (error) {
          ElMessage({
            type: 'error',
            message: error.message || '审核失败'
          })
        }
      })
    }

  // 打开弹窗
  const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
  }

  
  // 添加附件
  const addAttachment = () => {
    formData.value.emailAttachments.push({
      id: null,
      type: null,
      num: 1
    })
  }

  // 删除附件
  const removeAttachment = (index) => {
    formData.value.emailAttachments.splice(index, 1)
  }

  // 处理附件资源类型变化
  const handleAttachmentTypeChange = async (index, type) => {
    const attachment = formData.value.emailAttachments[index]
    // 清空资源ID选择
    attachment.id = null
    
    if (type) {
      // 如果该类型的资源列表还未加载，则加载
      if (!attachmentResourceLists.value[type]) {
        try {
          await fetchResourceList(type)
          // 从 resourceMap 中获取资源列表
          attachmentResourceLists.value[type] = resourceList.value
        } catch (error) {
          console.error(`加载资源类型 ${type} 失败:`, error)
          attachmentResourceLists.value[type] = []
        }
      }
    }
  }

  // 根据资源类型获取资源列表
  const getResourceListByType = (type) => {
    if (!type) return []
    // 优先从缓存中获取
    if (attachmentResourceLists.value[type]) {
      return attachmentResourceLists.value[type]
    }
    // 如果缓存中没有，从 store 的 resourceMap 中构建
    if (resourceMap.value[type]) {
      return Object.keys(resourceMap.value[type]).map(id => ({
        id: parseInt(id),
        name: resourceMap.value[type][id]
      }))
    }
    return []
  }

  // 关闭弹窗
  const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
      id: null,
      emailType: 1,  // 默认选中类型1
      emailTitle: initMultilingualData(),
      emailContent: initMultilingualData(),
      emailAttachments: [],
      emailRemark: '',
      startTime: null,
      areaIds: '', // 区服ID列表，为空表示全区服
      maxRegTime: null, // 最大注册时间
    }
    attachmentResourceLists.value = {}
    resetActiveTabs()
  }
  // 弹窗确定
  const enterDialog = async () => {
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      
      try {
        // 过滤掉无效的附件（只保留 type、id、num 都有效的附件）
        const validAttachments = formData.value.emailAttachments.filter(
          att => att.type !== null && att.id !== null && att.num > 0
        )
        
        // 构建提交数据
        const submitData = {
          ...formData.value,
          emailAttachments: validAttachments
        }
        
        switch (type.value) {
          case 'create':
            await sendSystemEmailAudit(submitData)
            break
          case 'update':
            await updateSystemEmailAudit(submitData)
            break
          default:
            await sendSystemEmailAudit(submitData)
            break
        }
        
        ElMessage({
          type: 'success',
          message: '创建/更改成功'
        })
        closeDialog()
      } catch (error) {
        ElMessage({
          type: 'error',
          message: error.message || '操作失败'
        })
      }
    })
  }
  
  watch(
      () => tableData.value,
      (newValue, oldValue) => {
        console.log('tableData 变化了')
        console.log('新值:', newValue)
        console.log('旧值:', oldValue)
        console.log('新值长度:', newValue?.length)
      },
      { deep: true, immediate: true }
  )
  
  
  // 初始化
  onMounted(async () => {
    try {
      // 并行获取资源类型和道具列表
      await Promise.all([
        fetchResourceTypes(),
        fetchSystemEmailAuditList()
      ])
    } catch (error) {
      console.error('初始化失败:', error)
    }
  })
  </script>
  
  <style scoped>
  .text-green-600 {
    color: #16a34a;
  }
  
  .text-red-600 {
    color: #dc2626;
  }

  /* 附件列表样式 */
  .attachments-list {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .attachment-tag {
    margin: 0;
  }

  /* 模板参数列表样式 */
  .template-params-list {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .param-tag {
    margin: 0;
  }

  /* 附件表单样式 */
  .attachments-form {
    width: 100%;
  }

  .attachment-item {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
    padding: 10px;
    background-color: #f5f7fa;
    border-radius: 4px;
  }

  .attachment-item:last-of-type {
    margin-bottom: 0;
  }
  </style>
  