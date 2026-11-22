from bitpack.batch.model.BatchPlan import BatchPlan


class BatchParser:

    def parse(
        self,
        source
    ):

        plan = BatchPlan()

        section = None

        for raw in source.splitlines():

            line = raw.strip()

            if not line:
                continue

            if line.startswith("batch"):

                parts = line.split('"')

                if len(parts) >= 2:
                    plan.name = parts[1]

                version_parts = line.split()

                if len(version_parts) >= 3:
                    plan.version = version_parts[-1]

                continue

            if line == "author:":
                section = "author"
                continue

            if line == "requires:":
                section = "requires"
                continue

            if line == "tags:":
                section = "tags"
                continue

            if line == "forge:":
                section = "forge"
                continue

            if line == "verify:":
                section = "verify"
                continue

            if line == "publish:":
                section = "publish"
                continue

            if line == "parallel:":
                section = "parallel"
                continue

            if line == "failfast:":
                section = "failfast"
                continue

            if section == "author":

                plan.author = line
                section = None

            elif section == "parallel":

                plan.parallel = (
                    line.lower() == "true"
                )
                section = None

            elif section == "failfast":

                plan.fail_fast = (
                    line.lower() == "true"
                )
                section = None

            elif section == "forge":

                plan.forge.append(line)

            elif section == "verify":

                plan.verify.append(line)

            elif section == "publish":

                plan.publish.append(line)

            elif section == "requires":

                plan.requires.append(line)

            elif section == "tags":

                plan.tags.append(line)

        return plan