class CertificateStore:

    def __init__(self):
        self.certificates = {}

    def add(self, cert):
        self.certificates[
            cert.id
        ] = cert

    def get(self, cert_id):
        return self.certificates.get(cert_id)