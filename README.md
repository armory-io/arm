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
Both the Dinghyfile and module repo must be available locally, there is an example folder build in.

For example, if we use an example from the folder we can:


```bash
$ arm dinghy render ./examples/dinghyfile_globals --modules ./examples/modules --rawdata ./examples/RawData.json --output ./testing
INFO[0000] Checking dinghyfile                          
INFO[0000] Reading rawdata file                         
INFO[0000] Parsing rawdata json                         
INFO[0000] Parsing dinghyfile                           
INFO[0000] Parsed dinghyfile                            
INFO[0000] Validating output json.                      
INFO[0000] Saving output file                           
INFO[0000] Output:                                      
{
  "application": "global_vars",
  "globals": {
    "waitTime": "42",
    "waitname": "default-name"
  },
  "pipelines": [
    {
      "application": "global_vars",
      "name": "Made By Armory Pipeline Templates",
      "stages": [
        {
          "name": "default-name",
          "waitTime": "42",
          "type": "wait"
        },
        {
          "name": "overwrite-name",
          "waitTime": "100",
          "type": "wait"
        }
      ]
    }
  ]
}
```

If final json file is valid you can see the mesage `Validation passed`, this means that the final JSON object is valid.