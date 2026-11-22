class Policies:

    def __init__(self):
        self.rules = []

    def add(
        self,
        rule
    ):
        self.rules.append(
            rule
        )

    def all(self):
        return self.rules