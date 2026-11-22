class ModuleState:

    CREATED = "created"
    STARTING = "starting"
    RUNNING = "running"
    STOPPED = "stopped"
    CRASHED = "crashed"

    def __init__(self):
        self.status = self.CREATED

    def set(self, status):
        self.status = status

    def is_running(self):
        return self.status == self.RUNNING