# containers-from-scratch

Proof of concept POSIX shell script that simulates main functionalities
of Docker.
Intended only to deep knowledge about GNU/Linux containers not for real world
usage.

Different tools for container management (Docker, rkt) on GNU/Linux are basically
"wrappers on steroids" around specific Linux kernel capabilities, most notably ones
are:

- chroots
- namespaces
- cgroups

# Usage

From shell of a VM container download Alpine Linux "image":

```
./smocker pull
```
launch main container to run a shell limiting its memory consumption to 10MB:
```
./smocker run -m 10000000 /bin/sh
```
attach to main container from another shell:
```
./smocker exec /bin/sh
```

# Note
For safety only run the script inside a test VM.
