class BPackValidator:

    def validate(
        self,
        graph
    ):

        if not graph.children:
            raise ValueError(
                "Assembly is empty."
            )

        return True