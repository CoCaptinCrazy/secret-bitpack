from bitpack.registry.client.RegistryAuth import RegistryAuth


class LoginCommand:

    def execute(self, args):

        auth = RegistryAuth()

        username = (
            args[0]
            if args
            else "local"
        )

        auth.login(
            username,
            "token"
        )

        print(
            f"Logged in as "
            f"{username}"
        )