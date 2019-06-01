#!/usr/bin/env bash

# Make sure we are on the right directory
cd $(dirname $0)

TAG=$(git describe)
CHECKSUM=$(sha256sum dist/shadowfox_mac_x64 | cut -d' ' -f1)
echo Tag=$TAG
echo Checksum=$CHECKSUM

# Clone repo
git clone https://github.com/SrKomodo/homebrew-tap.git
cd homebrew-tap
git remote rm origin
git remote add origin https://SrKomodo:$1@github.com/SrKomodo/homebrew-tap.git

# Compile template
IFS=''
while read line; do
  REPLACE1=${line//tag/$TAG}
  REPLACE2=${REPLACE1//checksum/$CHECKSUM}
  echo $REPLACE2
done < ../tap_template > Formula/shadowfox-updater.rb

# Push to tap
git add Formula/shadowfox-updater.rb
git commit -m $TAG
git push --set-upstream origin master