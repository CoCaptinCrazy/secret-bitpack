class Scaling:

    def scale(
        self,
        service,
        replicas
    ):

        return {
            "service":
                service,

            "replicas":
                replicas
        }