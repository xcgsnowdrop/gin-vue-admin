<template>
  <el-form-item :label="label" :prop="prop" :required="required">
    <el-tabs v-model="localActiveTab" type="border-card" class="multilingual-tabs">
      <el-tab-pane
        v-for="lang in languageOptions"
        :key="lang.code"
        :label="lang.label"
        :name="lang.code"
      >
        <RichEdit v-model="localValue[lang.code]" />
      </el-tab-pane>
    </el-tabs>
  </el-form-item>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import RichEdit from '@/components/richtext/rich-edit.vue'
import { languageOptions } from '@/composables/useMultilingual'

const props = defineProps({
  // 表单标签
  label: {
    type: String,
    required: true
  },
  // 表单验证 prop
  prop: {
    type: String,
    default: ''
  },
  // 是否必填
  required: {
    type: Boolean,
    default: false
  },
  // 多语言对象，格式如 { en: 'English', ja: '日本語', ... }
  modelValue: {
    type: Object,
    default: () => ({})
  },
  // 活动标签页（可选，用于同步多个输入组件的标签页）
  activeTab: {
    type: String,
    default: 'en'
  }
})

const emit = defineEmits(['update:modelValue', 'update:activeTab'])

// 本地活动标签页状态
const localActiveTab = ref(props.activeTab || 'en')

// 同步 props.activeTab 的变化
watch(
  () => props.activeTab,
  (newVal) => {
    if (newVal) {
      localActiveTab.value = newVal
    }
  },
  { immediate: true }
)

// 当本地标签页变化时，通知父组件
watch(localActiveTab, (newVal) => {
  emit('update:activeTab', newVal)
})

// 本地值，用于双向绑定
const localValue = computed({
  get: () => props.modelValue,
  set: (value) => {
    emit('update:modelValue', value)
  }
})
</script>

<style scoped>
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

