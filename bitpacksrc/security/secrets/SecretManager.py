class SecretManager:

    def __init__(self):
        self.store = SecretStore()

    def resolve(
        self,
        name
    ):
        secret = self.store.secrets.get(name)

        if secret:
            return secret.value

        return None