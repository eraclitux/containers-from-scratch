# Build

On a GNU/Linux systsem:
```
$ go build -o smoker
```
# Usage

From shell of a VM container download Alpine Linux "image":

```
./smoker pull
```
launch main container to run a shell limiting its memory consumption to 10MB:
```
./smoker run -m 10000000 /bin/sh
```
remove image from disk:
```
./smoker rmi
```

# Note
For safety only run this code inside a test VM.


