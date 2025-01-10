#!/bin/bash

# Function to process a single .go file
process_file() {
    local file="$1"
    # Create a temporary file
    temp_file=$(mktemp)

    # Add filename as comment and then original content
    echo "// $file" > "$temp_file"
    cat "$file" >> "$temp_file"

    # Replace original file with new content
    mv "$temp_file" "$file"

    echo "Processed: $file"
}

# Find all .go files recursively and process them
find . -type f -name "*.go" | while read -r file; do
    process_file "$file"
done