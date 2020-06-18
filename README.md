# xexec
A Go library that allows you to abstract the operations related to process execution.

# Unit Test
The unit tests depend on some simple executables. You can use the command below to build the executables.
```
mkdir bin
cd bin
go build ../internal/...
cd ..
```
Then, you can run the unit tests as below:
```
go test ./...
```
The unit tests use *mockery* + *testify* for mocking. If you want to generate new mock types, see 
[here](https://github.com/vektra/mockery) for the way to install *mockery*.
