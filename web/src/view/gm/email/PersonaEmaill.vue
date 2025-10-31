<template>
    <div>
      <warning-bar title="注：GM管理 - 个人邮件列表" />
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
        <el-table :data="tableData" row-key="email_id" v-loading="loading">
          <el-table-column align="left" label="ID" min-width="80" prop="email_id" />
          <el-table-column
            align="left"
            label="PlayerId"
            min-width="120"
            prop="player_id"
          />
          <el-table-column
            align="left"
            label="邮件类型"
            min-width="100"
            prop="type"
          />
          <el-table-column
            align="left"
            label="模板ID"
            min-width="100"
            prop="tpl_id"
          />
          <el-table-column
            align="left"
            label="模板参数"
            min-width="150"
          >
            <template #default="scope">
              <div v-if="scope.row.tpl_params && getTemplateParams(scope.row.tpl_params).length > 0" class="template-params-list">
                <el-tag
                  v-for="(param, index) in getTemplateParams(scope.row.tpl_params)"
                  :key="index"
                  size="small"
                  type="info"
                  class="param-tag"
                >
                  {{ param }}
                </el-tag>
              </div>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="创建时间"
            min-width="180"
            prop="create_time_formatted"
          />
          <el-table-column
            align="left"
            label="开始生效时间"
            min-width="180"
            prop="start_time_formatted"
          />
          <el-table-column
            align="left"
            label="结束生效时间"
            min-width="180"
            prop="end_time_formatted"
          />
          <el-table-column
            align="left"
            label="附件"
            min-width="300"
          >
            <template #default="scope">
              <div v-if="scope.row.attachments && scope.row.attachments.length > 0" class="attachments-list">
                <el-tag
                  v-for="(attachment, index) in scope.row.attachments"
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
            label="发件人"
            min-width="150"
          >
            <template #default="scope">
              <MultilingualCell :value="scope.row.senderI18n" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="标题"
            min-width="200"
          >
            <template #default="scope">
              <MultilingualCell :value="scope.row.titleI18n" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="内容"
            min-width="250"
          >
            <template #default="scope">
              <MultilingualCell :value="scope.row.contentI18n" :max-length="30" />
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="备注"
            min-width="150"
            prop="remark"
          />
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
          <el-form-item label="玩家PlayerID:" prop="player_id" required>
            <el-input v-model="formData.player_id" placeholder="请输入PlayerId" />
          </el-form-item>
          <el-form-item label="邮件发件人:" prop="senderI18n" required>
            <el-tabs v-model="activeTitleTab" type="border-card" class="multilingual-tabs">
              <el-tab-pane
                v-for="lang in languageOptions"
                :key="lang.code"
                :label="lang.label"
                :name="lang.code"
              >
                <el-input
                  v-model="formData.senderI18n[lang.code]"
                  :clearable="true"
                  :placeholder="`请输入${lang.label}发件人`"
                />
              </el-tab-pane>
            </el-tabs>
          </el-form-item>
          <el-form-item label="邮件标题:" prop="titleI18n" required>
            <el-tabs v-model="activeTitleTab" type="border-card" class="multilingual-tabs">
              <el-tab-pane
                v-for="lang in languageOptions"
                :key="lang.code"
                :label="lang.label"
                :name="lang.code"
              >
                <el-input
                  v-model="formData.titleI18n[lang.code]"
                  :clearable="true"
                  :placeholder="`请输入${lang.label}标题`"
                />
              </el-tab-pane>
            </el-tabs>
          </el-form-item>
          <el-form-item label="邮件内容:" prop="contentI18n" required>
            <el-tabs v-model="activeContentTab" type="border-card" class="multilingual-tabs">
              <el-tab-pane
                v-for="lang in languageOptions"
                :key="lang.code"
                :label="lang.label"
                :name="lang.code"
              >
                <RichEdit v-model="formData.contentI18n[lang.code]" />
              </el-tab-pane>
            </el-tabs>
          </el-form-item>
          <el-form-item label="邮件附件:" prop="attachments">
            <div class="attachments-form">
              <div
                v-for="(attachment, index) in formData.attachments"
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
          <el-form-item label="邮件备注:" prop="remark">
            <el-input v-model="formData.remark" placeholder="请输入邮件备注" />
          </el-form-item>
          <el-form-item label="开始生效时间:" prop="startTime" required>
            <el-date-picker
              v-model="formData.startTime"
              type="datetime"
              placeholder="请选择开始时间"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="邮件类型:" prop="type" required>
            <el-select
              v-model="formData.type"
              placeholder="请选择邮件类型"
              style="width: 100%"
              :clearable="true"
            >
              <el-option label="类型1" value="1" />
              <el-option label="类型2" value="2" />
              <el-option label="类型3" value="3" />
            </el-select>
          </el-form-item>
        </el-form>
      </el-drawer>
    </div>
  </template>
  
  <script setup>
  import { onMounted, watch, ref, reactive } from 'vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import RichEdit from '@/components/richtext/rich-edit.vue'
  import MultilingualCell from '@/components/multilingual/MultilingualCell.vue'
  import { useGMPersonalEmailStore } from '@/pinia/gm/personalEamil'
  import { storeToRefs } from 'pinia'
  import { ElMessage } from 'element-plus'
  import { languageOptions, initMultilingualData, initSenderI18nDefault, useMultilingual } from '@/composables/useMultilingual'
  
  defineOptions({
    name: 'GmPersonalEmail'
  })
  
  const gmPersonalEmailStore = useGMPersonalEmailStore()
  
  // 使用store中的状态
  const { 
    personalEmailList: tableData, 
    loading, 
    total, 
    page, 
    pageSize, 
    searchInfo,
    resourceTypes,
    resourceList,
    resourceMap
  } = storeToRefs(gmPersonalEmailStore)
  
  const {
    fetchPersonalEmailList,
    sendPersonalEmail,
    fetchResourceTypes,
    fetchResourceList,
    getResourceTypeName,
    getResourceName,
    resetSearchInfo,
    setPage,
    setPageSize
  } = gmPersonalEmailStore
  
  // 格式化附件显示
  const formatAttachment = (attachment) => {
    if (!attachment) return ''
    const typeName = getResourceTypeName(attachment.type)
    const resourceName = getResourceName(attachment.type, attachment.id)
    return `${typeName} - ${resourceName} × ${attachment.num}`
  }

  // 获取模板参数列表
  const getTemplateParams = (tplParams) => {
    if (!tplParams) return []
    
    // 如果已经是数组，直接返回
    if (Array.isArray(tplParams)) {
      return tplParams.filter(param => param !== null && param !== undefined)
    }
    
    // 如果是字符串，尝试解析为 JSON
    if (typeof tplParams === 'string') {
      try {
        const parsed = JSON.parse(tplParams)
        if (Array.isArray(parsed)) {
          return parsed.filter(param => param !== null && param !== undefined)
        }
      } catch (e) {
        // 如果解析失败，可能是逗号分隔的字符串
        return tplParams.split(',').map(s => s.trim()).filter(s => s)
      }
    }
    
    return []
  }
  
  // 查询数据
  const onSubmit = () => {
    setPage(1)
    fetchPersonalEmailList()
  }
  
  // 重置搜索
  const onReset = () => {
    resetSearchInfo()
    setPage(1)
    fetchPersonalEmailList()
  }
  
  // 分页处理
  const handleCurrentChange = (val) => {
    setPage(val)
    fetchPersonalEmailList()
  }
  
  const handleSizeChange = (val) => {
    setPageSize(val)
    fetchPersonalEmailList()
  }

  
    const elFormRef = ref()

    // 行为控制标记（弹窗内部需要增还是改）
    const type = ref('')

    // 弹窗控制标记
    const dialogFormVisible = ref(false)

    // 使用多语言 Composable
    const { activeTitleTab, activeContentTab, resetActiveTabs } = useMultilingual()

    // 附件资源列表缓存（每个类型对应一个资源列表）
    const attachmentResourceLists = ref({})

    // 验证规则
    const rule = reactive({
      player_id: [
        { required: true, message: '请输入玩家PlayerID', trigger: 'blur' }
      ],
      senderI18n: [
        { 
          validator: (rule, value, callback) => {
            if (!value || typeof value !== 'object') {
              callback(new Error('请输入邮件发件人'))
              return
            }
            // 检查是否至少有一个语言版本有值
            const hasValue = Object.values(value).some(v => v && v.trim())
            if (!hasValue) {
              callback(new Error('请输入至少一种语言的发件人'))
              return
            }
            callback()
          },
          trigger: 'blur'
        }
      ],
      titleI18n: [
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
      contentI18n: [
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
      type: [
        { required: true, message: '请选择邮件类型', trigger: 'change' }
      ]
    })

    // 自动化生成的字典（可能为空）以及字段
    const formData = ref({
      player_id: '',
      senderI18n: initSenderI18nDefault(),
      titleI18n: initMultilingualData(),
      contentI18n: initMultilingualData(),
      attachments: [],
      tpl_id: '',
      tpl_params: [],
      remark: '',
      startTime: null,
      endTime: null,
      type: '1'  // 默认选中类型1
    })


  // 打开弹窗
  const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
  }

  
  // 添加附件
  const addAttachment = () => {
    formData.value.attachments.push({
      id: null,
      type: null,
      num: 1
    })
  }

  // 删除附件
  const removeAttachment = (index) => {
    formData.value.attachments.splice(index, 1)
  }

  // 处理附件资源类型变化
  const handleAttachmentTypeChange = async (index, type) => {
    const attachment = formData.value.attachments[index]
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
      player_id: '',
      senderI18n: initSenderI18nDefault(),
      titleI18n: initMultilingualData(),
      contentI18n: initMultilingualData(),
      attachments: [],
      tpl_id: '',
      tpl_params: [],
      remark: '',
      startTime: null,
      endTime: null,
      type: '1'  // 默认选中类型1
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
        const validAttachments = formData.value.attachments.filter(
          att => att.type !== null && att.id !== null && att.num > 0
        )
        
        // 构建提交数据
        const submitData = {
          ...formData.value,
          attachments: validAttachments
        }
        
        switch (type.value) {
          case 'create':
            await sendPersonalEmail(submitData)
            break
          case 'update':
            // await updatePersonalEmail(submitData)
            break
          default:
            await sendPersonalEmail(submitData)
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
        fetchPersonalEmailList()
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
  