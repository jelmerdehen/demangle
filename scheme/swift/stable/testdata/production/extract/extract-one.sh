#!/bin/sh
# extract-one.sh — extract Swift mangled symbols from one SDK binary/tbd
# Usage: extract-one.sh <framework-name> <iphoneos-sdk-root> <swift-demangle>
#
# Tries candidates in order:
#   <sdk>/System/Library/Frameworks/<FW>.framework/<FW>        (dylib/Mach-O)
#   <sdk>/System/Library/Frameworks/<FW>.framework/<FW>.tbd    (stub)
#   <sdk>/usr/lib/swift/lib<FW>.tbd                            (core stub)
#
# Outputs one line per Swift symbol:
#   <mangled> ---> <demangled>

set -e

FW="$1"
IPHONEOS="$2"
SD="$3"

if [ -z "$FW" ] || [ -z "$IPHONEOS" ] || [ -z "$SD" ]; then
    echo "usage: extract-one.sh <framework> <iphoneos-sdk-root> <swift-demangle>" >&2
    exit 1
fi

# Emit sorted unique mangled symbols from a file.
# .tbd files: grep for _$s tokens.
# Mach-O/dylib: nm -gU then filter.
emit_symbols() {
    f="$1"
    case "$f" in
        *.tbd)
            grep -o '_\$s[A-Za-z0-9_$]*' "$f" 2>/dev/null | sort -u
            ;;
        *)
            xcrun nm -gU "$f" 2>/dev/null \
                | awk '/ [TtWwDd] /{ print $3 }' \
                | grep '^_\$s' \
                | sort -u
            ;;
    esac
}

found=
for candidate in \
    "$IPHONEOS/System/Library/Frameworks/$FW.framework/$FW" \
    "$IPHONEOS/System/Library/Frameworks/$FW.framework/$FW.tbd" \
    "$IPHONEOS/usr/lib/swift/lib$FW.tbd"; do
    if [ -f "$candidate" ]; then
        found="$candidate"
        break
    fi
done

if [ -z "$found" ]; then
    echo "WARNING: no binary found for $FW under $IPHONEOS, skipping" >&2
    exit 0
fi

emit_symbols "$found" | while IFS= read -r sym; do
    expected=$("$SD" "$sym" 2>/dev/null)
    printf '%s ---> %s\n' "$sym" "$expected"
done
