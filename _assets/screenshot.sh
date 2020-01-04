#!/bin/bash

function take_shot () {
  PID=$1
  echo "PID is $PID"
  FNAME=$2

  # Wait window
  sleep 3
  WID=$(xdotool search --pid $PID)
  echo "WID is $WID"

  xwininfo -id $WID
  xdotool search --pid $PID windowactivate
  sleep 1

  echo "Taking screenshot"
  import -window $WID "$FNAME.png"
  sleep 1
  kill $PID
  sleep 1

  # Remove scroll bar
  echo "Removing scroll bar from image"
  mogrify -gravity East -chop 80x0 "$FNAME.png"
  # Trim
  echo "Trimming image"
  convert "$FNAME.png" -trim +repage "$FNAME.png"
}

pushd ../bin

while IFS= read -r -u9 cmd; do
  tmpfile=../_assets/$(mktemp -u scrshot.XXXXX)
  echo "% " > $tmpfile.txt
  echo "% $cmd" >> $tmpfile.txt
  # Run command
  eval "$cmd" >> $tmpfile.txt
  echo "% " >> $tmpfile.txt
  echo -e "% " >> $tmpfile.txt
  konsole --notransparency --noclose --hide-tabbar --hide-menubar --separate -e cat $tmpfile.txt &
  take_shot $! $tmpfile
  sleep 1
  echo "deleting $tmpfile.txt"
  rm "$tmpfile.txt"
done 9< "../_assets/screenshots.txt"

