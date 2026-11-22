from bitpack.runtime.modules.Module import Module


class ModuleCloner:

    def clone(self, module):

        return Module(
            blueprint=module.blueprint,
            name=f"{module.name}-clone"
        )