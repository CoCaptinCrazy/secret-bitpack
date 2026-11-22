import zipfile


class Compressor:

    def compress(
        self,
        source,
        output
    ):
        with zipfile.ZipFile(
            output,
            "w",
            zipfile.ZIP_DEFLATED
        ) as archive:

            archive.write(
                source
            )