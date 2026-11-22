import shutil
import time
import os


class BackupCommand:

    def execute(self, args):

        backup_dir = (
            f"data/backups/"
            f"backup-{int(time.time())}"
        )

        os.makedirs(
            backup_dir,
            exist_ok=True
        )

        shutil.copytree(
            "data/blueprints",
            f"{backup_dir}/blueprints",
            dirs_exist_ok=True
        )

        print(
            f"Backup created: "
            f"{backup_dir}"
        )