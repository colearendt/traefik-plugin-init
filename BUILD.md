# Build Directions

Right now, builds are manual. You can build the SHA reference easily with something like:
```
just VERSION=$(git show-ref -s --head | head -1) build
```
