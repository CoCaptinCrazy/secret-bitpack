from bitpack.ignore.IgnoreRule import IgnoreRule


class IgnoreParser:

    def parse(self, path):

        rules = []

        with open(path) as f:

            for line in f:

                line = line.strip()

                if not line:
                    continue

                if line.startswith("#"):
                    continue

                negate = line.startswith("!")

                if negate:
                    line = line[1:]

                rules.append(
                    IgnoreRule(
                        line,
                        negate
                    )
                )

        return rules