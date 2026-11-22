class Resolver:

    def __init__(self, dns):
        self.dns = dns

    def lookup(
        self,
        hostname
    ):
        return self.dns.resolve(
            hostname
        )