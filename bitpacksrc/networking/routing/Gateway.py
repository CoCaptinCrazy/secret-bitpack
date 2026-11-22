class Gateway:

    def __init__(
        self,
        router
    ):
        self.router = router

    def forward(
        self,
        destination
    ):

        return self.router.route(
            destination
        )