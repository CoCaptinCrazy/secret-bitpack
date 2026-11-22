class Metric:

    def __init__(
        self,
        name,
        value=0
    ):
        self.name = name
        self.value = value

    def update(
        self,
        value
    ):
        self.value = value