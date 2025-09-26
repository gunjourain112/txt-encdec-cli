#!/bin/bash
sleep "$1"
current=$(wl-paste 2>/dev/null || xclip -selection clipboard -o 2>/dev/null || xsel --clipboard --output 2>/dev/null)
if [ "$current" = "$2" ]; then
    echo "" | wl-copy 2>/dev/null || echo "" | xclip -selection clipboard 2>/dev/null || echo "" | xsel --clipboard --input 2>/dev/null
fi
