class PolicyValidator:

    def validate(
        self,
        policy
    ):

        return bool(
            policy.name
        )