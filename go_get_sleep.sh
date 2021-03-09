#!/bin/bash

echo "Go Get Sleep Pushing into Github Now"

read -p "Commit message:" uservar
echo "<--------------Sleeping Now Commit and Pushing to Github-------------->"

git commit -a -m $uservar
git push origin Matthew

echo "Choo Choo"
sleep 1

sl




