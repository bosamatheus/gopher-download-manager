# Gopher Download Manager

Concurrent Download Manager built with Golang.

# Execution

To execute the program, you need to run the following command with arguments:

```bash
$ go run main.go <url> <filename> <threads>
```

For example, you can run the following command to download a beautiful gopher image:

```bash
$ go run main.go https://golang.org/doc/gopher/run.png gopher.png 10
```
