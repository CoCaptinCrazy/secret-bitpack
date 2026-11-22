class Firewall:

    def __init__(
        self,
        policies
    ):
        self.policies = policies

    def allowed(
        self,
        source,
        destination
    ):

        for rule in self.policies.all():

            if (
                rule.source == source
                and
                rule.destination == destination
            ):
                return (
                    rule.action
                    == "allow"
                )

        return False