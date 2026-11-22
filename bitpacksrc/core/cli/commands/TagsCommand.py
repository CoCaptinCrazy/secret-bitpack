from bitpack.registry.search.RegistryTags import RegistryTags


class TagsCommand:

    def execute(self, args):

        if not args:

            print(
                "Usage: bitpack tags "
                "<blueprint>"
            )

            return

        tags = RegistryTags()

        result = tags.get(
            args[0]
        )

        if not result:

            print(
                "No tags."
            )

            return

        for tag in result:

            print(tag)