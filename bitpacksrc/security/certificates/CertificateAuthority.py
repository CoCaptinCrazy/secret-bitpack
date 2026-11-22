from bitpack.security.certificates.Certificate import Certificate


class CertificateAuthority:

    def issue(
        self,
        owner
    ):
        return Certificate(owner)