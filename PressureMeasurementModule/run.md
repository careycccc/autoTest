# 终端 1：启动 Locust UI（别点 Start 按钮！）
locust -f locustfile.py --master --master-bind-host=0.0.0.0 --web-host=0.0.0.0 --skip-cpu-warning

# 终端 2：启动 Boomer（自动拉到 2w 用户）
cd pressure-test
go run main.go --spawn-count=20000 --spawn-rate=1000 --master-host=127.0.0.1 --master-port=5557