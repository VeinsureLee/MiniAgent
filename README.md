# Mini Agent in Golang

这是一个关于agent如何运行起来的实验项目。
记录了agent如何从调用api回答问题开始，
到一步步构建起可以做简单coding任务的cli agent项目，
整体使用Golang进行开发。

模型配置位于 `config/models.yaml`，由 `config.Load` 加载：

- `models` 下的键是配置别名，例如 `deepseek`；`default_model` 和 `/model` 都使用这个别名。
- `model` 是发送给 API 的模型 ID，例如 `deepseek-chat`，对应 Go 字段 `ModelConfig.ModelID`。
- `api_key_env` 是保存 API Key 的环境变量名，对应 `ModelConfig.APIKeyEnvVar`。

`Config.ModelsByAlias` 保存别名与配置的映射；`Agent.ActiveModel` 保存当前模型配置，`Agent.APIClient` 保存对应的 API 客户端。
