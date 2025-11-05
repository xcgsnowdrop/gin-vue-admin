<template>
  <div class="attachments-form">
    <div
      v-for="(attachment, index) in modelValue"
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
          v-for="resType in resourceTypesList"
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
</template>

<script setup>
import { ref, watch, computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  resourceTypes: {
    type: Array,
    default: () => []
  },
  resourceList: {
    type: Array,
    default: () => []
  },
  resourceMap: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue', 'fetch-resource-list'])

// 附件资源列表缓存（每个类型对应一个资源列表）
const attachmentResourceLists = ref({})

// 确保 resourceTypes 是响应式的
// 注意：props 在 Vue 3 中已经是响应式的，但为了确保正确渲染，使用 computed
const resourceTypesList = computed(() => {
  const types = Array.isArray(props.resourceTypes) ? props.resourceTypes : []
  // 调试：检查数据结构
  if (process.env.NODE_ENV === 'development') {
    if (types.length === 0 && props.resourceTypes) {
      console.warn('AttachmentForm: resourceTypes 不是数组格式:', props.resourceTypes)
    }
  }
  return types
})

// 根据资源类型获取资源列表
// 注意：这是一个函数，在模板中调用时会自动响应 props.resourceMap 的变化
const getResourceListByType = (type) => {
  if (!type) return []
  // 优先从缓存中获取
  if (attachmentResourceLists.value[type]) {
    return attachmentResourceLists.value[type]
  }
  // 如果缓存中没有，从 props 的 resourceMap 中构建
  const resourceMap = props.resourceMap || {}
  if (resourceMap[type]) {
    return Object.keys(resourceMap[type]).map(id => ({
      id: parseInt(id),
      name: resourceMap[type][id]
    }))
  }
  return []
}

// 添加附件
const addAttachment = () => {
  const newAttachments = [...props.modelValue, {
    id: null,
    type: null,
    num: 1
  }]
  emit('update:modelValue', newAttachments)
}

// 删除附件
const removeAttachment = (index) => {
  const newAttachments = props.modelValue.filter((_, i) => i !== index)
  emit('update:modelValue', newAttachments)
}

// 处理附件资源类型变化
const handleAttachmentTypeChange = async (index, type) => {
  const attachment = props.modelValue[index]
  // 清空资源ID选择
  const newAttachments = [...props.modelValue]
  newAttachments[index] = { ...attachment, id: null, type }
  emit('update:modelValue', newAttachments)
  
  if (type) {
    // 如果该类型的资源列表还未加载，则触发加载事件
    if (!attachmentResourceLists.value[type]) {
      emit('fetch-resource-list', type, (list) => {
        attachmentResourceLists.value[type] = list
      })
    }
  }
}

// 监听 resourceList 变化，更新缓存
watch(() => props.resourceList, (newList) => {
  if (newList && newList.length > 0) {
    // 根据 resourceList 中的第一个资源的 type 来更新缓存
    const type = newList[0]?.type
    if (type) {
      attachmentResourceLists.value[type] = newList
    }
  }
}, { deep: true })
</script>

<style scoped>
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

