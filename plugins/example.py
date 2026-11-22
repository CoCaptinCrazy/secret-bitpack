import os
import sys


def main():
    if os.environ.get("BITPACK_HOST") != "go":
        raise SystemExit("This extension must be started by Bitpack Go.")

    message = " ".join(sys.argv[1:]) or "hello from the optional Python extension"
    print(message)


if __name__ == "__main__":
    main()
