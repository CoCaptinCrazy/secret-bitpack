class RegistryTags:

    def __init__(self):
        self.tags = {}

    def add(
        self,
        blueprint,
        tag
    ):
        self.tags.setdefault(
            blueprint,
            []
        ).append(tag)

    def get(
        self,
        blueprint
    ):
        return self.tags.get(
            blueprint,
            []
        )