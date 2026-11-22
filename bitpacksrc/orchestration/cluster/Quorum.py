class Quorum:

    def healthy(
        self,
        total,
        online
    ):
        return online > (
            total // 2
        )