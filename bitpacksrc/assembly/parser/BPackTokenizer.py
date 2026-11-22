class BPackTokenizer:

    SYMBOLS = [
        "{",
        "}"
    ]

    def tokenize(self, source):

        tokens = []
        current = ""

        for char in source:

            if char in self.SYMBOLS:

                if current.strip():
                    tokens.append(current.strip())

                tokens.append(char)
                current = ""

            elif char.isspace():

                if current.strip():
                    tokens.append(current.strip())
                    current = ""

            else:
                current += char

        if current.strip():
            tokens.append(current.strip())

        return tokens