from bitpack.assembly.parser.BPackParser import BPackParser
from bitpack.assembly.deployment.Deployment import Deployment


class AssembleCommand:

    def execute(self, args):

        if not args:
            print("Usage: bitpack assemble <file.bpack>")
            return

        file = args[0]

        try:

            parser = BPackParser()

            graph = parser.parse(open(file).read())

            deployment = Deployment()

            deployment.deploy(graph)

            print(f"Assembly deployed: {file}")

        except Exception as e:

            print(f"Assembly failed: {e}")