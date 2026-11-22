class Tokenizer:

    SYMBOLS = [
        "<|",
        "|>",
    ]

    def tokenize(self, source):

        tokens = []

        current = ""

        i = 0

        while i < len(source):

            if source[i:i+2] == "<|":
                if current:
                    tokens.append(current)
                    current = ""

                tokens.append("<|")
                i += 2
                continue

            if source[i:i+2] == "|>":
                if current:
                    tokens.append(current)
                    current = ""

                tokens.append("|>")
                i += 2
                continue

            if source[i].isspace():
                if current:
                    tokens.append(current)
                    current = ""
            else:
                current += source[i]

            i += 1

        if current:
            tokens.append(current)

        return tokens