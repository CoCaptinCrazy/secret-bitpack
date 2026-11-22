import hashlib

from bitpack.security.signing.Signature import Signature


class Signer:

    def sign(
        self,
        content
    ):

        digest = hashlib.sha256(
            content.encode()
        ).hexdigest()

        return Signature(
            digest
        )