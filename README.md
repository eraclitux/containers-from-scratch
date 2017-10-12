# containers-from-scratch

Proof of concept implementations that simulate main functionalities
of Docker.
Intended only to deep knowledge about GNU/Linux containers not for real world
usage.

Different tools for container management (Docker, rkt) on GNU/Linux are basically
"wrappers on steroids" around specific Linux kernel capabilities, most notably ones
are:

- chroots
- namespaces
- cgroups

# Note
For safety only execute this code inside a test VM.
