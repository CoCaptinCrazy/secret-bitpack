class RegistryAuth:

    def __init__(self):
        self.token = None

    def login(
        self,
        username,
        token
    ):
        self.token = token

    def logout(self):
        self.token = None

    def authenticated(self):
        return self.token is not None