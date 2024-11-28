#!/bin/bash

last_value=""

for i in {1..100}
do
    echo "Execute the $i-th time"
    current_value=$(gcloc . | tail -n 2 | head -n 1 | awk -F ' ' '{print $5}')
    if [ -n "$last_value" ] && [ "$last_value" != "$current_value" ]; then
        echo "Values are different!"
        echo "Last value: $last_value"
        echo "Current value: $current_value"
        break
    fi
done
