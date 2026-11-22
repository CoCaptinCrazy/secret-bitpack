class RegistrySearch:

    def __init__(self):

        self.index = [
            "postgres",
            "redis",
            "nginx",
            "ubuntu",
            "website"
        ]

    def find(
        self,
        query
    ):

        query = query.lower()

        return [
            item
            for item
            in self.index
            if query in item.lower()
        ]