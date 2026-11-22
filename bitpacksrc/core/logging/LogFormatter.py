from datetime import datetime


class LogFormatter:
    def format(self, level, source, message):
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

        return (
            f"[{timestamp}] "
            f"[{level}] "
            f"[{source}] "
            f"{message}"
        )