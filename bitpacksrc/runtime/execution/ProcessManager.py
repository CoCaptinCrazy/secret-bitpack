import subprocess


class ProcessManager:

    def __init__(self):
        self.processes = {}

    def start(
        self,
        module_id,
        command,
        env=None
    ):

        process = subprocess.Popen(
            command,
            shell=True,
            env=env
        )

        self.processes[
            module_id
        ] = process

        return process

    def stop(self, module_id):

        process = self.processes.get(
            module_id
        )

        if process:
            process.terminate()

    def status(self, module_id):

        process = self.processes.get(
            module_id
        )

        if not process:
            return "missing"

        if process.poll() is None:
            return "running"

        return "stopped"