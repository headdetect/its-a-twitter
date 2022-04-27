#!/usr/bin/env bash

echo -n "This will reset this branch to master and deploy. Continue? (y/N) "
read answer

if [ "$answer" != "${answer#[Nn]}" ] || ["$answer" != ""]; then
  exit 0
fi

# reset state to latest
git fetch origin
git reset --hard origin/master

# build the app
npm run build

# Clear /docs
rm -rf ../docs
mkdir -p ../docs

# move to /docs
mv build/* ../docs

cd ..

# commit & push changes
git add .
git commit -am "gh-pages build"
git push origin build