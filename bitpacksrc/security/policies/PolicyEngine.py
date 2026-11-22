class PolicyEngine:

    def evaluate(
        self,
        policy,
        action
    ):

        return (
            action
            in policy.rules
        )