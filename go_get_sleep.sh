#!/bin/bash

echo "Go Get Sleep Pushing into Github Now"
echo "Github Commit Message:"

read userin
echo "<--------------Sleeping Now Commit and Pushing to Github-------------->"

git commit -a -m '$userin'
git push origin Matthew

echo "Choo Choo"
sleep 1

sl




