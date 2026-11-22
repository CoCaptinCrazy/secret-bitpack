class RollingUpdate:

    def update(
        self,
        old_version,
        new_version
    ):

        return {
            "from":
                old_version,

            "to":
                new_version
        }