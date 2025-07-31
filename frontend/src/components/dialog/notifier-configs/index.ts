// 通知服务配置组件
export { default as WechatWorkConfig } from './WechatWorkConfig.vue'
export { default as WechatWorkWebhookConfig } from './WechatWorkWebhookConfig.vue'
export { default as TelegramConfig } from './TelegramConfig.vue'
export { default as DingtalkConfig } from './DingtalkConfig.vue'

// 组件映射
import WechatWorkConfig from './WechatWorkConfig.vue'
import WechatWorkWebhookConfig from './WechatWorkWebhookConfig.vue'
import TelegramConfig from './TelegramConfig.vue'
import DingtalkConfig from './DingtalkConfig.vue'
import { NotifierTypeMap } from '@/common/types'

export const notifierConfigComponents = {
  [NotifierTypeMap.wechatWorkAPPBot]: WechatWorkConfig,
  [NotifierTypeMap.wechatWorkWebhookBot]: WechatWorkWebhookConfig,
  [NotifierTypeMap.telegramAppBot]: TelegramConfig,
  [NotifierTypeMap.dingTalkAppBot]: DingtalkConfig,
} as const

export type NotifierConfigType = keyof typeof notifierConfigComponents
