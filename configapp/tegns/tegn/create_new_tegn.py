import argparse
import os
import re
import sys
from pathlib import Path

DIR_PATH = Path(__file__).parent

CLASSNAME_REGEXP = re.compile(r"^[a-zA-Z0-9]+$")

def main():
    argv = sys.argv[1:]

    classname = input("Enter the name of the class: ")
    if not CLASSNAME_REGEXP.match(classname):
        print(f"Classname '{classname}' does not match the regexp '{CLASSNAME_REGEXP}'")
        sys.exit(1)
    
    filename = classname.lower() + ".go"
    new_file_path = (DIR_PATH / filename)

    if new_file_path.exists():
        print(f"File '{new_file_path}' already exists")
        sys.exit(1)

    template = DIR_PATH / ".go.template"
    content = template.read_text()

    content = content.replace("__MYSTRUCT__", classname)
    new_file_path.write_text(content)

    print(f"Successfully created {new_file_path}")

if __name__ == "__main__":
    main()
