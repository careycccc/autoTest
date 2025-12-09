go run main.go run constant login -r 20/s -d 6s


go run main.go run constant login -r 2/s -d 30s --verbose   # 生成json文件

http://127.0.0.1:3000/?orgId=1&from=now-6h&to=now&timezone=browser  # 查看数据,需要启动启动 grafana-server.exe

需要开服务，把json放入到服务下面  http-server -p 8088 --cors