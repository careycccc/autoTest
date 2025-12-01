from locust import User, task, constant

class DummyUser(User):
    wait_time = constant(0)   # Boomer 完全接管

    @task
    def dummy(self):
        pass  # 永远不会执行，只是让 Locust 启动