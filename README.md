go-v8
=====

v8 javascript engine binding for golang.

How to install
==============

Need 'Go 1.2' for auto compile C++ file.

Use the command below to auto download and compile go-v8:

```
git init && git remote add origin https://github.com/realint/go-v8 && git pull origin master && ./install.sh
```

Tools
=====

There is a shell script named 'install.sh'. It can auto download and compile v8 engine.

For example:

```
./install.sh
```

And there is a shell script named 'test.sh'. It can auto configure cgo environment variable and run go-v8 test.

For example:

```
./test.sh . .
```

The first argument of test.sh is test name pattern, second argument is benchmark name pattern.

For example:

```
./test.sh .
```

The command above will run test without benchmark.


