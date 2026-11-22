from bitpack.core.Bitpack import Bitpack


class RuntimeCommand:

    def execute(self, args):

        bp = Bitpack()

        modules = (
            bp.runtime.modules()
        )

        print(
            modules
        )