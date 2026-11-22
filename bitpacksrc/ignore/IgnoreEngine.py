import os

from bitpack.ignore.IgnoreParser import IgnoreParser
from bitpack.ignore.IgnoreMatcher import IgnoreMatcher


class IgnoreEngine:

    def __init__(self):

        self.parser = IgnoreParser()
        self.matcher = IgnoreMatcher()

    def load(self):

        if not os.path.exists(".bpackignore"):
            return []

        return self.parser.parse(
            ".bpackignore"
        )

    def should_ignore(
        self,
        path
    ):

        rules = self.load()

        ignored = False

        for rule in rules:

            if self.matcher.matches(
                path,
                rule
            ):

                ignored = not rule.negate

        return ignored

    def collect(self, root):
        files = {}

        for dirpath, _, filenames in os.walk(root):
            for filename in filenames:
                path = os.path.join(dirpath, filename)

                rel_path = os.path.relpath(
                    path,
                    root
                )

                if self.should_ignore(rel_path):
                    continue

                try:
                    with open(
                        path,
                        "r",
                        encoding="utf-8"
                    ) as file:
                        files[rel_path] = file.read()

                except Exception:
                    continue

        return files