from bitpack.registry.client.RegistryConfig import RegistryConfig
from bitpack.registry.client.RegistryAuth import RegistryAuth


class RegistryClient:

    def __init__(self):

        self.config = RegistryConfig()

        self.auth = RegistryAuth()

    def endpoint(self):

        return self.config.get_registry()