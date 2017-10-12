# Usage

From shell of a VM container download Alpine Linux "image":

```
$ ./smoker pull
```
launch main container to run a shell limiting its memory consumption to 10MB:
```
$ ./smoker run -m 10000000 /bin/sh
```
attach to main container from another shell:
```
$ ./smoker exec /bin/sh
```

# Note
For safety only run the script inside a test VM.

