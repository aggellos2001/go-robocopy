# go-robocopy

## Disclaimer

> This is an unofficial wrapped for robocopy. There is no affiliation with Microsoft or any of its subsidiaries in any way. Use it at your own risk.

Robocopy or Robust File Copy is a command line tool developed by Microsoft for copying files and directories from one location to another.

This is an **unofficial wrapper** written in Go to programmatically produce a robocopy command and execute it.

## Installation

You can install the package using the go get command.

```bash
go get github.com/aggellos2001/go-robocopy
```

After that you need to import the module in your code.

```go

import (
    gorobocopy "github.com/aggellos2001/go-robocopy"
)
```

## Examples

You can create a new command instance by calling the NewRobocopy constructor.

```go
cmd := gorobocopy.NewRobocopy(
    "C:\\source",
    "D:\\destination",
    "*.*",
)
```

You can then set some options for the command.

```go
cmd.SetCopyOptions(&gorobocopy.CopyOptions{
    E:    true,
    Mt:   4,
    Copy: copyflags.D | copyflags.A | copyflags.T,
    }
)
```

Finally, you can execute the command. You can specify the stdin, stdout and stderr in the parameters or leave them as nil if you want to suppress the console input/output.

```go
cmd.Run(nil,nil,nil)
```

You can also get the command object and call it yourself however you like. It returns a pointer to a *exec.Cmd object.

```go
cmd.GetCommand()
```

Also you can get the slice containing all the arguments to manually create the command if you please.

```go
cmd.GetCommandArgs()
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
