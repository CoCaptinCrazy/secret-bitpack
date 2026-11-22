from bitpack.registry.push.UploadWorker import UploadWorker


class PushManager:

    def __init__(self):
        self.worker = UploadWorker()

    def publish(
        self,
        blueprint
    ):
        return self.worker.upload(
            blueprint
        )