from bitpack.security.policies.Policy import Policy
from bitpack.security.policies.PolicyValidator import PolicyValidator


class PolicyCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bitpack policy <file.bppolicy>"
            )
            return

        policy = Policy(args[0])

        validator = PolicyValidator()

        valid = validator.validate(
            policy
        )

        print(
            f"Policy valid: {valid}"
        )