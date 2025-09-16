#!/bin/bash

name="EterLine (Andrew)"
year="2025"

LICENSE_TEXT="// Copyright (c) $year $name
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.

"

find . -type f -name "*.go" | while read -r file; do
    if ! grep -q "Copyright (c) $year $name" "$file"; then
        tmpfile=$(mktemp)
        echo "$LICENSE_TEXT" > "$tmpfile"
        cat "$file" >> "$tmpfile"
        mv "$tmpfile" "$file"
        echo "LICENSE added to: $file"
    fi
done