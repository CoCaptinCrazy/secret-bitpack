from bitpack.runtime.manager.RuntimeSupervisor import RuntimeSupervisor
from bitpack.runtime.manager.RuntimeScheduler import RuntimeScheduler


class RuntimeManager:

    def __init__(self):
        self.supervisor = RuntimeSupervisor()
        self.scheduler = RuntimeScheduler()

    def start(self, module):
        self.supervisor.register(module)
        self.supervisor.start(module)

    def stop(self, module_id):
        self.supervisor.stop(module_id)

    def restart(self, module_id):
        self.supervisor.restart(module_id)

    def status(self, module_id):
        return self.supervisor.status(module_id)

    def modules(self):
        return self.supervisor.modules()