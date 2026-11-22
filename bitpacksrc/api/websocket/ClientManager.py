class ClientManager:

    def __init__(self):
        self.clients = []

    def connect(
        self,
        client
    ):
        self.clients.append(client)

    def disconnect(
        self,
        client
    ):
        self.clients.remove(client)