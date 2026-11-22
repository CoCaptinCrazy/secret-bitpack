from bitpack.runtime.execution.ProcessManager import ProcessManager


class Executor:

    def __init__(self):
        self.processes = ProcessManager()

    def execute(
        self,
        module,
        command,
        environment
    ):

        return self.processes.start(
            module.id,
            command,
            environment.export()
        )