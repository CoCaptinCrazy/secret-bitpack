class IgnoreRule:

    def __init__(
        self,
        pattern,
        negate=False
    ):
        self.pattern = pattern
        self.negate = negate