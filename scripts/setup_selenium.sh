#!/usr/bin/env bash
set -euo pipefail

echo "fetching latest chrome info"
meta_data=$(curl 'https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions-with-downloads.json')

echo "downloading latest chrome"
chrome_url=$(echo "$meta_data" | jq -r '.channels.Stable.downloads.chrome[0].url')
curl $chrome_url -o ./tools/chrome-linux64.zip

echo "extracting downloaded chrome"
unzip ./build/chrome-linux64.zip -d ./build

echo "downloading chrome driver"
chrome_driver_url=$(echo "$meta_data" | jq -r '.channels.Stable.downloads.chromedriver[0].url')
curl $chrome_driver_url -o ./tools/chromedriver

echo "installing chrome dependencies"
sudo apt install ca-certificates fonts-liberation \
    libappindicator3-1 libasound2 libatk-bridge2.0-0 libatk1.0-0 libc6 \
    libcairo2 libcups2 libdbus-1-3 libexpat1 libfontconfig1 libgbm1 \
    libgcc1 libglib2.0-0 libgtk-3-0 libnspr4 libnss3 libpango-1.0-0 \
    libpangocairo-1.0-0 libstdc++6 libx11-6 libx11-xcb1 libxcb1 \
    libxcomposite1 libxcursor1 libxdamage1 libxext6 libxfixes3 libxi6 \
    libxrandr2 libxrender1 libxss1 libxtst6 lsb-release wget xdg-utils -y


