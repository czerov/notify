# 通知管理系统 - 前端

基于Vue 3 + TypeScript + Vuetify的现代化通知管理系统前端应用。

## 🚀 技术栈

- **框架**: Vue 3 (Composition API)
- **语言**: TypeScript
- **UI组件**: Vuetify 3
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **构建工具**: Vite
- **HTTP客户端**: Axios
- **工具库**: 
  - @vueuse/core (Vue工具函数)
  - lodash-es (工具函数)
  - dayjs (日期处理)
  - vue-toast-notification (消息提示)
  - vuetify-use-dialog (对话框)

## 📁 项目结构

```
frontend/
├── public/                     # 静态资源
├── src/
│   ├── assets/                # 资源文件
│   ├── common/                # 通用工具
│   │   ├── types.ts          # TypeScript类型定义
│   │   └── utils.ts          # 工具函数
│   ├── components/           # Vue组件
│   │   ├── AppsPage.vue     # 通知应用管理
│   │   ├── DashboardPage.vue # 系统仪表板
│   │   ├── Layout.vue       # 主布局
│   │   ├── LoginPage.vue    # 登录页面
│   │   ├── NotifiersPage.vue # 通知服务配置
│   │   └── TemplatesPage.vue # 模板管理
│   ├── constants/           # 常量定义
│   ├── plugins/            # 插件配置
│   │   └── vuetify/       # Vuetify配置
│   ├── router/            # 路由配置
│   │   └── index.ts
│   ├── store/            # Pinia状态管理
│   │   ├── apps.ts      # 应用状态
│   │   ├── auth.ts      # 认证状态
│   │   ├── notifiers.ts # 通知服务状态
│   │   └── templates.ts # 模板状态
│   ├── App.vue          # 根组件
│   └── main.ts         # 应用入口
├── index.html          # HTML模板
├── package.json       # 项目配置
├── vite.config.ts    # Vite配置
└── README.md        # 项目文档
```

## 🎯 主要功能

### 🔐 认证系统
- **自动检测**: 系统启动时自动检测是否需要认证
- **Basic认证**: 支持用户名密码认证
- **免认证模式**: 支持无认证环境的警告提示
- **状态持久化**: 认证状态自动保存和恢复

### 📊 系统仪表板
- **统计概览**: 实时显示应用、通知服务、模板数量
- **系统状态**: 监控后端服务、认证、配置状态
- **快速操作**: 一键跳转到各管理页面
- **测试功能**: 直接发送测试通知

### 📱 应用管理
- **完整CRUD**: 创建、读取、更新、删除应用
- **配置丰富**: 支持认证、字段映射、默认图片等
- **模板关联**: 通过模板ID关联消息模板
- **实时测试**: 支持发送测试通知验证配置

### 🔔 通知服务管理
- **多类型支持**: 企业微信、Telegram、钉钉
- **实例管理**: 支持同类型多实例配置
- **配置编辑**: 可视化编辑通知服务参数
- **连接测试**: 实时测试通知服务连通性

### 📝 模板管理
- **可视化编辑**: 支持Go模板语法编辑
- **实时预览**: 编辑时实时预览渲染效果
- **变量说明**: 内置模板变量使用指南
- **复用设计**: 一个模板可被多个应用使用

## 🎨 UI特性

### 现代化设计
- **Material Design**: 基于Vuetify 3的Material设计
- **响应式布局**: 适配桌面、平板、手机屏幕
- **暗色主题**: 支持明暗主题切换（待实现）
- **国际化**: 中文界面，支持多语言扩展

### 交互体验
- **实时反馈**: 操作结果即时Toast提示
- **加载状态**: 异步操作显示加载动画
- **确认对话框**: 危险操作二次确认
- **表单验证**: 实时表单验证和错误提示

### 动画效果
- **卡片悬停**: 卡片hover时的阴影动画
- **页面切换**: 路由切换动画
- **列表动画**: 数据更新时的动画效果

## 🔧 开发指南

### 环境要求
- Node.js >= 16
- pnpm >= 7 (推荐)

### 安装依赖
```bash
cd frontend
pnpm install
```

### 开发运行
```bash
# 启动开发服务器
pnpm dev

# 访问地址: http://localhost:5173
```

