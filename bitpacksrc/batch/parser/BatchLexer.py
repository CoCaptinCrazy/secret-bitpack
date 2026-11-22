from bitpack.batch.parser.BatchToken import BatchToken
from bitpack.batch.parser.BatchTokenType import BatchTokenType


class BatchLexer:

    KEYWORDS = {
        "batch",
        "author",
        "parallel",
        "failfast",
        "requires",
        "tags",
        "forge",
        "verify",
        "publish",
        "assemble",
        "ignite"
    }

    def lex(self, lines):

        tokens = []

        for line in lines:

            line = line.strip()

            if not line:
                continue

            if line.endswith(":"):

                key = line[:-1]

                if key in self.KEYWORDS:

                    tokens.append(
                        BatchToken(
                            BatchTokenType.KEYWORD,
                            key
                        )
                    )

                else:

                    tokens.append(
                        BatchToken(
                            BatchTokenType.VALUE,
                            line
                        )
                    )

                continue

            if line.startswith("batch"):

                tokens.append(
                    BatchToken(
                        BatchTokenType.KEYWORD,
                        "batch"
                    )
                )

                tokens.append(
                    BatchToken(
                        BatchTokenType.VALUE,
                        line
                    )
                )

                continue

            tokens.append(
                BatchToken(
                    BatchTokenType.VALUE,
                    line
                )
            )

        tokens.append(
            BatchToken(
                BatchTokenType.EOF,
                ""
            )
        )

        return tokens