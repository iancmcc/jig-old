# jig
A multirepo tool

jig handles the overhead of multiple repositories. It will clone them, keep
them in a sane structure, execute commands against them, and navigate between
them with only a few simple commands.

All command-line options can also be passed as environment variables, for
recreating environments programmatically.


* ``jig up [-f JIGFILE] [-j JIGROOT]``

    Realizes a Jigfile (serialized representation) to disk under JIGROOT. Will
    use Jigfile in the current directory if not specified. As a side benefit,
    executing this in an existing JIGROOT will fetch remote changes for all
    repositories.

* ``jig get URI [-j JIGROOT]``

    Check out the repository at URI into JIGROOT (defaults to closest ancestor
    JIGROOT)

* ``jig dir QUERY [-j JIGROOT]``

    Find a repository with fuzzy-matched SEARCHSTRING under JIGROOT (defaults to
    closest ancestor jig). Outputs the path to that repo. With no arguments,
    outputs JIGROOT.

* ``jig do [-q SEARCHSTRING [-q SEARCHSTRING]] [-r REPO [-r REPO]] [-j JIGROOT] COMMAND``
    For every repository matching the search string or repo specified, execute
    COMMAND in that directory

* ``jig save [-j JIGROOT] [-f JIGFILE]``

    Save the state of an existing JIGROOT to a JIGFILE so it can be reproduced
    later with “jig up".

* ``jig bootstrap``

    Prints the path to a bootstrap file that can be sourced in your shell to
    provide autocompletion and utility functions to navigate among
    repositories. Add ``source $(jig bootstrap)`` to ``~/.bashrc``.


## Build
Building can be accomplished either by using ``go get`` or my using Go 1.5
vendoring support. To do the latter, set:

    GO15VENDOREXPERIMENT=1

Then ``go build`` as normal.

