#!/bin/bash

echo "Go Get Sleep Pushing into Github Now"

read -p "Commit message:" uservar
git commit -a -m $uservar
git push origin Matthew

echo "Choo Choo"
sleep 2

sl
