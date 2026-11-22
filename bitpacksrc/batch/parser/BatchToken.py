class BatchToken:

    def __init__(
        self,
        token_type,
        value
    ):

        self.type = token_type

        self.value = value

    def __repr__(self):

        return (
            f"BatchToken("
            f"type={self.type}, "
            f"value={self.value}"
            f")"
        )