### 构建部署
```bash
# 构建生产版本
pnpm build

# 预览构建结果
pnpm preview
```

### 代码检查
```bash
# ESLint检查
pnpm lint

# 自动修复
pnpm lint:fix

# TypeScript类型检查
pnpm check
```

## 🔌 API集成

### 后端接口
- **基础URL**: `http://localhost:8089/api/v1`
- **认证方式**: Basic Authentication (可选)
- **数据格式**: JSON

### 主要接口
```typescript
// 系统状态
GET /health

// 应用管理
GET /admin/apps
POST /admin/apps
PUT /admin/apps/:appId
DELETE /admin/apps/:appId

// 通知服务管理  
GET /admin/notifiers
PUT /admin/notifiers/:notifier
DELETE /admin/notifiers/:notifier
POST /admin/notifiers/:notifier/test

// 模板管理
GET /admin/templates
POST /admin/templates
PUT /admin/templates/:templateId
DELETE /admin/templates/:templateId

// 发送通知
POST /notify/:appId
GET /notify/:appId
```

## 📱 页面说明

### 登录页面 (`/login`)
- 系统状态检测
- 认证表单（条件显示）
- 连接状态指示器
- 渐变背景设计

### 仪表板 (`/`)
- 统计数据卡片
- 系统状态监控
- 快速操作面板
- 系统信息显示

### 应用管理 (`/apps`)
- 应用列表表格
- 创建/编辑对话框
- 高级配置（认证、字段映射）
- 测试通知功能

### 通知服务配置 (`/notifiers`)
- 卡片式实例展示
- 分类型配置表单
- 敏感信息脱敏
- 连通性测试

### 模板管理 (`/templates`)
- 卡片式模板展示
- 分屏编辑预览
- 变量提示面板
- 实时渲染预览

## 🎯 状态管理

### Auth Store
```typescript
// 认证状态
isAuthenticated: boolean
username: string
isAuthRequired: boolean

// 方法
login(username, password)
logout()
checkSystemStatus()
```

### Apps Store
```typescript
// 应用数据
apps: Record<string, NotificationApp>

// 方法
fetchApps()
createApp(app)
updateApp(appId, app)
deleteApp(appId)
sendTestNotification(appId, data)
```

### Notifiers Store
```typescript
// 通知服务数据
notifiers: Record<string, NotifierInstance>

// 方法
fetchNotifiers()
updateNotifier(name, config)
deleteNotifier(name)
testNotifier(name, data)
```

### Templates Store
```typescript
// 模板数据
templates: Record<string, MessageTemplate>

// 方法
fetchTemplates()
createTemplate(template)
updateTemplate(templateId, template)
deleteTemplate(templateId)
previewTemplate(content, data)
```

## 🛠️ 开发规范

### 组件规范
- 使用Composition API
- TypeScript严格模式
- Props和Emits类型定义
- 组件文档注释

### 代码风格
- ESLint + Prettier
- 驼峰命名法
- 2空格缩进
- 单引号字符串

### 提交规范
- feat: 新功能
- fix: 修复bug
- docs: 文档更新
- style: 样式调整
- refactor: 重构
- test: 测试相关

## 🔮 待开发功能

### 即将实现
- [ ] 暗色主题支持
- [ ] 国际化 (i18n)
- [ ] 消息历史记录
- [ ] 高级搜索过滤
- [ ] 配置导入导出

### 长期规划
- [ ] 实时通知状态
- [ ] 可视化统计图表
- [ ] 移动端适配优化
- [ ] 离线PWA支持
- [ ] 插件化架构

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 发起 Pull Request

## 📄 许可证

MIT License

## 🆘 常见问题

### Q: 登录时提示连接失败？
A: 请检查后端服务是否正常运行，以及API地址配置是否正确。

### Q: 某些功能按钮无响应？
A: 请检查浏览器控制台是否有JavaScript错误，并确保认证状态正常。

### Q: 如何修改API地址？
A: 在 `src/store/auth.ts` 中修改 `apiBaseUrl` 配置。

### Q: 如何添加新的通知服务类型？
A: 需要同时修改前端的类型定义和后端的实现逻辑。

---

🎉 **享受现代化的通知管理体验！**
