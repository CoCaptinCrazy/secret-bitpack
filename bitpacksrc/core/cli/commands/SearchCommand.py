from bitpack.registry.search.RegistrySearch import RegistrySearch


class SearchCommand:

    def execute(self, args):

        if not args:

            print(
                "Usage: bitpack search "
                "<query>"
            )

            return

        search = RegistrySearch()

        results = search.find(
            args[0]
        )

        if not results:

            print(
                "No results."
            )

            return

        for item in results:

            print(item)