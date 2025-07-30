<template>
  <v-dialog v-model="dialogVisible" max-width="500px" @update:model-value="handleDialogClose">
    <v-card>
      <v-card-title>测试通知 - {{ app?.name }}</v-card-title>
      <v-card-text>
        <v-form>
          <v-text-field v-model="form.title" label="标题" class="mb-4"></v-text-field>
          <v-textarea v-model="form.content" label="内容" rows="3" class="mb-4"></v-textarea>
          <v-text-field v-model="form.image" label="图片" class="mb-4"></v-text-field>
          <v-text-field v-model="form.url" label="链接URL"></v-text-field>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text="取消" @click="handleCancel"></v-btn>
        <v-btn variant="flat" color="primary" text="发送" @click="handleSend" :loading="loading"></v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue'

interface TestNotificationData {
  title: string
  content: string
  url: string
  image: string
}

interface Props {
  app?: any
  loading?: boolean
}

interface Emits {
  (e: 'send', data: TestNotificationData): void
  (e: 'cancel'): void
}

withDefaults(defineProps<Props>(), {
  loading: false
})

const dialogVisible = defineModel<boolean>()

const emit = defineEmits<Emits>()

// 表单数据
const form = ref<TestNotificationData>({
  title: '测试通知',
  content: '这是一条测试通知消息',
  url: '',
  image: ''
})

// 监听对话框打开时重置表单
watch(dialogVisible, (visible) => {
  if (visible) {
    form.value = {
      title: '测试通知',
      content: '这是一条测试通知消息',
      url: '',
      image: ''
    }
  }
})

// 处理发送
const handleSend = () => {
  emit('send', { ...form.value })
}

// 处理取消
const handleCancel = () => {
  emit('cancel')
}

// 处理对话框关闭
const handleDialogClose = (value: boolean) => {
  if (!value) {
    emit('cancel')
  }
}
</script>
