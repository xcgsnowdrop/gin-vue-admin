<template>
  <div v-if="attachments && attachments.length > 0" class="attachment-list-vertical">
    <el-tag
      v-for="(attachment, index) in attachments"
      :key="index"
      size="small"
      type="success"
      class="attachment-tag"
    >
      {{ formatAttachmentFunc(attachment) }}
    </el-tag>
  </div>
  <span v-else>-</span>
</template>

<script setup>
import { computed } from 'vue'
import { useResource } from '@/composables/useResource'

const props = defineProps({
  attachments: {
    type: Array,
    default: () => []
  },
  formatAttachment: {
    type: Function,
    default: null
  }
})

// 如果没有传入 formatAttachment，使用默认的格式化函数（从 composable 获取）
// 注意：这仍然会创建一个新的状态实例，但只在没有传入 formatAttachment 时使用
const { formatAttachment: formatAttachmentDefault } = useResource()

// 使用传入的 formatAttachment 或默认的
const formatAttachmentFunc = computed(() => props.formatAttachment || formatAttachmentDefault)
</script>

<style scoped>
/* 附件列表样式 - 上下排列 */
.attachment-list-vertical {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.attachment-tag {
  margin: 0;
  width: fit-content;
}
</style>

