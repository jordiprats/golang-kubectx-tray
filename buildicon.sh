#!/bin/bash

INPUT_FILE=macIcon/main.png

ICONSET_NAME="macIcon/KCT.iconset"

mkdir -p $ICONSET_NAME

sips -z 16 16     $INPUT_FILE --out "${ICONSET_NAME}/icon_16x16.png"
sips -z 32 32     $INPUT_FILE --out "${ICONSET_NAME}/icon_16x16@2x.png"
sips -z 32 32     $INPUT_FILE --out "${ICONSET_NAME}/icon_32x32.png"
sips -z 64 64     $INPUT_FILE --out "${ICONSET_NAME}/icon_32x32@2x.png"
sips -z 128 128   $INPUT_FILE --out "${ICONSET_NAME}/icon_128x128.png"
sips -z 256 256   $INPUT_FILE --out "${ICONSET_NAME}/icon_128x128@2x.png"
sips -z 256 256   $INPUT_FILE --out "${ICONSET_NAME}/icon_256x256.png"
sips -z 512 512   $INPUT_FILE --out "${ICONSET_NAME}/icon_256x256@2x.png"
sips -z 512 512   $INPUT_FILE --out "${ICONSET_NAME}/icon_512x512.png"

iconutil -c icns $ICONSET_NAME

cp macIcon/KCT.icns KubeCtxTray.app/Contents/Resources/KCT.icns