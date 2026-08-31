#!/usr/bin/env bash
# Verify official harbor art hashes and copy into assets/.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
ART_SRC="${1:-$ROOT/pirate-wars-art}"

declare -A FILES=(
  ["harbor-midday-1536x1024.png"]="83b588416cd4dc13a0d98b28e75471c1cb1c64becb6d47bd153af70b7b455ed7"
  ["harbor-mask-1536x1024.png"]="e2d79c0c09f6c1c220c6283a5d19edee81ac2d44c94fb9ff04a8c24ba3310bd3"
  ["ship-player-white-8way-160.png"]="13f37014ce960486109ca85c4b0c95c11c000c00360bdc33b3d290f7856b9d65"
  ["ship-npc-red-8way-160.png"]="8e3c5d2e18a3b07735db7789b2c29239a5d03a0e5d0c00a0ecae0b53b791e03d"
)

mkdir -p "$ROOT/assets"
for name in "${!FILES[@]}"; do
  src="$ART_SRC/$name"
  want="${FILES[$name]}"
  if [[ ! -f "$src" ]]; then
    echo "MISSING: $src (want sha256 $want)" >&2
    exit 1
  fi
  got=$(sha256sum "$src" | awk '{print $1}')
  if [[ "$got" != "$want" ]]; then
    echo "HASH MISMATCH: $src got $got want $want" >&2
    exit 1
  fi
  cp "$src" "$ROOT/assets/$name"
  echo "OK $name"
done
echo "Harbor assets installed under assets/"
