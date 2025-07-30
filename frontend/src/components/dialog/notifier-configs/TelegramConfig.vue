<template>
  <div>
    <v-text-field v-model="config.bot_token" label="Bot Token *" :rules="[rules.required]" hint="从 @BotFather 获取"
      persistent-hint class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.chat_id" label="Chat ID *" :rules="[rules.required]" hint="群组或频道ID，可以是负数"
      persistent-hint class="mb-4" @input="handleConfigChange"></v-text-field>

    <v-text-field v-model="config.proxy" label="代理服务器" hint="可选，格式: http://proxy.example.com:8080" persistent-hint
      class="mb-4" @input="handleConfigChange"></v-text-field>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { TelegramConfig } from '@/common/types'

interface Props {
  modelValue: Partial<TelegramConfig>
}

interface Emits {
  (e: 'update:modelValue', value: Partial<TelegramConfig>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 内部配置状态
const config = ref<Partial<TelegramConfig>>({
  bot_token: '',
  chat_id: '',
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
    bot_token: '',
    chat_id: '',
    proxy: '',
    ...newValue
  }
}, { deep: true })

// 配置变化处理
const handleConfigChange = () => {
  emit('update:modelValue', { ...config.value })
}
</script>
