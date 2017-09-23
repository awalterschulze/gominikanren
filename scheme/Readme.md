# Docker Scheme with a copy of microKanren for sanity checks

This folder is meant for sanity checking the scheme version of microKanren against the Go implementation of microKanren.

This folder contains:
  - Scheme Docker
  - Copy of MicroKanren
  - Makefile to tie the two together

This repo is meant for a Go implementation of microKanren, 
so we assume that most developers will not have scheme installed.
This is why we include a dockered version of Scheme.
The Dockerfile was found here: https://hub.docker.com/r/kleinpa/chez/

The version of microKanren has been copied from https://github.com/jasonhemann/microKanren
this is why a LICENSE file has been included.

Lastly the Makefile is meant to tie the two together and make it easy to run the scheme version of microKanren, even if you do not have scheme installed.

Run miniKanren with the code included in hellokanren.scm

    make run

Run the scheme repl with miniKanren loaded

    make repl

Exit the repl

    (exit)

Run the code in helloworld.scm

    make helloworld

Discover more about scheme

    make help
