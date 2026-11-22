class RuntimeSupervisor:

    def __init__(self):
        self._modules = {}

    def register(self, module):

        self._modules[module.id] = {
            "module": module,
            "status": "registered"
        }

    def start(self, module):

        self._modules[module.id]["status"] = (
            "running"
        )

    def stop(self, module_id):

        if module_id in self._modules:
            self._modules[module_id][
                "status"
            ] = "stopped"

    def restart(self, module_id):

        if module_id in self._modules:
            self._modules[module_id][
                "status"
            ] = "running"

    def status(self, module_id):

        module = self._modules.get(
            module_id
        )

        if not module:
            return "unknown"

        return module["status"]

    def modules(self):
        return self._modules