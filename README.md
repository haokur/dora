# dora 百宝箱

## 初始化

1. go mod init dora
2. 安装 cobra，go get -u github.com/spf13/cobra@latest
3. 使用 air，go install github.com/cosmtrek/air@latest
4. air 初始化，air init，开发环境运行 air，或者 npm start
5. 打包，go get github.com/karalabe/xgo，执行 npm run build:all
6. 第5条打包要docker装镜像，简单点直接只打包当前的环境的，npm run build