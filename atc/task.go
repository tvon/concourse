package atc

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

type TaskConfig struct {
	// The platform the task must run on (e.g. linux, windows).
	Platform string `yaml:"platform,omitempty"`

	// Optional string specifying an image to use for the build. Depending on the
	// platform, this may or may not be required (e.g. Windows/OS X vs. Linux).
	RootfsURI string `yaml:"rootfs_uri,omitempty"`

	ImageResource *ImageResource `yaml:"image_resource,omitempty"`

	// Limits to set on the Task Container
	Limits ContainerLimits `yaml:"container_limits,omitempty"`

	// Parameters to pass to the task via environment variables.
	Params map[string]string `yaml:"params,omitempty"`

	// Script to execute.
	Run TaskRunConfig `yaml:"run,omitempty"`

	// The set of (logical, name-only) inputs required by the task.
	Inputs []TaskInputConfig `yaml:"inputs,omitempty"`

	// The set of (logical, name-only) outputs provided by the task.
	Outputs []TaskOutputConfig `yaml:"outputs,omitempty"`

	// Path to cached directory that will be shared between builds for the same task.
	Caches []CacheConfig `yaml:"caches,omitempty"`
}

type ContainerLimits struct {
	CPU    *uint64 `yaml:"cpu,omitempty"`
	Memory *uint64 `yaml:"memory,omitempty"`
}

type ImageResource struct {
	Type   string `yaml:"type"`
	Source Source `yaml:"source"`

	Params  *Params  `yaml:"params,omitempty"`
	Version *Version `yaml:"version,omitempty"`
}

func NewTaskConfig(configBytes []byte) (TaskConfig, error) {
	var config TaskConfig
	err := yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return TaskConfig{}, err
	}

	err = config.Validate()
	if err != nil {
		return TaskConfig{}, err
	}

	return config, nil
}

func (config TaskConfig) Validate() error {
	messages := []string{}

	if config.Platform == "" {
		messages = append(messages, "  missing 'platform'")
	}

	if config.Run.Path == "" {
		messages = append(messages, "  missing path to executable to run")
	}

	messages = append(messages, config.validateInputsAndOutputs()...)

	if len(messages) > 0 {
		return fmt.Errorf("invalid task configuration:\n%s", strings.Join(messages, "\n"))
	}

	return nil
}

func (config TaskConfig) validateInputsAndOutputs() []string {
	messages := []string{}

	messages = append(messages, config.validateInputContainsNames()...)
	messages = append(messages, config.validateOutputContainsNames()...)

	return messages
}

func (config TaskConfig) validateDotPath() []string {
	messages := []string{}

	pathCount := 0
	dotPath := false

	for _, input := range config.Inputs {
		path := strings.TrimPrefix(input.resolvePath(), "./")

		if path == "." {
			dotPath = true
		}

		pathCount++
	}

	for _, output := range config.Outputs {
		path := strings.TrimPrefix(output.resolvePath(), "./")

		if path == "." {
			dotPath = true
		}

		pathCount++
	}

	if pathCount > 1 && dotPath {
		messages = append(messages, "  you may not have more than one input or output when one of them has a path of '.'")
	}

	return messages
}

type pathCounter struct {
	inputCount  map[string]int
	outputCount map[string]int
}

func (counter *pathCounter) foundInBoth(path string) bool {
	_, inputFound := counter.inputCount[path]
	_, outputFound := counter.outputCount[path]

	return inputFound && outputFound
}

func (counter *pathCounter) registerInput(input TaskInputConfig) {
	path := strings.TrimPrefix(input.resolvePath(), "./")

	if val, found := counter.inputCount[path]; !found {
		counter.inputCount[path] = 1
	} else {
		counter.inputCount[path] = val + 1
	}
}

func (counter *pathCounter) registerOutput(output TaskOutputConfig) {
	path := strings.TrimPrefix(output.resolvePath(), "./")

	if val, found := counter.outputCount[path]; !found {
		counter.outputCount[path] = 1
	} else {
		counter.outputCount[path] = val + 1
	}
}

func (config TaskConfig) validateOutputContainsNames() []string {
	messages := []string{}

	for i, output := range config.Outputs {
		if output.Name == "" {
			messages = append(messages, fmt.Sprintf("  output in position %d is missing a name", i))
		}
	}

	return messages
}

func (config TaskConfig) validateInputContainsNames() []string {
	messages := []string{}

	for i, input := range config.Inputs {
		if input.Name == "" {
			messages = append(messages, fmt.Sprintf("  input in position %d is missing a name", i))
		}
	}

	return messages
}

type TaskRunConfig struct {
	Path string   `yaml:"path"`
	Args []string `yaml:"args,omitempty"`
	Dir  string   `yaml:"dir,omitempty"`

	// The user that the task will run as (defaults to whatever the docker image specifies)
	User string `yaml:"user,omitempty"`
}

type TaskInputConfig struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path,omitempty"`
	Optional bool   `yaml:"optional,omitempty"`
}

func (input TaskInputConfig) resolvePath() string {
	if input.Path != "" {
		return input.Path
	}
	return input.Name
}

type TaskOutputConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path,omitempty"`
}

func (output TaskOutputConfig) resolvePath() string {
	if output.Path != "" {
		return output.Path
	}
	return output.Name
}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CacheConfig struct {
	Path string `yaml:"path,omitempty"`
}
