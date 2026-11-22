from bitpack.core.Bitpack import Bitpack


class ConfigCommand:

    def execute(self, args):

        bp = Bitpack()

        if not args:

            print(
                bp.config.user_config.values
            )

            return

        if args[0] == "get":

            print(
                bp.config.user_config.get(
                    args[1]
                )
            )

            return

        if args[0] == "set":

            bp.config.user_config.set(
                args[1],
                args[2]
            )

            print(
                f"Updated {args[1]}"
            )