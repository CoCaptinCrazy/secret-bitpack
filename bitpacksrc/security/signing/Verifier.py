import hashlib


class Verifier:

    def verify(
        self,
        content,
        signature
    ):

        digest = hashlib.sha256(
            content.encode()
        ).hexdigest()

        return (
            digest ==
            signature.digest
        )