# go-dotfile

Simple dotfile management command line.

- [Concept and Limitation](#concept-and-limitation)
- [Install](#install)
- [Build](#build)
- [Configuration](#configuration)
- [Testing](#testing)
- [License](#license)

<!--more-->

### Concept and Limitation

- Not a drop-in replacement of Stow.
- Only top level directories and files are dotted in target location (`DirDest`)
- Symlink directory is copied as normal directory
- Files removed from source, will not be deleted from target location
- Files that should keep out of go-dotfile management
  - `~/.ssh/known_hosts`
  - history files
  - cache files

### Install

Go install

```sh
go install github.com/J-Siu/go-dotfile@latest
```

Download

- https://github.com/J-Siu/go-dotfile/releases

### Build

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

Install

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

### License

The MIT License (MIT)

Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
