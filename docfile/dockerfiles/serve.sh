#!/bin/bash

# 创建日志目录
mkdir -p $LOG_PATH

# 修改目录权限
chmod 777 -R $LOG_PATH

# 运行程序
$APP_PATH/main -c $APP_PATH/datariver.conf.json -g medication
