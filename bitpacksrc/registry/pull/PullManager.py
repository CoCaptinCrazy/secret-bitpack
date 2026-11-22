from bitpack.registry.pull.DownloadWorker import DownloadWorker


class PullManager:

    def __init__(self):
        self.worker = DownloadWorker()

    def acquire(
        self,
        blueprint
    ):

        return self.worker.download(
            blueprint
        )