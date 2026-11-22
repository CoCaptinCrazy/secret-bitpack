import os
from bitpack.core.logging.LogFormatter import LogFormatter


class LogManager:
    LOG_DIR = "./data/logs"
    LOG_FILE = LOG_DIR + "/bitpack.log"

    def __init__(self):
        os.makedirs(
            self.LOG_DIR,
            exist_ok=True
        )
        self.formatter = LogFormatter()

    def write(self, level, source, message):
        entry = self.formatter.format(level, source, message)

        print(entry)

        with open(self.LOG_FILE, "a", encoding="utf-8") as log:
            log.write(entry + "\n")