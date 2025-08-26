curl https://cursor.com/install -fsS | bash
Next Steps

1. Add ~/.local/bin to your PATH:
   For zsh:
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
   source ~/.zshrc

2. Start using Cursor Agent:
   cursor-agent


###  快速打开聊天窗口的快捷键

- Cmd + K - 打开聊天窗口

- Cmd + L - 聚焦到聊天输入框
- /fix 请修复这个函数的bug
- /test 为这个函数生成单元测试
- /@currentfile 请优化这个文件的性能
- /@filename.go 解释这个文件的作用

请用Go语言生成一个HTTP服务器，包含：
- 端口8080
- /health 健康检查端点
- /api/users GET端点返回用户列表
- 使用gin框架

请解释这段Go代码的并发模式：
/fix 这个函数有竞态条件，请修复：
基于我当前打开的main.go文件，请建议如何添加错误处理
请比较main.go和handler.go中的函数，找出重复代码
/@*.go 在所有Go文件中查找使用过时的API
分析我的项目结构，建议更好的包组织方式

- Cmd/Ctrl + I - 快速插入代码片段

- Cmd/Ctrl + / - 对选中代码进行重构/解释

- 右键代码 → "Ask Cursor" - 针对特定代码提问

自定义指令

#style 保持代码简洁，使用Go最佳实践
#tone 专业但友好
请生成一个配置解析器

调试帮助
为什么这个goroutine会泄漏？

小技巧
使用 @ 符号 引用特定文件或符号

用 ``` 包裹代码 获得更准确的回答

明确指定语言 如 "用Go语言实现..."

要求提供示例 如 "请给出使用示例"