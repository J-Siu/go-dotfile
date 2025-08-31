# go-dotfile

Simple dotfile management command line.

### Concept and Limitation

- Only top level directories and files are dotted in target location (`DirDest`)
- Symlink directory is copied as normal directory

### Minimum GO version

go > 1.21

### Install

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

To install

```sh
go install
```

### Configuration

`go-dotfile.sample.conf`

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

### Testing

```sh
cp -r examples/df_test $HOME/
mkdir $HOME/tmp
go-dotfile -c examples/go-dotfile.sample.json
ls -a $HOME/tmp
```

### Misc

Files that should keep out of go-dotfile management

- `~/.ssh/known_hosts`
- history files
- cache files

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