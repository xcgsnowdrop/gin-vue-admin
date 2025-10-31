<template>
  <div class="multilingual-cell">
    <span class="main-text">{{ displayText }}</span>
    <el-popover
      v-if="hasMultipleLanguages"
      placement="top"
      :width="300"
      trigger="hover"
      popper-class="multilingual-popover"
    >
      <template #reference>
        <div class="language-tags">
          <el-tag
            v-for="lang in availableLanguages"
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
          v-for="lang in availableLanguages"
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

<script setup>
import { computed } from 'vue'
import { ElTag, ElPopover } from 'element-plus'

const props = defineProps({
  // 多语言对象，格式如 { en: 'English', ja: '日本語', ... }
  value: {
    type: [Object, String],
    default: null
  },
  // 最大显示长度，超过则截断
  maxLength: {
    type: Number,
    default: null
  },
  // 语言优先级顺序
  languagePriority: {
    type: Array,
    default: () => ['ja', 'en', 'zh-TW', 'ko']
  },
  // 语言代码到语言名称的映射
  languageMap: {
    type: Object,
    default: () => ({
      'en': '英文',
      'zh-TW': '繁中',
      'ja': '日文',
      'ko': '韩文'
    })
  }
})

// 获取多语言文本（返回默认语言的文本）
const getMultilingualText = (multilingualObj) => {
  if (!multilingualObj || typeof multilingualObj !== 'object') {
    return '-'
  }

  // 按优先级查找第一个有值的语言
  for (const lang of props.languagePriority) {
    if (multilingualObj[lang] && multilingualObj[lang].trim()) {
      let text = multilingualObj[lang]
      if (props.maxLength && text.length > props.maxLength) {
        text = text.substring(0, props.maxLength) + '...'
      }
      return text
    }
  }

  // 如果没有找到，返回第一个有值的语言
  for (const key in multilingualObj) {
    if (multilingualObj[key] && multilingualObj[key].trim()) {
      let text = multilingualObj[key]
      if (props.maxLength && text.length > props.maxLength) {
        text = text.substring(0, props.maxLength) + '...'
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
  for (const lang of props.languagePriority) {
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
  for (const lang of props.languagePriority) {
    if (multilingualObj[lang]) {
      languages.push({
        code: lang,
        label: props.languageMap[lang] || lang.toUpperCase(),
        text: multilingualObj[lang],
        isDefault: lang === defaultLanguage
      })
    }
  }

  // 如果有其他语言不在优先级列表中，也添加进去
  for (const key in multilingualObj) {
    if (!props.languagePriority.includes(key) && multilingualObj[key]) {
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

// 计算显示的文本
const displayText = computed(() => {
  return getMultilingualText(props.value)
})

// 获取可用语言列表
const availableLanguages = computed(() => {
  return getAvailableLanguages(props.value)
})

// 判断是否有多个语言版本
const hasMultipleLanguages = computed(() => {
  if (!props.value || typeof props.value !== 'object') {
    return false
  }
  const languages = Object.keys(props.value).filter(key => props.value[key])
  return languages.length > 1
})
</script>

<style scoped>
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
</style>

<style>
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
.multilingual-popover {
  max-width: 350px;
}
</style>

