package app

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"text/template"
	"time"

	"notify/internal/config"
	"notify/internal/logger"
	"notify/internal/notifier"
)

var funcMap = template.FuncMap{
	"contains":   strings.Contains,
	"hasSuffix":  strings.HasSuffix,
	"hasPrefix":  strings.HasPrefix,
	"index":      strings.Index,
	"lastIndex":  strings.LastIndex,
	"replace":    strings.Replace,
	"replaceAll": strings.ReplaceAll,
	"split":      strings.Split,
	"join":       strings.Join,
}

// NotificationApp 通知应用
type NotificationApp struct {
	configManager *config.ConfigManager
	notifiers     map[string]notifier.Notifier
}

// NewNotificationApp 创建通知应用实例
func NewNotificationApp(configManager *config.ConfigManager) *NotificationApp {
	app := &NotificationApp{
		configManager: configManager,
		notifiers:     make(map[string]notifier.Notifier),
	}

	// 初始化通知服务
	app.InitNotifiers()

	return app
}

// initNotifiers 初始化通知服务
func (app *NotificationApp) InitNotifiers() {
	// 遍历所有通知服务实例
	cfg := app.configManager.GetConfig()
	for instanceName, instance := range cfg.Notifiers {
		if !instance.Enabled {
			continue
		}

		switch instance.Type {
		case config.WechatWorkAPPBot:
			// 将map[string]interface{}转换为WechatWorkConfig
			if config, err := app.parseWechatWorkConfig(instance.Config); err == nil {
				config.Enabled = instance.Enabled
				wechatNotifier := notifier.NewWechatWorkNotifier(config)
				app.notifiers[instanceName] = wechatNotifier
			}
		case config.WechatWorkWebhookBot:
			// 将map[string]interface{}转换为WechatWorkWebhookConfig
			if config, err := app.parseWechatWorkWebhookConfig(instance.Config); err == nil {
				config.Enabled = instance.Enabled
				wechatWebhookNotifier := notifier.NewWechatWorkWebhookNotifier(config)
				app.notifiers[instanceName] = wechatWebhookNotifier
			}
		case config.TelegramAppBot:
			// 将map[string]interface{}转换为TelegramConfig
			if config, err := app.parseTelegramConfig(instance.Config); err == nil {
				config.Enabled = instance.Enabled
				telegramNotifier := notifier.NewTelegramNotifier(config)
				app.notifiers[instanceName] = telegramNotifier
			}
		case config.DingTalkAppBot:
			// 将map[string]interface{}转换为DingTalkConfig
			if config, err := app.parseDingTalkConfig(instance.Config); err == nil {
				config.Enabled = instance.Enabled
				dingtalkNotifier := notifier.NewDingTalkNotifier(config)
				app.notifiers[instanceName] = dingtalkNotifier
			}
		}
	}
	logger.Debug("notifiers", cfg.Notifiers)
}

// parseWechatWorkConfig 解析企业微信配置
func (app *NotificationApp) parseWechatWorkConfig(configData map[string]interface{}) (config.WechatWorkConfig, error) {
	cfg := config.WechatWorkConfig{}

	if corpID, ok := configData["corp_id"].(string); ok {
		cfg.CorpID = corpID
	}
	if agentID, ok := configData["agent_id"].(string); ok {
		cfg.AgentID = agentID
	}
	if secret, ok := configData["secret"].(string); ok {
		cfg.Secret = secret
	}
	if proxy, ok := configData["proxy"].(string); ok {
		cfg.Proxy = proxy
	}
	if targets, ok := configData["targets"].(string); ok {
		cfg.Targets = targets
	}

	if cfg.CorpID == "" || cfg.AgentID == "" || cfg.Secret == "" {
		return cfg, fmt.Errorf("企业微信配置不完整")
	}

	return cfg, nil
}

// parseTelegramConfig 解析Telegram配置
func (app *NotificationApp) parseTelegramConfig(configData map[string]interface{}) (config.TelegramConfig, error) {
	cfg := config.TelegramConfig{}

	if botToken, ok := configData["bot_token"].(string); ok {
		cfg.BotToken = botToken
	}
	if chatID, ok := configData["chat_id"].(string); ok {
		cfg.ChatID = chatID
	}
	if proxy, ok := configData["proxy"].(string); ok {
		cfg.Proxy = proxy
	}

	if cfg.BotToken == "" {
		return cfg, fmt.Errorf("Telegram配置不完整")
	}

	return cfg, nil
}

