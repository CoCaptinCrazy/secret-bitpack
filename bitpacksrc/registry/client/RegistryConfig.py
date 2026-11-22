class RegistryConfig:

    DEFAULT_REGISTRY = (
        "https://registry.bitpack.io"
    )

    def __init__(self):
        self.registry = (
            self.DEFAULT_REGISTRY
        )

    def set_registry(
        self,
        registry
    ):
        self.registry = registry

    def get_registry(self):
        return self.registry