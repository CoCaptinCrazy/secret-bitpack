from bitpack.core.logging.LogManager import LogManager


class Logger:
    def __init__(self, name="Bitpack"):
        self.name = name
        self.manager = LogManager()

    def info(self, message):
        self.manager.write("INFO", self.name, message)

    def warn(self, message):
        self.manager.write("WARN", self.name, message)

    def error(self, message):
        self.manager.write("ERROR", self.name, message)

    def debug(self, message):
        self.manager.write("DEBUG", self.name, message)