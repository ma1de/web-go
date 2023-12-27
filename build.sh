#!/bin/sh 
choice=$(printf "build\nrun" | fzf)

if [ $choice = "build" ]
then 
  printf "Building...\n"
  go build .
fi 

if [ $choice = "run" ]
then 
  printf "Running...\n"
  go run .
fi 

echo "Done!"
