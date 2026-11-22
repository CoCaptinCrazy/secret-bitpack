class ModuleInspector:

    def inspect(self, module):

        return {
            "id":
                module.id,

            "name":
                module.name,

            "blueprint":
                module.blueprint,

            "status":
                module.state.status
        }