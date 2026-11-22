class DNS:

    def __init__(self):
        self.records = {}

    def register(
        self,
        name,
        address
    ):
        self.records[name] = address

    def unregister(
        self,
        name
    ):
        self.records.pop(
            name,
            None
        )

    def resolve(
        self,
        name
    ):
        return self.records.get(name)