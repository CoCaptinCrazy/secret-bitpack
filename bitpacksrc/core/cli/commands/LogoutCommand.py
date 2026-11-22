from bitpack.registry.client.RegistryAuth import RegistryAuth


class LogoutCommand:

    def execute(self, args):

        auth = RegistryAuth()

        auth.logout()

        print(
            "Logged out."
        )