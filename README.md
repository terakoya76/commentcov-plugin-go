# commentcov-plugin-go

[commentcov-plugin-go](https://github.com/commentcov/commentcov-plugin-go) is the [commentcov](https://github.com/commentcov/commentcov) plugin for go

## How to use

Specify plugin info on your commentcov configuraiton (default is `<path_to_project>/.commentcov.yaml`).

```bash
$ pwd
/home/commentcov/workspace/github.com/kubernetes/kubernetes

$ cat .commentcov.yaml
plugins:
  - extension: .go
    install_command: go install github.com/commentcov/commentcov-plugin-go@latest
    execute_command: commentcov-plugin-go
exclude_paths:
  - "vendor/**/**"
mode: scope
```

Then, run commentcov. You could get the comment coverage of the `.go` files
```bash
$ commentcov coverage
{"@level":"info","@message":"Install Plugin","@module":"commentcov","@timestamp":"2022-06-04T09:49:57.276670+09:00","plugin":"commentcov-plugin-for-go"}
{"@level":"debug","@message":"starting plugin","@module":"commentcov","@timestamp":"2022-06-04T09:49:59.637761+09:00","args":["commentcov-plugin-go"],"path":"/home/commentcov/go/bin/commentcov-plugin-go"}
{"@level":"debug","@message":"plugin started","@module":"commentcov","@timestamp":"2022-06-04T09:49:59.637963+09:00","path":"/home/commentcov/go/bin/commentcov-plugin-go","pid":302317}
{"@level":"debug","@message":"waiting for RPC address","@module":"commentcov","@timestamp":"2022-06-04T09:49:59.638000+09:00","path":"/home/commentcov/go/bin/commentcov-plugin-go"}
{"@level":"debug","@message":"plugin address","@module":"commentcov.commentcov-plugin-go","@timestamp":"2022-06-04T09:49:59.641455+09:00","address":"/tmp/plugin3894134805","network":"unix","timestamp":"2022-06-04T09:49:59.641+0900"}
{"@level":"debug","@message":"using plugin","@module":"commentcov","@timestamp":"2022-06-04T09:49:59.641511+09:00","version":1}
{"@level":"trace","@message":"waiting for stdio data","@module":"commentcov.stdio","@timestamp":"2022-06-04T09:49:59.642178+09:00"}
PUBLIC_TYPE: 86.67992047713717
PRIVATE_FUNCTION: 26.929215170859933
PUBLIC_FUNCTION: 46.19191641462744
PUBLIC_CLASS: 86.1316662413881
PUBLIC_VARIABLE: 69.15447974449488
PRIVATE_TYPE: 36.721311475409834
FILE: 8.397812854637367
PRIVATE_CLASS: 38.396509408235616
PRIVATE_VARIABLE: 14.998234463276836
```

## Commentcov CoverageItem Scope Mapping

The mapping from the type of Go code comment to the CoverageItem Scope is below.

| Scope Name                    | Golang Node                          |
|-------------------------------|--------------------------------------|
| CoverageItem_UNKNONW          | N/A                                  |
| CoverageItem_FILE             | Package Comment                      |
| CoverageItem_PUBLIC_MODULE    | doc.go (Not Supported yet)           |
| CoverageItem_PRIVATE_MODULE   | N/A                                  |
| CoverageItem_PUBLIC_CLASS     | Exported Struct, Interface Comment   |
| CoverageItem_PRIVATE_CLASS    | Unexported Struct, Interface Comment |
| CoverageItem_PUBLIC_TYPE      | Exported Type Alias Comment          |
| CoverageItem_PRIVATE_TYPE     | Unexported Type Alias Comment        |
| CoverageItem_PUBLIC_FUNCTION  | Exported Function Comment            |
| CoverageItem_PRIVATE_FUNCTION | Unexported Function Comment          |
| CoverageItem_PUBLIC_VARIABLE  | Exported Var, Const Comment          |
| CoverageItem_PRIVATE_VARIABLE | Unexported Var, Const Comment        |

