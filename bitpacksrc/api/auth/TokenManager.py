import secrets


class TokenManager:

    def issue(self):

        return secrets.token_hex(
            32
        )