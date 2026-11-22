class HelpGenerator:

    def show(self):
        print("""
Bitpack

Usage:
  bp forge <path>
  bp ignite <blueprint>
  bp modules
  bp blueprints
  bp stop <module>
  bp rm <module>
  bp assemble <file>

Examples:
  bp forge .
  bp ignite hello-world
  bp modules
""")

    def unknown(self, command):
        print(f"Unknown command: {command}")
        self.show()