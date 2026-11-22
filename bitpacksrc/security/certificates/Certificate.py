import uuid


class Certificate:

    def __init__(self, owner):

        self.id = str(uuid.uuid4())

        self.owner = owner

        self.valid = True