// parseWechatWorkWebhookConfig 解析企业微信群机器人配置
func (app *NotificationApp) parseWechatWorkWebhookConfig(configData map[string]interface{}) (config.WechatWorkWebhookConfig, error) {
	cfg := config.WechatWorkWebhookConfig{}

	if key, ok := configData["key"].(string); ok {
		cfg.Key = key
	}
	if proxy, ok := configData["proxy"].(string); ok {
		cfg.Proxy = proxy
	}

	if cfg.Key == "" {
		return cfg, fmt.Errorf("企业微信群机器人配置不完整：缺少 key")
	}

	return cfg, nil
}

// parseDingTalkConfig 解析钉钉配置
func (app *NotificationApp) parseDingTalkConfig(configData map[string]interface{}) (config.DingTalkConfig, error) {
	cfg := config.DingTalkConfig{}

	if accessToken, ok := configData["access_token"].(string); ok {
		cfg.AccessToken = accessToken
	}
	if secret, ok := configData["secret"].(string); ok {
		cfg.Secret = secret
	}
	if proxy, ok := configData["proxy"].(string); ok {
		cfg.Proxy = proxy
	}
	if targets, ok := configData["targets"].(string); ok {
		cfg.Targets = targets
	}

	if cfg.AccessToken == "" {
		return cfg, fmt.Errorf("钉钉配置不完整")
	}

	return cfg, nil
}

// Send 发送通知
func (app *NotificationApp) Send(ctx context.Context, appConfig config.NotificationApp, req *map[string]any) error {
	// 获取通知应用配置
	appConfig, exists := app.configManager.GetConfig().NotificationApps[appConfig.AppID]
	if !exists {
		return fmt.Errorf("通知应用 %s 不存在", appConfig.Name)
	}

	if !appConfig.Enabled {
		return fmt.Errorf("通知应用 %s 未启用", appConfig.Name)
	}

	// 根据TemplateID查找模板内容
	template, err := app.getTemplateContent(appConfig.TemplateID)
	if err != nil {
		return fmt.Errorf("获取模板失败: %w", err)
	}
	title, err := app.renderTemplate(appConfig.TemplateID+"_title", template.Title, req)
	if err != nil {
		return fmt.Errorf("渲染消息模板失败: %w", err)
	}
	// 渲染消息模板
	content, err := app.renderTemplate(appConfig.TemplateID+"_content", template.Content, req)
	if err != nil {
		return fmt.Errorf("渲染消息模板失败: %w", err)
	}
	url, _ := app.renderTemplate(appConfig.TemplateID+"_url", template.URL, req)
	image, _ := app.renderTemplate(appConfig.TemplateID+"_image", template.Image, req)
	if image == "" {
		image = appConfig.DefaultImage
	}
	targetsStr, _ := app.renderTemplate(appConfig.TemplateID+"_targets", template.Targets, req)
	targets := []string{}
	if targetsStr != "" {
		targets = strings.Split(targetsStr, ",")
	}
	// 创建通知消息
	message := &notifier.NotificationMessage{
		Title:     title,
		Content:   content,
		Image:     image,
		URL:       url,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	if len(appConfig.Notifiers) == 0 {
		return fmt.Errorf("通知应用 %s 未配置任何通知服务", appConfig.Name)
	}
	// 发送到配置的通知服务 - 并发发送，最多10个协程
	const maxConcurrentNotifiers = 10

	var (
		errors    []error
		errorsMux sync.Mutex // 保护errors切片的并发安全
		wg        sync.WaitGroup
		semaphore = make(chan struct{}, maxConcurrentNotifiers) // 信号量，限制并发数
	)

	// 遍历所有通知服务，为每个启动一个协程
	for _, notifierName := range appConfig.Notifiers {
		// 提前检查通知服务是否存在和启用
		notifierInstance, exists := app.notifiers[notifierName]
		if !exists {
			errorsMux.Lock()
			errors = append(errors, fmt.Errorf("通知服务 %s 不存在", notifierName))
			errorsMux.Unlock()
			continue
		}

		if !notifierInstance.IsEnabled() {
			errorsMux.Lock()
			errors = append(errors, fmt.Errorf("通知服务 %s 未启用", notifierName))
			errorsMux.Unlock()
			continue
		}

		// 为每个有效的通知服务启动协程
		wg.Add(1)
		go func(name string, notifier notifier.Notifier) {
			defer wg.Done()

			// 获取信号量，控制并发数
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 发送通知
			if err := notifier.Send(ctx, message, targets); err != nil {
				errorsMux.Lock()
				errors = append(errors, fmt.Errorf("通知服务 %s 发送失败: %w", name, err))
				errorsMux.Unlock()
			}
		}(notifierName, notifierInstance)
	}

	// 等待所有协程完成
	wg.Wait()

	// 如果有错误，返回合并的错误信息
	if len(errors) > 0 {
		var errorMsg string
		for i, err := range errors {
			if i > 0 {
				errorMsg += "\n "
			}
			errorMsg += err.Error()
		}
		return fmt.Errorf("发送通知时发生错误: %s", errorMsg)
	}

	return nil
}

// renderTemplate 渲染消息模板
func (app *NotificationApp) renderTemplate(name string, templateStr string, data *map[string]any) (string, error) {
	if templateStr == "" {
		// 如果没有模板，使用默认格式
		return "", fmt.Errorf("模板不能为空")
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, *data); err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}
	txt := buf.String()
	txt = strings.ReplaceAll(txt, "<no value>", "")
	return txt, nil
}

