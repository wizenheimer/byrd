#!/bin/bash

# Function to process a single .go file
process_file() {
    local file="$1"

    # Check if first line matches the pattern "// filename"
    first_line=$(head -n 1 "$file")
    if [[ $first_line == "// $file" ]]; then
        # Create a temporary file without the first line
        temp_file=$(mktemp)
        tail -n +2 "$file" > "$temp_file"

        # Replace original file with new content
        mv "$temp_file" "$file"
        echo "Processed: $file"
    else
        echo "Skipped: $file (doesn't match expected pattern)"
    fi
}

# Find all .go files recursively and process them
find . -type f -name "*.go" | while read -r file; do
    process_file "$file"
done
