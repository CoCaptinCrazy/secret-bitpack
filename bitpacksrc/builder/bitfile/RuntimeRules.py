class RuntimeRules:

    REQUIRED = [
        "blueprint",
        "using",
        "ignite"
    ]

    VALID_ENVIRONMENTS = [
        "minimal",
        "standard",
        "server",
        "desktop",
        "ai",
        "embedded"
    ]

    @classmethod
    def validate_environment(cls, env):
        return env in cls.VALID_ENVIRONMENTS