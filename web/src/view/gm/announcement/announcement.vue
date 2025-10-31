<template>
    <div>
      <!-- <div class="gva-search-box">
        <el-form
          ref="elSearchFormRef"
          :inline="true"
          :model="searchInfo"
          class="demo-form-inline"
          :rules="searchRule"
          @keyup.enter="onSubmit"
        >
          <el-form-item label="创建时间" prop="createdAt">
            <template #label>
              <span>
                创建时间
                <el-tooltip
                  content="搜索范围是开始时间（包含）至结束时间（不包含）"
                >
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </span>
            </template>
            <el-date-picker
              v-model="searchInfo.startCreatedAt"
              type="datetime"
              placeholder="开始时间"
              :disabled-date="
                (time) =>
                  searchInfo.endCreatedAt
                    ? time.getTime() > searchInfo.endCreatedAt.getTime()
                    : false
              "
            />
            —
            <el-date-picker
              v-model="searchInfo.endCreatedAt"
              type="datetime"
              placeholder="结束时间"
              :disabled-date="
                (time) =>
                  searchInfo.startCreatedAt
                    ? time.getTime() < searchInfo.startCreatedAt.getTime()
                    : false
              "
            />
          </el-form-item>
  
          <template v-if="showAllQuery">
           // 将需要控制显示状态的查询条件添加到此范围内
          </template>
  
          <el-form-item>
            <el-button type="primary" icon="search" @click="onSubmit">
              查询
            </el-button>
            <el-button icon="refresh" @click="onReset"> 重置 </el-button>
            <el-button
              v-if="!showAllQuery"
              link
              type="primary"
              icon="arrow-down"
              @click="showAllQuery = true"
            >
              展开
            </el-button>
            <el-button
              v-else
              link
              type="primary"
              icon="arrow-up"
              @click="showAllQuery = false"
            >
              收起
            </el-button>
          </el-form-item>
        </el-form>
      </div> -->
      <div class="gva-table-box">
        <div class="gva-btn-list">
          <el-button type="primary" icon="plus" @click="openDialog">
            新增
          </el-button>
          <el-button
            icon="delete"
            style="margin-left: 10px"
            :disabled="!multipleSelection.length"
            @click="onDelete"
          >
            删除
          </el-button>
        </div>
        <el-table
          ref="multipleTable"
          style="width: 100%"
          tooltip-effect="dark"
          :data="tableData"
          row-key="index"
          v-loading="loading"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
  
          <el-table-column align="left" label="排序" prop="index" width="80" />
          <el-table-column align="left" label="标题" width="200">
            <template #default="scope">
              <div class="multilingual-cell">
                <span class="main-text">{{ getMultilingualText(scope.row.title) }}</span>
                <el-popover
                  placement="top"
                  :width="300"
                  trigger="hover"
                  popper-class="multilingual-popover"
                >
                  <template #reference>
                    <div class="language-tags">
                      <el-tag
                        v-for="lang in getAvailableLanguages(scope.row.title)"
                        :key="lang.code"
                        size="small"
                        :type="lang.isDefault ? 'primary' : 'info'"
                        class="lang-tag"
                      >
                        {{ lang.label }}
                      </el-tag>
                    </div>
                  </template>
                  <div class="multilingual-detail">
                    <div
                      v-for="lang in getAvailableLanguages(scope.row.title)"
                      :key="lang.code"
                      class="lang-item"
                    >
                      <el-tag size="small" :type="lang.isDefault ? 'primary' : 'info'">
                        {{ lang.label }}
                      </el-tag>
                      <span class="lang-text">{{ lang.text || '(空)' }}</span>
                    </div>
                  </div>
                </el-popover>
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left" label="内容" width="250">
            <template #default="scope">
              <div class="multilingual-cell">
                <span class="main-text">{{ getMultilingualText(scope.row.content, 30) }}</span>
                <el-popover
                  placement="top"
                  :width="300"
                  trigger="hover"
                  popper-class="multilingual-popover"
                >
                  <template #reference>
                    <div class="language-tags">
                      <el-tag
                        v-for="lang in getAvailableLanguages(scope.row.content)"
                        :key="lang.code"
                        size="small"
                        :type="lang.isDefault ? 'primary' : 'info'"
                        class="lang-tag"
                      >
                        {{ lang.label }}
                      </el-tag>
                    </div>
                  </template>
                  <div class="multilingual-detail">
                    <div
                      v-for="lang in getAvailableLanguages(scope.row.content)"
                      :key="lang.code"
                      class="lang-item"
                    >
                      <el-tag size="small" :type="lang.isDefault ? 'primary' : 'info'">
                        {{ lang.label }}
                      </el-tag>
                      <div class="lang-text">{{ lang.text || '(空)' }}</div>
                    </div>
                  </div>
                </el-popover>
              </div>
            </template>
          </el-table-column>
          <el-table-column align="left" label="开始时间" prop="start_time_formatted" width="180" />
          <el-table-column align="left" label="结束时间" prop="end_time_formatted" width="180" />
          <el-table-column align="left" label="创建时间" prop="create_time_formatted" width="180" />
          <el-table-column align="left" label="公告类型" width="120">
            <template #default="scope">
              <el-tag :type="getAnnouncementTypeTag(scope.row.type)" size="default">
                {{ getAnnouncementTypeName(scope.row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column align="left" label="是否展示" width="100">
            <template #default="scope">
              <el-tag 
                :type="getShowStatusTag(scope.row.isShow)" 
                size="default"
                effect="plain"
              >
                <el-icon v-if="scope.row.isShow" style="margin-right: 4px;">
                  <View />
                </el-icon>
                <el-icon v-else style="margin-right: 4px;">
                  <Hide />
                </el-icon>
                {{ getShowStatusText(scope.row.isShow) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            align="left"
            label="操作"
            fixed="right"
            min-width="300"
          >
            <template #default="scope">
              <el-button
                type="primary"
                link
                icon="edit"
                class="table-button"
                @click="updateRow(scope.row)"
              >
                变更
              </el-button>
              <el-button
                type="warning"
                link
                icon="top"
                :disabled="scope.row.index === 1"
                @click="toppingRow(scope.row)"
              >
                置顶
              </el-button>
              <el-button
                type="primary"
                link
                icon="delete"
                @click="deleteRow(scope.row)"
              >
                删除
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
          <el-form-item label="标题:" prop="title">
            <el-tabs v-model="activeTitleTab" type="border-card" class="multilingual-tabs">
              <el-tab-pane
                v-for="lang in languageOptions"
                :key="lang.code"
                :label="lang.label"
                :name="lang.code"
              >
                <el-input
                  v-model="formData.title[lang.code]"
                  :clearable="true"
                  :placeholder="`请输入${lang.label}标题`"
                />
              </el-tab-pane>
            </el-tabs>
          </el-form-item>
          <el-form-item label="内容:" prop="content">
            <el-tabs v-model="activeContentTab" type="border-card" class="multilingual-tabs">
              <el-tab-pane
                v-for="lang in languageOptions"
                :key="lang.code"
                :label="lang.label"
                :name="lang.code"
              >
                <RichEdit v-model="formData.content[lang.code]" />
              </el-tab-pane>
            </el-tabs>
          </el-form-item>
          <el-form-item label="开始时间:" prop="startTime">
            <el-date-picker
              v-model="formData.startTime"
              type="datetime"
              placeholder="请选择开始时间"
            />
          </el-form-item>
          <el-form-item label="结束时间:" prop="endTime">
            <el-date-picker
              v-model="formData.endTime"
              type="datetime"
              placeholder="请选择结束时间"
              :disabled-date="disabledEndDate"
            />
          </el-form-item>
          <el-form-item label="类型:" prop="type">
            <el-select
              v-model="formData.type"
              placeholder="请选择类型"
              style="width: 100%"
              :clearable="true"
            >
              <el-option label="版本更新" value="1" />
              <el-option label="活动预告" value="2" />
              <el-option label="FAQ" value="3" />
            </el-select>
          </el-form-item>
          <el-form-item label="是否展示:" prop="isShow">
            <el-switch
              v-model="formData.isShow"
              :active-value="true"
              :inactive-value="false"
            />
          </el-form-item>
        </el-form>
      </el-drawer>
    </div>
  </template>
  
  <script setup>
    // 富文本组件
    import RichEdit from '@/components/richtext/rich-edit.vue'
    // 文件选择组件
    import { useGMAnnouncementStore } from '@/pinia/gm/announcement'
    import { storeToRefs } from 'pinia'
  
    // 全量引入格式化工具 请按需保留
    import { ElMessage, ElMessageBox } from 'element-plus'
    import { ref, reactive, onMounted, watch } from 'vue'
    import { View, Hide } from '@element-plus/icons-vue'
  
    defineOptions({
      name: 'GmAnnouncement'
    })
  
    const gmAnnouncementStore = useGMAnnouncementStore()
    // 使用store中的状态
    const { 
      announcementList: tableData, 
      loading, 
      total, 
      page,
      pageSize,
    //   searchInfo,
    } = storeToRefs(gmAnnouncementStore)

    // 使用store中的方法
    const { 
        fetchAnnouncementList,
        addAnnouncement,
        deleteAnnouncement,
        updateAnnouncement,
        toppingAnnouncement,
        // resetSearchInfo,
        setPage,
        setPageSize
    } = gmAnnouncementStore

    // 控制更多查询条件显示/隐藏状态
    // const showAllQuery = ref(false)

    // 语言选项配置
    const languageOptions = [
      { code: 'en', label: '英文' },
      { code: 'zh-TW', label: '繁中' },
      { code: 'ja', label: '日文' },
      { code: 'ko', label: '韩文' }
    ]

    // 多语言表单标签页活动状态
    const activeTitleTab = ref('en')
    const activeContentTab = ref('en')
  
    // 初始化多语言数据对象
    const initMultilingualData = () => {
      const multilingual = {}
      languageOptions.forEach(lang => {
        multilingual[lang.code] = ''
      })
      return multilingual
    }

    // 自动化生成的字典（可能为空）以及字段
    const formData = ref({
      title: initMultilingualData(),
      content: initMultilingualData(),
      startTime: null,
      endTime: null,
      type: '',
    //   id: undefined,
      isShow: 1
    })

    // 禁用结束时间：必须在今天及开始时间（若有）之后的日期
    const disabledEndDate = (time) => {
      const todayStart = new Date()
      todayStart.setHours(0, 0, 0, 0)
      let minDateMs = todayStart.getTime()
      if (formData.value.startTime) {
        const start = new Date(formData.value.startTime)
        const startDateOnly = new Date(start.getFullYear(), start.getMonth(), start.getDate())
        minDateMs = Math.max(minDateMs, startDateOnly.getTime())
      }
      return time.getTime() < minDateMs
    }

    // 验证规则
    const rule = reactive({
      endTime: [
        {
          validator: (rule, value, callback) => {
            if (!value) {
              callback()
              return
            }
            const endMs = new Date(value).getTime()
            const nowMs = Date.now()
            if (endMs <= nowMs) {
              callback(new Error('结束时间必须大于当前时间'))
              return
            }
            if (formData.value.startTime) {
              const startMs = new Date(formData.value.startTime).getTime()
              if (endMs <= startMs) {
                callback(new Error('结束时间必须大于开始时间'))
                return
              }
            }
            callback()
          },
          trigger: 'change'
        }
      ]
    })
  
    // const searchRule = reactive({
    //   createdAt: [
    //     {
    //       validator: (rule, value, callback) => {
    //         if (
    //           searchInfo.value.startCreatedAt &&
    //           !searchInfo.value.endCreatedAt
    //         ) {
    //           callback(new Error('请填写结束日期'))
    //         } else if (
    //           !searchInfo.value.startCreatedAt &&
    //           searchInfo.value.endCreatedAt
    //         ) {
    //           callback(new Error('请填写开始日期'))
    //         } else if (
    //           searchInfo.value.startCreatedAt &&
    //           searchInfo.value.endCreatedAt &&
    //           (searchInfo.value.startCreatedAt.getTime() ===
    //             searchInfo.value.endCreatedAt.getTime() ||
    //             searchInfo.value.startCreatedAt.getTime() >
    //               searchInfo.value.endCreatedAt.getTime())
    //         ) {
    //           callback(new Error('开始日期应当早于结束日期'))
    //         } else {
    //           callback()
    //         }
    //       },
    //       trigger: 'change'
    //     }
    //   ]
    // })
  
    const elFormRef = ref()
    // const elSearchFormRef = ref()
  
    // 重置
    // const onReset = () => {
    //   resetSearchInfo()
    //   setPage(1)
    //   fetchAnnouncementList()
    // }
  
    // 搜索
    // const onSubmit = () => {
    //   elSearchFormRef.value?.validate(async (valid) => {
    //     if (!valid) return
    //     page.value = 1
    //     fetchAnnouncementList()
    //   })
    // }
  
    // 分页
    const handleSizeChange = (val) => {
      setPageSize(val)
      fetchAnnouncementList()
    }
  
    // 修改页面容量
    const handleCurrentChange = (val) => {
      setPage(val)
      fetchAnnouncementList()
    }

  
    // ============== 表格控制部分结束 ===============
  
    // 获取需要的字典 可能为空 按需保留
    const setOptions = async () => {}
  
    // 获取需要的字典 可能为空 按需保留
    setOptions()
  
    // 多选数据
    const multipleSelection = ref([])
    // 多选
    const handleSelectionChange = (val) => {
      multipleSelection.value = val
    }
  
    // 删除行
    const deleteRow = (row) => {
      ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteInfoFunc(row)
      })
    }

    // 置顶行
    const toppingRow = (row) => {
      ElMessageBox.confirm('确定要置顶此公告吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        toppingInfoFunc(row)
      })
    }
  
    // 多选删除
    const onDelete = async () => {
      ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        const indexes = []
        if (multipleSelection.value.length === 0) {
          ElMessage({
            type: 'warning',
            message: '请选择要删除的数据'
          })
          return
        }
        multipleSelection.value &&
          multipleSelection.value.map((item) => {
            indexes.push(item.index)
          })
        
        try {
          await deleteAnnouncement({indexes: indexes})
          ElMessage({
            type: 'success',
            message: '删除成功'
          })
          // 如果删除后当前页没有数据且不是第一页，则回到上一页
          if (tableData.value.length === indexes.length && page.value > 1) {
            setPage(page.value - 1)
          }
        } catch (error) {
          ElMessage({
            type: 'error',
            message: error.message || '删除失败'
          })
        }
      })
    }
  
    // 行为控制标记（弹窗内部需要增还是改）
    const type = ref('')
  
    // 更新行
    const updateRow = async (row) => {
        type.value = 'update'
        // 确保多语言数据格式正确
        const data = { ...row }
      
        // 确保所有语言字段都存在
        languageOptions.forEach(lang => {
          if (!data.title[lang.code]) {
            data.title[lang.code] = ''
          }
          if (!data.content[lang.code]) {
            data.content[lang.code] = ''
          }
        })

        // 转换时间戳为 Date 对象（用于日期选择器）
        if (data.startTime && typeof data.startTime === 'number') {
            data.startTime = new Date(data.startTime * 1000) // 秒级时间戳转 Date
        } else if (data.startTime && typeof data.startTime === 'string') {
            // 如果是字符串格式的时间戳
            const timestamp = Number(data.startTime)
            if (!isNaN(timestamp)) {
            data.startTime = new Date(timestamp * 1000)
            }
        }

        if (data.endTime && typeof data.endTime === 'number') {
            data.endTime = new Date(data.endTime * 1000) // 秒级时间戳转 Date
        } else if (data.endTime && typeof data.endTime === 'string') {
            // 如果是字符串格式的时间戳
            const timestamp = Number(data.endTime)
            if (!isNaN(timestamp)) {
            data.endTime = new Date(timestamp * 1000)
            }
        }

        formData.value = data
        
        // 设置默认活动标签页为第一个有内容的语言
        for (const lang of languageOptions) {
            if (data.title[lang.code] && data.title[lang.code].trim()) {
            activeTitleTab.value = lang.code
            break
            }
        }
        for (const lang of languageOptions) {
            if (data.content[lang.code] && data.content[lang.code].trim()) {
            activeContentTab.value = lang.code
            break
            }
        }
        
        dialogFormVisible.value = true
    }
  
    // 删除行
    const deleteInfoFunc = async (row) => {
      try {
        await deleteAnnouncement({ indexes: [row.index] })
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        // 如果删除后当前页没有数据且不是第一页，则回到上一页
        if (tableData.value.length === 1 && page.value > 1) {
          setPage(page.value - 1)
        }
      } catch (error) {
        ElMessage({
          type: 'error',
          message: error.message || '删除失败'
        })
      }
    }

    // 置顶行
    const toppingInfoFunc = async (row) => {
      try {
        await toppingAnnouncement({ index: row.index })
        ElMessage({
          type: 'success',
          message: '置顶成功'
        })
      } catch (error) {
        ElMessage({
          type: 'error',
          message: error.message || '置顶失败'
        })
      }
    }
  
    // 弹窗控制标记
    const dialogFormVisible = ref(false)
  
    // 打开弹窗
    const openDialog = () => {
      type.value = 'create'
      dialogFormVisible.value = true
    }
  
    // 关闭弹窗
    const closeDialog = () => {
      dialogFormVisible.value = false
      formData.value = {
        title: initMultilingualData(),
        content: initMultilingualData(),
        startTime: null,
        endTime: null,
        type: '',
        id: undefined,
        isShow: 1
      }
      activeTitleTab.value = 'en'
      activeContentTab.value = 'en'
    }
    // 弹窗确定
    const enterDialog = async () => {
      elFormRef.value?.validate(async (valid) => {
        if (!valid) return
        
        try {
          switch (type.value) {
            case 'create':
              await addAnnouncement(formData.value)
              break
            case 'update':
              await updateAnnouncement(formData.value)
              break
            default:
              await addAnnouncement(formData.value)
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

    // 公告类型映射
    const announcementTypeMap = {
      '1': { name: '版本更新', tagType: 'success' },
      '2': { name: '活动预告', tagType: 'warning' },
      '3': { name: 'FAQ', tagType: 'info' }
    }

    // 获取公告类型名称
    const getAnnouncementTypeName = (type) => {
      if (!type) return '未知'
      const typeStr = String(type)
      return announcementTypeMap[typeStr]?.name || '未知'
    }

    // 获取公告类型标签样式
    const getAnnouncementTypeTag = (type) => {
      if (!type) return 'info'
      const typeStr = String(type)
      return announcementTypeMap[typeStr]?.tagType || 'info'
    }

    // 获取展示状态文本
    const getShowStatusText = (isShow) => {
      // 兼容 true, false, 1, 0, "true", "false" 等情况
      if (isShow === true || isShow === 1 || isShow === '1' || isShow === 'true') {
        return '展示中'
      }
      return '已隐藏'
    }

    // 获取展示状态标签样式
    const getShowStatusTag = (isShow) => {
      if (isShow === true || isShow === 1 || isShow === '1' || isShow === 'true') {
        return 'success'
      }
      return 'info'
    }

    // 语言代码到语言名称的映射
    const languageMap = {
      'en': '英文',
      'zh-TW': '繁中',
      'ja': '日文',
      'ko': '韩文'
    }

    // 语言优先级顺序（用于选择默认显示语言）
    const languagePriority = ['ja', 'en', 'zh-TW', 'ko']

    // 获取多语言文本（返回默认语言的文本）
    const getMultilingualText = (multilingualObj, maxLength = null) => {
      if (!multilingualObj || typeof multilingualObj !== 'object') {
        return '-'
      }

      // 按优先级查找第一个有值的语言
      for (const lang of languagePriority) {
        if (multilingualObj[lang] && multilingualObj[lang].trim()) {
          let text = multilingualObj[lang]
          if (maxLength && text.length > maxLength) {
            text = text.substring(0, maxLength) + '...'
          }
          return text
        }
      }

      // 如果没有找到，返回第一个有值的语言
      for (const key in multilingualObj) {
        if (multilingualObj[key] && multilingualObj[key].trim()) {
          let text = multilingualObj[key]
          if (maxLength && text.length > maxLength) {
            text = text.substring(0, maxLength) + '...'
          }
          return text
        }
      }

      return '-'
    }

    // 获取所有可用语言列表
    const getAvailableLanguages = (multilingualObj) => {
      if (!multilingualObj || typeof multilingualObj !== 'object') {
        return []
      }

      const languages = []
      let defaultLanguage = null

      // 先找出默认语言（优先级最高的有值语言）
      for (const lang of languagePriority) {
        if (multilingualObj[lang] && multilingualObj[lang].trim()) {
          defaultLanguage = lang
          break
        }
      }

      // 如果没有找到，使用第一个有值的语言作为默认语言
      if (!defaultLanguage) {
        for (const key in multilingualObj) {
          if (multilingualObj[key] && multilingualObj[key].trim()) {
            defaultLanguage = key
            break
          }
        }
      }

      // 构建语言列表
      for (const lang of languagePriority) {
        if (multilingualObj[lang]) {
          languages.push({
            code: lang,
            label: languageMap[lang] || lang.toUpperCase(),
            text: multilingualObj[lang],
            isDefault: lang === defaultLanguage
          })
        }
      }

      // 如果有其他语言不在优先级列表中，也添加进去
      for (const key in multilingualObj) {
        if (!languagePriority.includes(key) && multilingualObj[key]) {
          languages.push({
            code: key,
            label: key.toUpperCase(),
            text: multilingualObj[key],
            isDefault: key === defaultLanguage
          })
        }
      }

      return languages
    }

    watch(
      () => tableData.value,
      (newValue, oldValue) => {
        console.log('tableData 变化了')
        console.log('新值:', newValue)
        console.log('旧值:', oldValue)
      },
      { deep: true, immediate: true }
    )

    onMounted(async () => {
      try {
        await Promise.all([
          fetchAnnouncementList()
        ])
    } catch (error) {
      console.error('初始化失败:', error)
    }
  })
  </script>
  
  <style>
    .file-list {
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
    }
  
    .fileBtn {
      margin-bottom: 10px;
    }
  
    .fileBtn:last-child {
      margin-bottom: 0;
    }

    /* 多语言单元格样式 */
    .multilingual-cell {
      display: flex;
      flex-direction: column;
      gap: 4px;
    }

    .multilingual-cell .main-text {
      font-weight: 500;
      color: #303133;
      word-break: break-word;
    }

    .multilingual-cell .language-tags {
      display: flex;
      flex-wrap: wrap;
      gap: 4px;
      margin-top: 2px;
    }

    .multilingual-cell .lang-tag {
      cursor: pointer;
      font-size: 11px;
      height: 18px;
      line-height: 16px;
      padding: 0 6px;
    }

    /* 多语言详情弹窗样式 */
    .multilingual-detail {
      max-height: 400px;
      overflow-y: auto;
    }

    .multilingual-detail .lang-item {
      display: flex;
      align-items: flex-start;
      margin-bottom: 12px;
      padding-bottom: 12px;
      border-bottom: 1px solid #ebeef5;
    }

    .multilingual-detail .lang-item:last-child {
      margin-bottom: 0;
      padding-bottom: 0;
      border-bottom: none;
    }

    .multilingual-detail .lang-item .el-tag {
      margin-right: 8px;
      flex-shrink: 0;
      min-width: 50px;
      text-align: center;
    }

    .multilingual-detail .lang-text {
      flex: 1;
      color: #606266;
      word-break: break-word;
      line-height: 1.5;
      white-space: pre-wrap;
    }

    /* 弹窗样式优化 */
    :deep(.multilingual-popover) {
      max-width: 350px;
    }

    /* 多语言表单标签页样式 */
    .multilingual-tabs {
      margin-top: 8px;
    }

    :deep(.multilingual-tabs .el-tabs__header) {
      margin-bottom: 16px;
    }

    :deep(.multilingual-tabs .el-tabs__item) {
      padding: 0 20px;
      font-size: 14px;
    }

    :deep(.multilingual-tabs .el-tab-pane) {
      padding: 0;
    }
  </style>
  