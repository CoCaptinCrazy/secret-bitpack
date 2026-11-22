import json
import os

class UserConfig:
    CONFIG_FILE = "./config.json"

    def __init__(self):
        self.values = {
            "registry": "https://registry.bitpack.io",
            "default_runtime": "standard",
            "telemetry": False
        }

    def load(self):
        if os.path.exists(self.CONFIG_FILE):
            with open(self.CONFIG_FILE, "r") as f:
                self.values.update(json.load(f))

    def get(self, key, default=None):
        return self.values.get(key, default)

    def set(self, key, value):
        self.values[key] = value