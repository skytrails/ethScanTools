#!/bin/bash
echo "kill eth-scan service"
killall eth-scan # kill go-admin service
nohup ./eth-scan server -c=config/settings.dev.yml >> access.log 2>&1 & #后台启动服务将日志写入access.log文件
echo "run eth-scan success"
ps -aux | grep eth-scan
