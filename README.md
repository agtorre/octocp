# Octocp
Experimenting with file system tools and parallel code

## This is early days
Right now there's a simple dispatcher that can recursively walk a tree and stat everything it finds. I'd like to eventually introduce a plugin system to do any operation on the files you might like (chown, chmod, rm, cp, whatever)

## Building
go build

## Running recursively
./octocp -r src destination