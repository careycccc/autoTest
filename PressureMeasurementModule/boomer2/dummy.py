# dummy.py  —— 永远不需要改的版本
from locust import HttpUser, task, between

# 这段代码永远不会被执行，只是让 locust --master 能正常启动
class DummyUser(HttpUser):
    wait_time = between(1, 1)
    
    @task
    def dummy_task(self):
        pass