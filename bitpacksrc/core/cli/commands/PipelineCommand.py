from bitpack.batch.BatchRunner import BatchRunner
from bitpack.batch.parser.BatchParser import BatchParser


class PipelineCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bitpack pipeline <file>"
            )
            return

        parser = BatchParser()

        pipeline = parser.parse_file(
            args[0]
        )

        runner = BatchRunner()

        runner.run(
            pipeline
        )