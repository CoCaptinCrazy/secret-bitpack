class Validator:

    REQUIRED = [
        "BLUEPRINT",
        "USING",
        "IGNITE"
    ]

    def validate(self, ast):

        found = {
            child.type
            for child in ast.children
        }

        missing = []

        for item in self.REQUIRED:
            if item not in found:
                missing.append(item)

        if missing:
            raise ValueError(
                f"Missing sections: {missing}"
            )

        return True