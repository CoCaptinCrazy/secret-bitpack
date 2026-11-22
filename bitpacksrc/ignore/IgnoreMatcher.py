import fnmatch


class IgnoreMatcher:

    def matches(
        self,
        path,
        rule
    ):

        return fnmatch.fnmatch(
            path,
            rule.pattern
        )