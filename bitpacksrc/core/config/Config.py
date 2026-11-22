from bitpack.core.config.GlobalConfig import GlobalConfig
from bitpack.core.config.UserConfig import UserConfig
from bitpack.core.config.RuntimeConfig import RuntimeConfig


class Config:
    def __init__(self):
        self.global_config = GlobalConfig()
        self.user_config = UserConfig()
        self.runtime_config = RuntimeConfig()

    def load(self):
        self.global_config.load()
        self.user_config.load()
        self.runtime_config.load()