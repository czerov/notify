<template>
  <div>
    <v-text-field v-model="config.access_token" label="Access Token *" :rules="[rules.required]" class="mb-4"
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.secret" label="签名密钥 *" class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.proxy" label="代理服务器" hint="可选，格式: http://proxy.example.com:8080" persistent-hint
      class="mb-4" @input="handleConfigChange"></v-text-field>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { DingTalkConfig } from '@/common/types'

interface Props {
  modelValue: Partial<DingTalkConfig>
}

interface Emits {
  (e: 'update:modelValue', value: Partial<DingTalkConfig>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 内部配置状态
const config = ref<Partial<DingTalkConfig>>({
  access_token: '',
  secret: '',
  proxy: '',
  ...props.modelValue
})

// 验证规则
const rules = {
  required: (value: any) => !!value || '此字段为必填项',
  url: (value: string) => {
    if (!value) return true
    try {
      new URL(value)
      return true
    } catch {
      return '请输入有效的URL'
    }
  }
}

// 监听 props 变化
watch(() => props.modelValue, (newValue) => {
  config.value = {
    access_token: '',
    secret: '',
    proxy: '',
    ...newValue
  }
}, { deep: true })

// 配置变化处理
const handleConfigChange = () => {
  emit('update:modelValue', { ...config.value })
}
</script>
