from bitpack.security.secrets.SecretManager import SecretManager


class SecretsCommand:

    def execute(self, args):

        manager = SecretManager()

        secrets = (
            manager.store.secrets
        )

        if not secrets:

            print(
                "No secrets."
            )

            return

        for name in secrets:

            print(name)