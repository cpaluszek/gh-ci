#!/bin/bash

# Function to choose contrasting text color
contrast_color() {
    local bg=$1
    if (( bg < 232 && (bg % 36) / 6 + (bg % 6) < 4 )) || (( bg >= 232 && bg < 244 )); then
        echo "7" # white text for dark backgrounds
    else
        echo "0" # black text for light backgrounds
    fi
}

# Print 16 basic colors with indices
echo "16 Basic ANSI Colors:"
for i in {0..15}; do
    text_color=$(contrast_color $i)
    printf "\e[48;5;${i}m\e[38;5;${text_color}m  %3d  \e[0m" $i
    if (( (i+1) % 8 == 0 )); then echo; fi
done
echo

# Print 6x6x6 color cube (216 colors) with indices
echo "6x6x6 Color Cube (216 colors):"
for i in {16..231}; do
    text_color=$(contrast_color $i)
    printf "\e[48;5;${i}m\e[38;5;${text_color}m  %3d  \e[0m" $i
    if (( (i-15) % 6 == 0 )); then echo -n " "; fi
    if (( (i-15) % 36 == 0 )); then echo; fi
done
echo

# Print grayscale colors (24 colors) with indices
echo "Grayscale (24 colors):"
for i in {232..255}; do
    text_color=$(contrast_color $i)
    printf "\e[48;5;${i}m\e[38;5;${text_color}m  %3d  \e[0m" $i
    if (( (i+1) % 6 == 0 )); then echo; fi
done
echo