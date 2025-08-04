Simple dotfile management command line.

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

### Concept and Limitation

- All top level directories and files are dotted in target location (`DirDest`)
- Cannot handle symlink directory

### Misc

Files that should keep out of dotfile management
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