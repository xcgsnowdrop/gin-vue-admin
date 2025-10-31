/**
 * 多语言相关的 Composable
 * 提供语言选项、多语言数据初始化等公共功能
 */

import { ref } from 'vue'

// 语言选项配置
export const languageOptions = [
  { code: 'en', label: '英文' },
  { code: 'zh-TW', label: '繁中' },
  { code: 'ja', label: '日文' },
  { code: 'ko', label: '韩文' }
]

// 语言代码到语言名称的映射
export const languageMap = {
  'en': '英文',
  'zh-TW': '繁中',
  'ja': '日文',
  'ko': '韩文'
}

// 语言优先级顺序（用于选择默认显示语言）
export const languagePriority = ['ja', 'en', 'zh-TW', 'ko']

/**
 * 初始化多语言数据对象
 * @returns {Object} 包含所有语言代码的空字符串对象
 */
export const initMultilingualData = () => {
  const multilingual = {}
  languageOptions.forEach(lang => {
    multilingual[lang.code] = ''
  })
  return multilingual
}

/**
 * 初始化发件人默认值（各语言的"GM系统管理员"）
 * @returns {Object} 包含各语言默认值的对象
 */
export const initSenderI18nDefault = () => {
  return {
    'en': 'GM System Administrator',
    'zh-TW': 'GM系統管理員',
    'ja': 'GMシステム管理者',
    'ko': 'GM 시스템 관리자'
  }
}

/**
 * 多语言 Composable
 * 提供多语言表单相关的响应式状态和方法
 * @returns {Object} 包含多语言相关的状态和方法
 */
export const useMultilingual = () => {
  // 多语言表单标签页活动状态
  const activeTitleTab = ref('en')
  const activeContentTab = ref('en')

  /**
   * 重置活动标签页
   */
  const resetActiveTabs = () => {
    activeTitleTab.value = 'en'
    activeContentTab.value = 'en'
  }

  /**
   * 设置默认活动标签页为第一个有内容的语言
   * @param {Object} data - 包含多语言数据的对象
   * @param {Object} data.title - 标题的多语言对象
   * @param {Object} data.content - 内容的多语言对象
   */
  const setActiveTabsFromData = (data) => {
    if (data.title) {
      for (const lang of languagePriority) {
        if (data.title[lang] && data.title[lang].trim()) {
          activeTitleTab.value = lang
          break
        }
      }
    }
    if (data.content) {
      for (const lang of languagePriority) {
        if (data.content[lang] && data.content[lang].trim()) {
          activeContentTab.value = lang
          break
        }
      }
    }
  }

  return {
    activeTitleTab,
    activeContentTab,
    resetActiveTabs,
    setActiveTabsFromData
  }
}

