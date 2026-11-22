from bitpack.batch.parser.BatchParser import BatchParser
from bitpack.batch.runtime.BatchEngine import BatchEngine


class BatchCommand:

    def execute(self, args):

        if not args:

            print(
                "Usage: bp batch <file.bpbatch>"
            )

            return

        file = args[0]

        try:

            with open(
                file,
                "r",
                encoding="utf-8"
            ) as f:

                source = f.read()

            parser = BatchParser()

            plan = parser.parse(
                source
            )

            engine = BatchEngine()

            engine.execute(
                plan
            )

        except Exception as error:

            print(
                f"Batch failed: {error}"
            )