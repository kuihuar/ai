## 部署环境
```bash
# 部署环境
export GO111MODULE="on"
export GONOPROXY="*.gitlab.com,*.gitee.com,*.100tal.com"
export GONOSUMDB="*.gitlab.com,*.gitee.com,*.100tal.com"
export GOPROXY="https://goproxy.cn,http://go.xesv5.com/proxy/,direct"
export GOPRIVATE="git.100tal.com"
export GOOS=linux
export GOARCH=amd64


git config --global url."ssh://git@git.100tal.com".insteadOf "https://git.100tal.com"

git config --global --get-regexp url

git clean --modcache
git mod tidy
go env

replace git.100tal.com => ssh://git@git.100tal.com
```