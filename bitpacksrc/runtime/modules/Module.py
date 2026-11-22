import uuid

from bitpack.runtime.modules.ModuleState import ModuleState
from bitpack.runtime.modules.ModuleMetadata import ModuleMetadata


class Module:

    def __init__(
        self,
        blueprint,
        name=None
    ):
        self.id = str(uuid.uuid4())

        self.name = (
            name or f"module-{self.id[:8]}"
        )

        self.blueprint = blueprint

        self.state = ModuleState()

        self.metadata = ModuleMetadata(
            blueprint
        )

    def to_dict(self):
        return {
            "id": self.id,
            "name": self.name,
            "blueprint": self.blueprint,
            "state": self.state.status
        }