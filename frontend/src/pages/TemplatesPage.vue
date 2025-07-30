<template>
  <!-- 页面加载状态 -->
  <LoadingView :loading="pageLoading" text="正在加载模板数据..." />

  <!-- 页面内容 -->
  <div>
    <!-- 页面标题 -->
    <v-row>
      <v-col>
        <h1 class="text-h6">
          <v-icon icon="mdi-file-document-multiple" class="mr-2"></v-icon>
          模板管理
        </h1>
        <p class="text-body-2 text-medium-emphasis mt-1">
          管理消息模板，支持模板复用和实时预览
        </p>
      </v-col>
      <v-col cols="auto">
        <v-btn color="primary" @click="showCreateDialog = true" size="small">
          <v-icon icon="mdi-plus" class="mr-2"></v-icon>
          创建模板
        </v-btn>
      </v-col>
    </v-row>

    <!-- 模板列表 -->
    <v-row v-if="Object.keys(templatesStore.templates).length > 0">
      <v-col v-for="(template, templateId) in templatesStore.templates" :key="templateId" cols="12" md="6" lg="4">
        <v-card class="template-card" elevation="2">
          <v-card-title class="d-flex align-center">
            <v-icon icon="mdi-file-document" color="primary" class="mr-3"></v-icon>
            <div class="flex-grow-1">
              <div class="text-h6">{{ template.name }}</div>
              <div class="text-caption text-medium-emphasis">ID: {{ template.id }}</div>
            </div>
          </v-card-title>

          <v-card-text class="template-item">
            <div class="text-body-2 mb-2">标题:</div>
            <v-card variant="outlined" class="template-preview pa-2" style="max-height: 150px; overflow-y: auto;">
              <pre class="text-caption">{{ template.title }}</pre>
            </v-card>
          </v-card-text class="template-item">
          <v-card-text class="template-item">
            <div class="text-body-2 mb-2">内容:</div>
            <v-card variant="outlined" class="template-preview pa-2" style="max-height: 150px; overflow-y: auto;">
              <pre class="text-caption">{{ template.content }}</pre>
            </v-card>
          </v-card-text class="template-item">
          <v-card-text class="template-item">
            <div class="text-body-2 mb-2">链接:</div>
            <v-card variant="outlined" class="template-preview pa-2" style="max-height: 150px; overflow-y: auto;">
              <pre class="text-caption">{{ template.url }}</pre>
            </v-card>
          </v-card-text>
          <v-card-text>
            <div class="text-body-2 mb-2">图片:</div>
            <v-card variant="outlined" class="template-preview pa-2" style="max-height: 150px; overflow-y: auto;">
              <pre class="text-caption">{{ template.image }}</pre>
            </v-card>
          </v-card-text>

          <v-card-actions>
            <v-btn variant="flat" size="small" @click="editTemplate(template)" :loading="templatesStore.loading">
              <v-icon icon="mdi-pencil" class="mr-1"></v-icon>
              编辑
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn variant="text" size="small" color="error" @click="deleteTemplate(template.id)"
              :loading="templatesStore.loading">
              <v-icon icon="mdi-delete"></v-icon>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- 无数据状态 -->
    <v-row v-else-if="!templatesStore.loading">
      <v-col cols="12" class="text-center py-8">
        <v-icon icon="mdi-file-document-multiple-outline" size="80" class="text-medium-emphasis mb-4"></v-icon>
        <h3 class="text-h6 text-medium-emphasis mb-2">暂无模板</h3>
        <p class="text-body-2 text-medium-emphasis mb-4">点击上方的"创建模板"按钮来创建您的第一个消息模板</p>
        <v-btn color="primary" variant="outlined" @click="showCreateDialog = true">
          <v-icon icon="mdi-plus" class="mr-2"></v-icon>
          创建模板
        </v-btn>
      </v-col>
    </v-row>

    <!-- 加载状态 -->
    <v-row v-else>
      <v-col cols="12" class="text-center py-8">
        <v-progress-circular indeterminate color="primary" size="40"></v-progress-circular>
        <p class="text-body-2 text-medium-emphasis mt-4">加载模板数据中...</p>
      </v-col>
    </v-row>

    <!-- 创建/编辑模板对话框 -->
    <TemplateEditDialog v-model="showCreateDialog" :editing-template="editingTemplate" :loading="templatesStore.loading"
      @save="handleSaveTemplate" @cancel="handleCancelEdit" />

  </div>
</template>

<script setup lang="ts">
import LoadingView from '@/components/LoadingView.vue'
import TemplateEditDialog from '@/components/dialog/TemplateEditDialog.vue'
import { useTemplatesStore, type IMessageTemplate } from '@/store/templates'
import { onMounted, ref } from 'vue'
import { useConfirm } from 'vuetify-use-dialog'

const templatesStore = useTemplatesStore()
const createConfirm = useConfirm()
// 页面加载状态
const pageLoading = ref(true)

// 对话框状态
const showCreateDialog = ref(false)

// 编辑状态
const editingTemplate = ref<IMessageTemplate | null>(null)

// 加载模板
const loadTemplates = async () => {
  await templatesStore.fetchTemplates()
}

// 编辑模板
const editTemplate = (template: IMessageTemplate) => {
  editingTemplate.value = template
  showCreateDialog.value = true
}


// 删除模板
const deleteTemplate = async (templateId: string) => {
  const template = templatesStore.templates[templateId]

  const isConfirmed = await createConfirm({
    title: '确认删除',
    content: `您确定要删除模板 "${template?.name}" 吗？此操作不可撤销。`,
    dialogProps: {
      width: '500px',
    },
    confirmationText: '删除',
    confirmationButtonProps: {
      color: 'error',
      variant: 'outlined'
    },
    cancellationText: '取消'
  })

  if (isConfirmed) {
    try {
      await templatesStore.deleteTemplate(templateId)
    } catch (error) {
      console.error('删除模板失败:', error)
    }
  }
}

// 处理保存模板
const handleSaveTemplate = async (formData: any) => {
  try {
    if (editingTemplate.value) {
      await templatesStore.updateTemplate(editingTemplate.value.id, formData)
    } else {
      await templatesStore.createTemplate(formData)
    }
    handleCancelEdit()
  } catch (error) {
    console.error('保存模板失败:', error)
  }
}

// 处理取消编辑
const handleCancelEdit = () => {
  showCreateDialog.value = false
  editingTemplate.value = null
}

// 页面加载时获取数据
onMounted(async () => {
  try {
    await loadTemplates()
  } catch (error) {
    console.error('页面初始化失败:', error)
  } finally {
    pageLoading.value = false
  }
})
</script>

<style scoped>
pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.template-item {
  padding-block: 8px;
}
</style>
