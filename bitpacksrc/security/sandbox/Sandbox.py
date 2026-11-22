class Sandbox:

    def __init__(self):
        self.permissions = set()

    def allow(self, permission):
        self.permissions.add(permission)

    def deny(self, permission):
        self.permissions.discard(permission)

    def has(self, permission):
        return permission in self.permissions