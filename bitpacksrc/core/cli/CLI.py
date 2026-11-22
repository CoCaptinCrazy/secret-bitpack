from bitpack.core.bootstrap.Bootstrap import Bootstrap
from bitpack.core.cli.ArgumentParser import ArgumentParser
from bitpack.core.cli.CommandRouter import CommandRouter


def main():


    bitpack = Bootstrap().initialize()


    parser = ArgumentParser()

    command, args = parser.parse()


    router = CommandRouter()


    router.route(command, args)


if __name__ == "__main__":
    main()