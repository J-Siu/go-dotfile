# go-dotfile

Simple dotfile management command line.

### Concept and Limitation

- Not a drop-in replacement of Stow.
- Only top level directories and files are dotted in target location (`DirDest`)
- Symlink directory is copied as normal directory
- Files removed from source, will not be deleted from target location
- Files that should keep out of go-dotfile management
  - `~/.ssh/known_hosts`
  - history files
  - cache files

### Minimum GO version

go > 1.21

### Build and Install

```sh
git clone https://github.com/J-Siu/go-dotfile.git
```

This is build with go 1.24.5. If you have go < 1.24.5:

```sh
rm go.md go.sum
go mod init github.com/J-Siu/go-dotfile
go mod tidy
go get
```

To install

```sh
go install
```

### Configuration

Configuration must exist at `$HOME/.config/go-dotfile.json`, or supplied by the `-c` option.

Sample configuration:

```json
{
  "DirDest": "$HOME/tmp",
  "DirAP": [
    "$HOME/df_test/df_pub/append",
    "$HOME/df_test/df_pri/append"
  ],
  "DirCP": [
    "$HOME/df_test/df_pub/base",
    "$HOME/df_test/df_pri/base"
  ]
}
```

Variable|Default|Usage
--|--|--
DirDest|$HOME|Target location of dotfiles and directories
DirCP|n/a|Directories to be copied to target location
DirAP|n/a|Files in these directories will be be copied to target location if not already exist, else appended

### Testing

```sh
cp -r examples/df_test $HOME/
mkdir $HOME/tmp
go-dotfile -c examples/go-dotfile.sample.json
ls -a $HOME/tmp
```

### Change Log

- v1.0.7
  - Fix Github workflows
- v1.0.8
  - Upgrade go-helper to v1.1.8
  - Fix config logic
  - Fix debug logic
  - Fix receiver name
- v1.0.9
  - Upgrade go-helper to v1.1.10
  - Add command line version
  - Move -dryRun from base to "upgrade" command
  - Move default config file `~/.go-dotfile.json` -> `~/.config/go-dotfile.json`
- v1.0.10
  - Fix TypeConf.setDefault overwrite command line config file option
  - conf.go
    - check DirDest exist
  - dotfile.go
    - Init() - Chdir() error check
    - Process() - queue error
    - Change some function to local
  - root.go
    - Print error queue
- v1.0.11
  - Fix version
- v1.1.0
  - Skip copy if source and destination files have same modification time and size
  - Copy source file modification time to destination file
- v1.1.1
  - TypeDotfile
    - fix MyType mismatch
  - fileChanged
    - fix logical err: should ignore destination file stat() err

### License

The MIT License (MIT)

Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
