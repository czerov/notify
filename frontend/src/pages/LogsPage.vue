<template>
  <div>
    <div class="header">
      日志
    </div>
    <div class="log-container">
      <div v-for="(log, index) in logs" :key="index" class="log-line">
        <v-chip class="ma-1 level-chip" size="x-small" :color="levelColor(log.level)" label>{{ log.level }}</v-chip>
        <span class="timestamp">{{ log.time }}</span>
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
    logs.value.unshift(entry)
    // 控制日志数量，避免内存过高
    if (logs.value.length > 500) {
      logs.value.pop()
    }
    // 滚动到底部
    const container = document.querySelector('.log-container')
    if (container) {
      container.scrollTop = container.scrollHeight
    }
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
  overflow-x: auto;
  overflow-y: hidden;
  .scrollbar();
}

.log-line {
  display: flex;
  flex-direction: row;
  gap: 4px;
  align-items: center;
  font-size: 0.9rem;
  line-height: 1.4;
  padding: 4px 0;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  width: fit-content;

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
    min-width: 400px;
  }
}
</style>
