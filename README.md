# qq机器人

基于 [go-cqhttp](https://github.com/Mrs4s/go-cqhttp) 实现的qq机器人，在其上额外增加了一个直接本地处理消息并根据配置自动回复的逻辑, 具体改动可见 *qqrobot* 目录

# 使用说明

## 配置机器人

请参考 [go-cqhttp的文档](https://docs.go-cqhttp.org/)

## 配置自动回复逻辑

复制qqrobot/config.toml到qq_robot.exe所在目录 然后打开qq_robot/config.go和config.toml，按照注释，自行调整配置，并添加各种规则

# 注意
## 腾讯云部署
需要修改dns为
101.6.6.6
223.6.6.6
114.114.114.114

不能使用8.8.8.8，已经被腾讯云污染了
具体可见 https://github.com/Mrs4s/go-cqhttp/issues/1115#issuecomment-949298739

# TODO

- [ ] 看看是否有更多事件值得接入，目前仅处理了群聊、私聊和加群这三个事件
