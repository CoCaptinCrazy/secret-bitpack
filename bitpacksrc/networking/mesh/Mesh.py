import uuid


class Mesh:

    def __init__(self, name):
        self.id = str(uuid.uuid4())

        self.name = name

        self.modules = {}

    def join(self, module):

        self.modules[module.name] = {
            "id": module.id,
            "module": module
        }

    def leave(self, module_name):

        self.modules.pop(
            module_name,
            None
        )