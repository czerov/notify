<template>
  <div>
    <v-text-field v-model="config.corp_id" label="企业ID *" :rules="[rules.required]" class="mb-4"
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.agent_id" label="应用ID *" :rules="[rules.required]" class="mb-4"
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.secret" label="应用密钥 *" :rules="[rules.required]" class="mb-4"
      @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.proxy" label="代理服务器" hint="可选，格式: http://proxy.example.com:8080" persistent-hint
      class="mb-4" @input="handleConfigChange"></v-text-field>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { WechatWorkConfig } from '@/common/types'

interface Props {
  modelValue: Partial<WechatWorkConfig>
}

interface Emits {
  (e: 'update:modelValue', value: Partial<WechatWorkConfig>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 内部配置状态
const config = ref<Partial<WechatWorkConfig>>({
  corp_id: '',
  agent_id: '',
  secret: '',
  proxy: '',
  ...props.modelValue
})

// 验证规则
const rules = {
  required: (value: any) => !!value || '此字段为必填项'
}

// 监听 props 变化
watch(() => props.modelValue, (newValue) => {
  config.value = {
    corp_id: '',
    agent_id: '',
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
