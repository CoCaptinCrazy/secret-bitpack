from bitpack.core.state.StateStore import StateStore


class StateManager:
    def __init__(self):
        self.store = StateStore()
        self.state = self.store.load()

    def get(self, key, default=None):
        return self.state.get(key, default)

    def set(self, key, value):
        self.state[key] = value
        self.store.save(self.state)

    def remove(self, key):
        if key in self.state:
            del self.state[key]
            self.store.save(self.state)