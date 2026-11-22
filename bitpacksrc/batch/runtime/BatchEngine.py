from bitpack.core.cli.commands.ForgeCommand import ForgeCommand
from bitpack.core.cli.commands.AssembleCommand import AssembleCommand


class BatchEngine:

    def execute(self, plan):

        print(f"[BATCH] {plan.name}")
        print()

        forge = ForgeCommand()
        assemble = AssembleCommand()

        if plan.forge:

            print("[FORGE]")

            for target in plan.forge:

                print(f"Running forge for {target}")

                # Later you can change ForgeCommand
                # to support filenames.

                forge.execute([target])

            print()

        if plan.verify:

            print("[VERIFY]")

            for target in plan.verify:

                print(f"Verifying {target}")

            print()

        if plan.publish:

            print("[PUBLISH]")

            for target in plan.publish:

                print(f"Publishing {target}")

            print()

        if hasattr(plan, "assemble"):

            print("[ASSEMBLE]")

            for target in plan.assemble:

                print(f"Deploying {target}")

                assemble.execute(
                    [target]
                )

            print()

        if hasattr(plan, "ignite"):

            print("[IGNITE]")

            for target in plan.ignite:

                print(
                    f"Starting {target}"
                )

            print()

        print(
            "[SUCCESS] Batch complete"
        )

        return True