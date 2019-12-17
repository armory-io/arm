# arm

The Armory CLI

To use this clone the project.
Make sure to have go installed [https://golang.org/dl/](https://golang.org/dl/)

cd into the project and run `go build`

Add the tool to your path by doing
ln -sf $PWD/arm /usr/local/bin/arm  

to use it just type arm into your console

```bash
$ arm
The Armory Platform CLI

Usage:
  arm [command]

Available Commands:
  dinghy      Run dinghy subcommands
  help        Help about any command

Flags:
  -h, --help   help for arm

Use "arm [command] --help" for more information about a command.
```

## Sample Usage
Both the Dinghyfile and module repo must be available locally

For example, if we have a dinghyfile at `/tmp/dinghyfile`, and a repo at `git@github.com:justinrlee/sales-pipelines-modules`, then you can do this:

```bash
# Clone the module repository
$ git clone git@github.com/justinrlee/sales-pipelines-modules /tmp/sales-pipelines-modules

# Both the /tmp/dinghyfile and /tmp/sales-pipelines-modules support either relative or absolute paths
# The first item is the dinghyfile file we are rendering
# The modules item should be a directory; in can be modules.  Modules can be in subdirectories; in this case, the modules should be referenced in the dinghyfile by their relative path within the module repository

$ arm dinghy render /tmp/dinghyfile --modules /tmp/sales-pipelines-modules
INFO[0000] Output:
{
  "application": "multi-many",
  "pipelines": [
...
  ]
}
```

You can pipe this through jq to validate that it produces valid json:
```bash
arm dinghy render dinghyfile --modules test/modules | jq '.'
```
