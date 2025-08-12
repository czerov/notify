<template>
  <div>
    <div class="header">
      日志
    </div>
    <div class="log-container" ref="containerRef">
      <div v-for="log in logs" :key="log._k" class="log-line">
        <div class="log-line-left">
          <v-chip class="ma-1 level-chip" size="x-small" :color="levelColor(log.level)" label>{{ log.level }}</v-chip>
          <span class="timestamp">{{ log.time }}</span>
        </div>
        <span class="message">
          <span>{{ log.msg }}</span>
          <span>{{ _.omit(log, ['level', 'time', 'msg']) }}</span>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as _ from 'lodash-es'

interface LogEntry {
  level: string
  time: string
  msg: string
  [key: string]: any
}

const logs = ref<LogEntry[]>([])
let eventSource: EventSource | null = null
const containerRef = ref<HTMLDivElement | null>(null)

// 批量缓冲，避免高频渲染
const buffer: LogEntry[] = []
let flushScheduled = false
let nextId = 1



function scheduleFlush() {
  if (flushScheduled) return
  flushScheduled = true
  requestAnimationFrame(async () => {
    try {
      if (buffer.length === 0) return
      // 批量一次性 push，减少响应触发次数
      logs.value.push(...buffer)
      buffer.length = 0
      // 控制日志数量上限
      if (logs.value.length > 500) {
        const removeCount = logs.value.length - 500
        logs.value.splice(0, removeCount)
      }
    } finally {
      flushScheduled = false
      // 若在 flush 过程中又有新数据，继续安排下一次
      if (buffer.length > 0) scheduleFlush()
    }
  })
}



function levelColor(level: string) {
  switch (level) {
    case 'DEBUG':
      return 'indigo'
    case 'INFO':
      return 'green'
    case 'WARN':
      return 'orange'
    case 'ERROR':
      return 'red'
    default:
      return 'grey'
  }
}

onMounted(() => {
  eventSource = new EventSource('/api/v1/logs/stream')
  eventSource.onmessage = (e) => {
    const entry = JSON.parse(e.data)
    entry.time = new Date(entry.time).toLocaleString()
    // 为每条日志生成稳定 key
    entry._k = nextId++
    buffer.push(entry)
    scheduleFlush()
  }
})

onBeforeUnmount(() => {
  if (eventSource) {
    eventSource.close()
  }
})
</script>

<style scoped lang="less">
@import '@/styles/mix.less';

.header {
  padding: 16px;
  font-size: 1.2rem;
  font-weight: 500;
}

.log-container {
  display: flex;
  flex-direction: column-reverse;
  overflow-y: auto;
  overflow-x: hidden;
  .scrollbar();
}

.log-line {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  font-size: 0.9rem;
  line-height: 1.4;
  padding: 4px 0;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));

  .level-chip {
    width: 60px;
    justify-content: center;
    padding: 0 2px;
    width: 70px;
    flex: 0 0 70px;
  }

  .timestamp {
    color: gray;
    margin: 0 4px;
    flex: 0 0 150px;
  }

  .message {
    word-break: break-all;
    display: flex;
    flex-direction: column;
    gap: 4px;
    flex: 1;
  }
}
</style>
