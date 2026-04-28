import argparse
from dataclasses import dataclass
from pathlib import Path
import sys
from typing import List, Optional


@dataclass
class GrepGoConfig:
    """Configuration for grepping Go files and writing them as markdown."""
    directories: List[Path]
    output_file: Optional[Path] = None


def get_args(args: List[str]) -> GrepGoConfig:
    """Create config from command line arguments."""
    parser = argparse.ArgumentParser(
        description="Grep all .go files in directories and write them as markdown."
    )
    parser.add_argument(
        "directories",
        nargs="+",
        type=Path,
        help="Directories to search for .go files"
    )
    parser.add_argument(
        "-o", "--output",
        type=Path,
        help="Output file path (if not specified, prints to stdout)"
    )
    
    args = parser.parse_args(args)
    return GrepGoConfig(
        directories=args.directories,
        output_file=args.output
    )


class GoFileGrep:
    """Greps Go files and formats them as markdown."""
    
    def __init__(self, path_list: List[Path]):
        self.path_list = path_list
        
    def find_go_files(self) -> List[Path]:
        """Find all .go files in the specified directories."""
        go_files = []
        for entry in self.path_list:
            if not entry.exists():
                print(f"Warning: Entry {entry} does not exist, skipping...", file=sys.stderr)
                continue
            if entry.is_dir():
                # Recursively find all .go files
                go_files.extend(entry.rglob("*.go"))
                continue

            if entry.is_file():
                go_files.append(entry)

            print(f"Warning: {entry} is neither a directory nor a file, skipping...", file=sys.stderr)

        
        return sorted(go_files)  # Sort for consistent output
    
    @staticmethod
    def _format_file_as_markdown(file_path: Path) -> List[str]:
        lines = []

        try:
            # Read file content
            content = file_path.read_text(encoding='utf-8')
            
            # Add code block with go language specification
            lines.append("```go")
            lines.append(content.rstrip())  # Remove trailing whitespace
            if not lines[-1]:  # Ensure newline after content if empty
                lines.append("")
            lines.append("```\n")

        except Exception as e:
            print(f"*Error reading file: {e.with_traceback()}*\n", file=sys.stderr)
            return []
            
        return [
            f"## {file_path}\n",
            *lines
        ]

    def format_as_markdown(self, file_paths: List[Path]) -> str:
        """Format all files (passed by paths) as markdown with headers and code blocks."""
        
        lines = ["# Code\n"]
        
        for file_path in file_paths:
            lines.extend(
                self._format_file_as_markdown(file_path)
            )
        
        return "\n".join(lines)

    def run(self) -> str:
        """Execute the grep and output operation."""
        go_files = self.find_go_files()
        markdown_content = self.format_as_markdown(go_files)
    
        return markdown_content

def main() -> None:
    """Main entry point."""
    args = get_args(sys.argv[1:])

    grepper = GoFileGrep(args.directories)
    files_content = grepper.run()

    if args.output_file:
        # Write to file
        args.output_file.write_text(files_content, encoding='utf-8')
        print(f"Output written to '{args.output_file}'", file=sys.stderr)
    else:
        print(files_content, file=sys.stdout)

    print(f'Size: {len(files_content)}', file=sys.stderr)


if __name__ == "__main__":
    main()
