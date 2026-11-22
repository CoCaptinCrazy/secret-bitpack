class ErrorHandler:

    def handle(
        self,
        error
    ):
        return {
            "error": str(error)
        }