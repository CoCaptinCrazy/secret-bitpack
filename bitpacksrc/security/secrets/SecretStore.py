class SecretStore:

    def __init__(self):
        self.secrets = {}

    def save(
        self,
        secret
    ):
        self.secrets[
            secret.name
        ] = secret