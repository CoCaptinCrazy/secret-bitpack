from bitpack.core.Bitpack import Bitpack


class Bootstrap:

    def __init__(self):
        self.bitpack = None

    def initialize(self):


        self.bitpack = Bitpack()


        self.bitpack.initialize()


        return self.bitpack