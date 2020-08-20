package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/magiconair/properties/assert"
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
		errorMsg error
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
				nil,
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
nil,
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
nil,
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
nil,
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
nil,
		},

		{ "TestLocalModulesDinghy" , args{[]string{"../examples/dinghyfile_localmodule"}}, map[string]string{"modules": "../examples/modules", "rawdata": "../examples/RawData.json"},
			`{
  "application": "localmodules",
  "globals": {
    "waitTime": "42",
    "waitname": "localmodule default-name"
  },
  "pipelines": [
    {
      "application": "localmodules",
      "name": "Made By Armory Pipeline Templates",
      "stages": [
        {
          "name": "localmodule default-name",
          "waitTime": "42",
          "type": "wait"
        },
        {
          "name": "localmodule overwrite-name",
          "waitTime": "100",
          "type": "wait"
        },
        {
          "name": "global module overwrite-name",
          "waitTime": "100",
          "type": "wait"
        }
      ]
    }
  ]
}`,
nil,
		},
		{ "TestLocalModulesWithParameter" , args{[]string{"../examples/dinghyfile_localmodule_parameter"}}, map[string]string{"modules": "../examples/modules", "rawdata": "../examples/RawData.json", "local_modules" : "../"},
			`{
  "application": "localmodules",
  "globals": {
    "waitTime": "42",
    "waitname": "localmodule default-name"
  },
  "pipelines": [
    {
      "application": "localmodules",
      "name": "Made By Armory Pipeline Templates",
      "stages": [
        {
          "name": "localmodule default-name",
          "waitTime": "42",
          "type": "wait"
        },
        {
          "name": "localmodule overwrite-name",
          "waitTime": "100",
          "type": "wait"
        },
        {
          "name": "global module overwrite-name",
          "waitTime": "100",
          "type": "wait"
        }
      ]
    }
  ]
}`,
			nil,
		},
		{ "TestValidatePipelines" , args{[]string{"../test_dinghyfiles/TestValidatePipelines"}}, map[string]string{"modules": "../examples/modules", "rawdata": "../examples/RawData.json"},
			``,
			errors.New("Final Dinghyfile failed validations, please correct them and retry. mj2 refers to itself. Circular references are not supported"),
		},
		{ "TestValidateAppNotifications" , args{[]string{"../test_dinghyfiles/TestValidateAppNotifications"}}, map[string]string{"modules": "../examples/modules", "rawdata": "../examples/RawData.json"},
			``,
			errors.New("Final Dinghyfile failed application notification validations, please correct them and retry. application notifications format is invalid for email"),
		},
		{ "TestPipelineID" , args{[]string{"../examples/dinghyfile_pipelineID"}}, map[string]string{"modules": "../examples/modules", "rawdata": "../examples/RawData.json"},
			`{
  "application": "pipelineidexample",
  "pipelines": [
    {
      "application": "pipelineidexample",
      "name": "my-pipeline-name",
      "stages": [
        {
          "name": "reference pipeline",
          "application": "default-app",
          "pipeline": "mock-default-pipeline-id",
          "type": "pipeline"
        }
      ]
    }
  ]
}`,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for keyViper := range tt.viper {
				viper.Set(keyViper, tt.viper[keyViper])
			}

			got, err := dinghyRender(tt.args.args)
			if err != nil {
				if tt.errorMsg != nil {
					assert.Equal(t, tt.errorMsg, err)
				} else  {
					t.Fail()
				}
			} else {
				if got != "" {
					var gotBuffer bytes.Buffer
					var wantBuffer bytes.Buffer
					json.Indent(&gotBuffer, []byte(got), "", "  ")
					json.Indent(&wantBuffer, []byte(tt.want), "", "  ")
					if gotBuffer.String() != wantBuffer.String() {
						t.Errorf("dinghyRender() = %v, want %v", got, tt.want)
						t.Fail()
					}
				} else {
					t.Errorf("dinghyRender() = %v, want %v", got, tt.want)
					t.Fail()
				}
			}
		})
	}
}
