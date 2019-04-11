9t
==============================

[![Build Status](https://travis-ci.org/gongo/9t.svg?branch=master)](https://travis-ci.org/gongo/9t)
[![Coverage Status](https://coveralls.io/repos/gongo/9t/badge.svg?branch=master)](https://coveralls.io/r/gongo/9t?branch=master)

9t (nine-tailed fox in Japanese) is a multi-file tailer (like `tail -f a.log b.log ...`).

Usage
------------------------------

```
$ 9t file1 [file2 ...]
```

This folk version has additional options to work with named pipes.

```
9t -l lable1:file1 -p named_pipe1 -l -p label2:named_pipe2 [file2 ...]
```

You can follow logs without creating files
using [Bash Process Substitution](https://www.gnu.org/software/bash/manual/html_node/Process-Substitution.html).
For example, to work with `ping`, `kubectl logs` and `kubectl get events`:

```
9t \
    -p -l ping-local:<(ping 127.0.0.1) \
    -p -l somepod:<(kubectl -f logs somepod) \
    -p -l k8s-events:<(kubectl get events --watch)
```

### Demo

![Demo](./images/9t.gif)

1. Preparation for demo

    ```sh
    $ yukari() { echo '世界一かわいいよ!!' }
    $ while :; do       yukari >> tamura-yukari.log ; sleep 0.2 ; done
    $ while :; do echo $RANDOM >> random.log        ; sleep 3   ; done
    $ while :; do         date >>      d.log        ; sleep 1   ; done
    ```

1. Run

    ```
    $ 9t tamura-yukari.log random.log d.log
    ```

Installation
------------------------------

```
$ go get github.com/gongo/9t/cmd/9t
```

Motivation
------------------------------

So far, Multiple file display can be even `tail -f`.

![Demo](./images/tailf.gif)

But, I wanted to see in a similar format as the `heroku logs --tail`.

```
app[web.1]: foo bar baz
app[worker.1]: pizza pizza
app[web.1]: foo bar baz
app[web.2]: just do eat..soso..
.
.
```

License
------------------------------

MIT License
