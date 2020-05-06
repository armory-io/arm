package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"testing"
)

func Test_dinghyRender(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		viper map[string]string
		want string
	}{
		//Basic dinghy file
		{ "TestBasicDinghy" , args{[]string{"../examples/dinghyfile_basic"}}, map[string]string{"modules": "../examples/modules"},
		`{
				  "application": "helloworldapp",
				  "pipelines": [
					{
					  "application": "helloworldapp",
					  "name": "my-pipeline-name",
					  "stages": [
						{
						  "name": "one",
						  "type": "wait",
						  "waitTIme": 10
						}
					  ]
					}
				  ]
				}`,
		},

		{ "TestConditionalsDinghy" , args{[]string{"../examples/dinghyfile_conditionals"}}, map[string]string{"modules": "../examples/modules", "rawdata": "../examples/RawData.json"},
			`{
  "application": "conditionals",
  "pipelines": [
    {
      "application": "conditionals",
      "name": "my-pipeline-name",
      "stages": [
        {
          "name": "this_is_true",
          "waitTime": 10,
          "type": "wait"
        }
      ]
    }
  ]
}`,
		},

		{ "TestGlobalsDinghy" , args{[]string{"../examples/dinghyfile_globals"}}, map[string]string{"modules": "../examples/modules"},
			`{
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
}`,
		},

		{ "TestMakeSliceDinghy" , args{[]string{"../examples/dinghyfile_makeSlice"}}, map[string]string{"modules": "../examples/modules"},
			`{
  "application": "makeSlice",
  "pipelines": [
    {
      "name": "Loop Example",
      "application": "makeSlice",
      "stages": [
        {
          "name": "First Wait",
          "waitTime": "10",
          "type": "wait"
        },
        {
          "name": "Second Wait",
          "waitTime": "10",
          "type": "wait"
        },
        {
          "name": "Final Wait",
          "waitTime": "10",
          "type": "wait"
        }
      ]
    }
  ]
}`,
		},

		{ "TestRawDataDinghy" , args{[]string{"../examples/dinghyfile_rawdata"}}, map[string]string{"modules": "../examples/modules", "rawdata": "../examples/RawData.json"},
			`{
  "application": "rawdata",
  "pipelines": [
    {
      "application": "rawdata",
      "name": "my-pipeline-name",
      "stages": [
        {
          "name": "Samtribal",
          "type": "wait",
          "waitTIme": 10
        }
      ]
    }
  ]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for keyViper := range tt.viper {
				viper.Set(keyViper, tt.viper[keyViper])
			}

			got := dinghyRender(tt.args.args)
			if got != "" {
				var gotBuffer bytes.Buffer
				var wantBuffer bytes.Buffer
				json.Indent(&gotBuffer, []byte(got), "", "  ")
				json.Indent(&wantBuffer, []byte(tt.want), "", "  ")
				if gotBuffer.String() != wantBuffer.String() {
					t.Errorf("dinghyRender() = %v, want %v", got, tt.want)
				}
			} else {
				t.Errorf("dinghyRender() = %v, want %v", got, tt.want)
			}
		})
	}
}
