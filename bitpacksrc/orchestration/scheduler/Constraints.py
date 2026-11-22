class Constraints:

    def __init__(self):
        self.rules = {}

    def set(
        self,
        key,
        value
    ):
        self.rules[key] = value

    def get(
        self,
        key
    ):
        return self.rules.get(key)