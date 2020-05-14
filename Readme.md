# 说明

## 配置文件说明

`./resource/config.json` 配置文件中主要包括:

- paths: 需要添加Tag的项目目录(全路径)
- branches: 需要添加Tag的分支
- tag: 要打的Tag版本号后缀(如: 1.1.0.5)
- message: Tag备注

## 运行

将配置文件配置好后, 直接运行 `go run main.go`

## 其他

如果出现panic错误, 说明远程分支已经存在对应的Tag, 需要将其删除
