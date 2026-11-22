class APIServer:

    def __init__(self):
        self.routes = []

    def start(
        self,
        host="0.0.0.0",
        port=8080
    ):
        print(
            f"Bitpack API listening on {host}:{port}"
        )