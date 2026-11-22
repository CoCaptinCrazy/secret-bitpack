import os


class Environment:

    def __init__(self):
        self.variables = {}

    def set(self, key, value):
        self.variables[key] = str(value)

    def get(self, key):
        return self.variables.get(key)

    def export(self):
        env = os.environ.copy()

        env.update(
            self.variables
        )

        return env