import time


class Trace:

    def __init__(self):
        self.start = time.time()

    def stop(self):
        return (
            time.time() -
            self.start
        )