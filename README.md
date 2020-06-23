# xexec
It's a Go library that allows you to abstract the operations related to process execution.  
![CI Status](https://github.com/absgo/xexec/workflows/CI/badge.svg)
## Install
Install using the following command:
```bash
go get -u github.com/absgo/xexec
```

## Usage
You can create an instance of `Exec` as below:
```go
import github.com/absgo/xexec
func main() {
    exec := xexec.NewOsExec()
}
```
Then, all the operations related to process execution can be achieved via the `Exec` instance.  
For example, to start a process, you can use:
```go
proc, err := exec.StartProcess(&ProcessConf{
		Name: "ls",
		Args: []string{"ls", "/"},
	})
if err != nil {
    log.Fatalf("failed to start process: %v", err)
}
_, err = proc.Wait()
if err != nil {
    log.Fatalf("failed to wait on process: %v", err)
}
log.Print("Success")
```

## Unit Test
If you want to run the unit tests by yourself, follow the steps below.  
The unit tests depend on some simple executables. You can use the command below to build the executables.
```bash
mkdir bin
cd bin
go build ../internal/...
cd ..
```
Then, you can run the unit tests as below:
```bash
go test ./...
```
The unit tests use *mockery* + *testify* for mocking. If you want to generate new mock types, see 
[here](https://github.com/vektra/mockery) for the way to install *mockery*.  
There're `go:genenerate` commands in the files named `generate.go`.  You can use the following command to 
generate the mock types again:
```go
go generate ./...
```
