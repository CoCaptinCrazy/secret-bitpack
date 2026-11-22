from bitpack.core.config.Config import Config

from bitpack.core.logging.Logger import Logger

from bitpack.core.state.StateManager import StateManager

from bitpack.core.events.EventBus import EventBus

from bitpack.storage.StorageManager import StorageManager

from bitpack.runtime.manager.RuntimeManager import RuntimeManager

from bitpack.registry.client.RegistryClient import RegistryClient


class Bitpack:

    _instance = None

    def __new__(cls):

        if cls._instance is None:

            cls._instance = super().__new__(cls)

        return cls._instance

    def __init__(self):

        if hasattr(self, "_initialized"):
            return

        self._initialized = True

        self.config = Config()

        self.logger = Logger("Bitpack")

        self.state = StateManager()

        self.events = EventBus()

        self.storage = StorageManager()

        self.runtime = RuntimeManager()

        self.registry = RegistryClient()

    def initialize(self):

        self.storage.initialize()

        self.config.load()

        print("      ")
        print("========== Bitpack v0.1.17 ==========")
        print("      ")


    def shutdown(self):

        self.logger.info(
            "Bitpack shutting down"
        )