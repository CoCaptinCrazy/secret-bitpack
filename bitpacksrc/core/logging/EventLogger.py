from bitpack.core.logging.Logger import Logger


class EventLogger:
    def __init__(self):
        self.logger = Logger("Events")

    def blueprint_created(self, name):
        self.logger.info(f"Blueprint created: {name}")

    def module_started(self, module_id):
        self.logger.info(f"Module started: {module_id}")

    def module_stopped(self, module_id):
        self.logger.info(f"Module stopped: {module_id}")

    def snapshot_created(self, snapshot_id):
        self.logger.info(f"Snapshot created: {snapshot_id}")