// GetNotificationApps 获取所有通知应用
func (app *NotificationApp) GetNotificationApps() map[string]config.NotificationApp {
	return app.configManager.GetConfig().NotificationApps
}

// GetNotifiers 获取所有通知服务
func (app *NotificationApp) GetNotifiers() map[string]notifier.Notifier {
	return app.notifiers
}

// getTemplateContent 根据模板ID获取模板内容
func (app *NotificationApp) getTemplateContent(templateID string) (*config.MessageTemplate, error) {
	if templateID == "" {
		return nil, fmt.Errorf("模板ID不能为空")
	}

	template, exists := app.configManager.GetConfig().Templates[templateID]
	if !exists {
		return nil, fmt.Errorf("模板ID '%s' 不存在", templateID)
	}

	return &template, nil
}

// ValidateConfig 验证配置
func (app *NotificationApp) ValidateConfig() error {
	// 验证通知服务配置
	for instanceName, instance := range app.configManager.GetConfig().Notifiers {
		if !instance.Enabled {
			continue
		}

		switch instance.Type {
		case config.WechatWorkAPPBot:
			if _, err := app.parseWechatWorkConfig(instance.Config); err != nil {
				return fmt.Errorf("通知服务实例 %s (企业微信应用) 配置错误: %v", instanceName, err)
			}
		case config.WechatWorkWebhookBot:
			if _, err := app.parseWechatWorkWebhookConfig(instance.Config); err != nil {
				return fmt.Errorf("通知服务实例 %s (企业微信群机器人) 配置错误: %v", instanceName, err)
			}
		case config.TelegramAppBot:
			if _, err := app.parseTelegramConfig(instance.Config); err != nil {
				return fmt.Errorf("通知服务实例 %s (Telegram) 配置错误: %v", instanceName, err)
			}
		case config.DingTalkAppBot:
			if _, err := app.parseDingTalkConfig(instance.Config); err != nil {
				return fmt.Errorf("通知服务实例 %s (钉钉) 配置错误: %v", instanceName, err)
			}
		default:
			return fmt.Errorf("通知服务实例 %s 使用了未知的类型: %s", instanceName, instance.Type)
		}
	}

	// 验证通知应用配置
	for name, appConfig := range app.configManager.GetConfig().NotificationApps {
		if !appConfig.Enabled {
			continue
		}

		// 验证应用ID不为空
		if appConfig.AppID == "" {
			return fmt.Errorf("通知应用 %s 的 app_id 不能为空", name)
		}

		// 验证是否配置了通知服务
		// if len(appConfig.Notifiers) == 0 {
		// 	return fmt.Errorf("通知应用 %s 未配置任何通知服务", name)
		// }

		// 验证引用的通知服务实例是否存在
		for _, notifierName := range appConfig.Notifiers {
			if _, exists := app.configManager.GetConfig().Notifiers[notifierName]; !exists {
				return fmt.Errorf("通知应用 %s 引用了不存在的通知服务实例: %s", name, notifierName)
			}
		}

		// 验证应用级别的认证配置
		if appConfig.Auth != nil && appConfig.Auth.Enabled && appConfig.Auth.Token == "" {
			return fmt.Errorf("通知应用 %s 启用了认证但未配置token", name)
		}
	}

	return nil
}
