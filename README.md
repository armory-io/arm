# Arm installation

The Armory CLI

To use this you can download the binaries fom the [release section](https://github.com/armory-io/arm/releases) and unzip the files depending on your OS.

You can add the tool to your path by doing
ln -sf $PWD/arm /usr/local/bin/arm  

To use it just type arm in your console(if you added it to your path) or execute arm-cli directly.


```bash
➜  arm-0.0.3-osx-amd64 ./arm-0.0.3-darwin-amd64
The Armory Platform CLI

Usage:
  arm [command]

Available Commands:
  dinghy      Run dinghy subcommands
  help        Help about any command
  version     Prints version information

Flags:
  -h, --help              help for arm
  -l, --loglevel string   log level (default "info")

Use "arm [command] --help" for more information about a command.
```

## Mac OS
Since arm-cli is not signed you may receive a couple of messages regarding security. To execute the binaries you need to:

1. When you execute arm-cli on the console you may see this message.

<img src="docs/img/01_developer_verification.jpg" width="50%" />

2. Open Spotlight search in your Mac OS and search for Security & Privacy. 

<img src="docs/img/02_open_privacy.jpg" width="50%" />

3. Once you opened this option go to General Tab and you will see on the bottom a button to "Allow Anyway" pointing at the bin, click on it.
 
 <img src="docs/img/03_allow_anyway.jpg" width="50%" />

4. Once you click on it the button will dissapear.

5. Try to execute again the binary in your console and arm-cli should be working.

 <img src="docs/img/05_working.jpg" width="50%" />

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

If final json file is valid you can see the message `Validation passed`, this means that the final JSON object is valid.

