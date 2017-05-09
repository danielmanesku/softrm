# softrm

By default there is no trash can in Unix CLI environment. Once you execute built in `rm args` all the files are gone.

Never again delete an important file by mistake. `softrm` is a CLI tool 
that offers trash can like capabilities (that are commonly found in desktop environments).

Note that this is still work in progress, a lot of functionality is missing at the moment.

## Available commands
```
softrm rm [args]      # move file(s) to trash
softrm list           # show trash contents (not implemented yet)
softrm restore [id]   # restore group of files deleted with [id] (not implemented yet)
softrm flush [id]     # permanently delete group of files [id]
```

In the future it is planned to add support for automatic flushing of files older than specified time.

## Other
Note that softrm is not intended to integrate with desktop environment's trash can. This is due to differences
in how softrm works.

## License
Apache License, Version 2.0
