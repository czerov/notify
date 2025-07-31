import { defineStore } from 'pinia'
import http from '@/common/axiosConfig'
import { ref } from 'vue'
import { useToast } from 'vue-toast-notification'

const toast = useToast()

export interface IMessageTemplate {
  id: string
  name: string
  content: string
  title: string
  image: string
  url: string
  targets: string
}

export const useTemplatesStore = defineStore('templates', () => {
  // 状态
  const templates = ref<Record<string, IMessageTemplate>>({})
  const loading = ref(false)

  // 获取所有模板
  const fetchTemplates = async () => {
    loading.value = true
    try {
      const response = await http.get('/admin/templates')
      templates.value = response.data
    } catch (error: any) {
      console.error('获取模板列表失败:', error)
      toast.error('获取模板列表失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取单个模板
  const getTemplate = async (templateId: string) => {
    try {
      const response = await http.get(`/admin/templates/${templateId}`)
      return response.data
    } catch (error: any) {
      console.error('获取模板失败:', error)
      toast.error('获取模板失败')
      throw error
    }
  }

  // 创建模板
  const createTemplate = async (templateData: IMessageTemplate) => {
    loading.value = true
    try {
      const response = await http.post('/admin/templates', templateData)
      if (response.code === 0) {
        // 更新本地状态
        templates.value[templateData.id] = templateData

        toast.success(`模板 ${templateData.name} 创建成功`)
      } else {
        toast.error(response.msg || '创建模板失败')
      }
      return response.data
    } catch (error: any) {
      console.error('创建模板失败:', error)
      if (error.response?.status === 409) {
        toast.error('模板ID已存在')
      } else {
        toast.error('创建模板失败')
      }
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新模板
  const updateTemplate = async (templateId: string, templateData: IMessageTemplate) => {
    loading.value = true
    try {
      const response = await http.request({
        method: 'PUT',
        url: `/admin/templates/${templateId}`,
        data: templateData,
      })
      if (response.code === 0) {
        // 更新本地状态
        templates.value[templateId] = templateData
        fetchTemplates()
        toast.success(`模板 ${templateData.name} 更新成功`)
      } else {
        toast.error(response.msg || '更新模板失败')
      }
      return response.data
    } catch (error: any) {
      console.error('更新模板失败:', error)
      toast.error('更新模板失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  // 删除模板
  const deleteTemplate = async (templateId: string) => {
    loading.value = true
    try {
      const response = await http.delete(`/admin/templates/${templateId}`)

      if (response.code === 0) {
        // 从本地状态中移除
        delete templates.value[templateId]
        fetchTemplates()
        toast.success(`模板 ${templateId} 删除成功`)
      } else {
        toast.error(response.msg || '删除模板失败')
      }
    } catch (error: any) {
      toast.error('删除模板失败')
    } finally {
      loading.value = false
    }
  }

  // 获取模板选项列表（用于下拉选择）
  const getTemplateOptions = () => {
    return Object.entries(templates.value).map(([id, template]) => ({
      value: id,
      title: template.name,
      subtitle: `ID: ${id}`,
    }))
  }

  // 预览模板渲染结果
  const previewTemplate = (templateContent: string, data: Record<string, any>) => {
    try {
      // 简单的模板预览，将变量替换为示例数据
      let preview = templateContent

      // 替换常见变量
      const variables: Record<string, any> = {
        title: data.title || '示例标题',
        content: data.content || '示例内容',
        level: data.level || 'info',
        message: data.message || '示例消息',
        timestamp: new Date().toLocaleString('zh-CN'),
        image: data.image || '',
        url: data.url || '',
        ...data,
      }

      Object.entries(variables).forEach(([key, value]) => {
        // 替换 {{.Key}} 格式的变量
        const regex = new RegExp(`{{\\s*\\.${key}\\s*}}`, 'g')
        preview = preview.replace(regex, String(value))

        // 替换带条件的变量 {{if .Key}}content{{end}}
        const ifRegex = new RegExp(`{{\\s*if\\s+\\.${key}\\s*}}([^{]*){{\\s*end\\s*}}`, 'g')
        preview = preview.replace(ifRegex, value ? '$1' : '')
      })

      // 处理 upper 过滤器
      preview = preview.replace(/{{\\s*\\.(\w+)\\s*\|\\s*upper\\s*}}/g, (match, key) => {
        const value = variables[key]
        return value ? String(value).toUpperCase() : match
      })

      return preview
    } catch (error) {
      console.error('模板预览失败:', error)
      return '模板预览失败'
    }
  }

  return {
    // 状态
    templates,
    loading,

    // 方法
    fetchTemplates,
    getTemplate,
    createTemplate,
    updateTemplate,
    deleteTemplate,
    getTemplateOptions,
    previewTemplate,
  }
})
