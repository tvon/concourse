package atc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const ConfigVersionHeader = "X-Concourse-Config-Version"
const DefaultPipelineName = "main"
const DefaultTeamName = "main"

type Tags []string

type Config struct {
	Groups        GroupConfigs    `yaml:"groups,omitempty"`
	Resources     ResourceConfigs `yaml:"resources,omitempty"`
	ResourceTypes ResourceTypes   `yaml:"resource_types,omitempty"`
	Jobs          JobConfigs      `yaml:"jobs,omitempty"`
}

type GroupConfig struct {
	Name      string   `yaml:"name"`
	Jobs      []string `yaml:"jobs,omitempty"`
	Resources []string `yaml:"resources,omitempty"`
}

type GroupConfigs []GroupConfig

func (groups GroupConfigs) Lookup(name string) (GroupConfig, int, bool) {
	for index, group := range groups {
		if group.Name == name {
			return group, index, true
		}
	}

	return GroupConfig{}, -1, false
}

type ResourceConfig struct {
	Name         string  `yaml:"name"`
	Public       bool    `yaml:"public,omitempty"`
	WebhookToken string  `yaml:"webhook_token,omitempty"`
	Type         string  `yaml:"type"`
	Source       Source  `yaml:"source"`
	CheckEvery   string  `yaml:"check_every,omitempty"`
	CheckTimeout string  `yaml:"check_timeout,omitempty"`
	Tags         Tags    `yaml:"tags,omitempty"`
	Version      Version `yaml:"version,omitempty"`
}

type ResourceType struct {
	Name                 string `yaml:"name"`
	Type                 string `yaml:"type"`
	Source               Source `yaml:"source"`
	Privileged           bool   `yaml:"privileged,omitempty"`
	CheckEvery           string `yaml:"check_every,omitempty"`
	Tags                 Tags   `yaml:"tags,omitempty"`
	Params               Params `yaml:"params,omitempty"`
	CheckSetupError      string `yaml:"check_setup_error,omitempty"`
	CheckError           string `yaml:"check_error,omitempty"`
	UniqueVersionHistory bool   `yaml:"unique_version_history,omitempty"`
}

type ResourceTypes []ResourceType

func (types ResourceTypes) Lookup(name string) (ResourceType, bool) {
	for _, t := range types {
		if t.Name == name {
			return t, true
		}
	}

	return ResourceType{}, false
}

func (types ResourceTypes) Without(name string) ResourceTypes {
	newTypes := ResourceTypes{}
	for _, t := range types {
		if t.Name != name {
			newTypes = append(newTypes, t)
		}
	}

	return newTypes
}

type Hooks struct {
	Abort   *PlanConfig
	Failure *PlanConfig
	Ensure  *PlanConfig
	Success *PlanConfig
}

// A PlanSequence corresponds to a chain of Compose plan, with an implicit
// `on: [success]` after every Task plan.
type PlanSequence []PlanConfig

// A VersionConfig represents the choice to include every version of a
// resource, the latest version of a resource, or a pinned (specific) one.
type VersionConfig struct {
	Every  bool
	Latest bool
	Pinned Version
}

func (c *VersionConfig) UnmarshalJSON(version []byte) error {
	var data interface{}

	err := json.Unmarshal(version, &data)
	if err != nil {
		return err
	}

	switch actual := data.(type) {
	case string:
		c.Every = actual == "every"
		c.Latest = actual == "latest"
	case map[string]interface{}:
		version := Version{}

		for k, v := range actual {
			if s, ok := v.(string); ok {
				version[k] = strings.TrimSpace(s)
			}
		}

		c.Pinned = version
	default:
		return errors.New("unknown type for version")
	}

	return nil
}

func (c *VersionConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data interface{}

	err := unmarshal(&data)
	if err != nil {
		return err
	}

	switch actual := data.(type) {
	case string:
		c.Every = actual == "every"
		c.Latest = actual == "latest"
	case map[interface{}]interface{}:
		version := Version{}

		for k, v := range actual {
			if ks, ok := k.(string); ok {
				if vs, ok := v.(string); ok {
					version[ks] = strings.TrimSpace(vs)
				}
			}
		}

		c.Pinned = version
	default:
		return errors.New("unknown type for version")
	}

	return nil
}

func (c *VersionConfig) MarshalYAML() (interface{}, error) {
	if c.Latest {
		return VersionLatest, nil
	}

	if c.Every {
		return VersionEvery, nil
	}

	if c.Pinned != nil {
		return c.Pinned, nil
	}

	return nil, nil
}

func (c *VersionConfig) MarshalJSON() ([]byte, error) {
	if c.Latest {
		return json.Marshal(VersionLatest)
	}

	if c.Every {
		return json.Marshal(VersionEvery)
	}

	if c.Pinned != nil {
		return json.Marshal(c.Pinned)
	}

	return json.Marshal("")
}

// A InputsConfig represents the choice to include every artifact within the
// job as an input to the put step or specific ones.
type InputsConfig struct {
	All       bool
	Specified []string
}

func (c *InputsConfig) UnmarshalJSON(inputs []byte) error {
	var data interface{}

	err := json.Unmarshal(inputs, &data)
	if err != nil {
		return err
	}

	switch actual := data.(type) {
	case string:
		c.All = actual == "all"
	case []interface{}:
		inputs := []string{}

		for _, v := range actual {
			str, ok := v.(string)
			if !ok {
				return fmt.Errorf("non-string put input: %v", v)
			}

			inputs = append(inputs, strings.TrimSpace(str))
		}

		c.Specified = inputs
	default:
		return errors.New("unknown type for put inputs")
	}

	return nil
}

func (c *InputsConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data interface{}

	err := unmarshal(&data)
	if err != nil {
		return err
	}

	switch actual := data.(type) {
	case string:
		c.All = actual == "all"
	case []interface{}:
		inputs := []string{}

		for _, v := range actual {
			str, ok := v.(string)
			if !ok {
				return fmt.Errorf("non-string put input: %v", v)
			}

			inputs = append(inputs, strings.TrimSpace(str))
		}

		c.Specified = inputs
	default:
		return errors.New("unknown type for put inputs")
	}

	return nil
}

func (c InputsConfig) MarshalYAML() (interface{}, error) {
	if c.All {
		return InputsAll, nil
	}

	if c.Specified != nil {
		return c.Specified, nil
	}

	return nil, nil
}

func (c InputsConfig) MarshalJSON() ([]byte, error) {
	if c.All {
		return json.Marshal(InputsAll)
	}

	if c.Specified != nil {
		return json.Marshal(c.Specified)
	}

	return json.Marshal("")
}

// A PlanConfig is a flattened set of configuration corresponding to
// a particular Plan, where Source and Version are populated lazily.
type PlanConfig struct {
	// makes the Plan conditional
	// conditions on which to perform a nested sequence

	// compose a nested sequence of plans
	// name of the nested 'do'
	RawName string `yaml:"name,omitempty"`

	// a nested chain of steps to run
	Do *PlanSequence `yaml:"do,omitempty"`

	// corresponds to an Aggregate plan, keyed by the name of each sub-plan
	Aggregate *PlanSequence `yaml:"aggregate,omitempty"`

	// corresponds to Get and Put resource plans, respectively
	// name of 'input', e.g. bosh-stemcell
	Get string `yaml:"get,omitempty"`
	// jobs that this resource must have made it through
	Passed []string `yaml:"passed,omitempty"`
	// whether to trigger based on this resource changing
	Trigger bool `yaml:"trigger,omitempty"`

	// name of 'output', e.g. rootfs-tarball
	Put string `yaml:"put,omitempty"`

	// corresponding resource config, e.g. aws-stemcell
	Resource string `yaml:"resource,omitempty"`

	// inputs to a put step either a list (e.g. [artifact-1, aritfact-2]) or all (e.g. all)
	Inputs *InputsConfig `yaml:"inputs,omitempty"`

	// corresponds to a Task plan
	// name of 'task', e.g. unit, go1.3, go1.4
	Task string `yaml:"task,omitempty"`
	// run task privileged
	Privileged bool `yaml:"privileged,omitempty"`
	// task config path, e.g. foo/build.yml
	TaskConfigPath string `yaml:"file,omitempty"`
	// task variables, if task is specified as external file via TaskConfigPath
	TaskVars Params `yaml:"vars,omitempty"`
	// inlined task config
	TaskConfig *TaskConfig `yaml:"config,omitempty"`

	// used by Get and Put for specifying params to the resource
	Params Params `yaml:"params,omitempty"`

	// used to pass specific inputs/outputs as generic inputs/outputs in task config
	InputMapping  map[string]string `yaml:"input_mapping,omitempty"`
	OutputMapping map[string]string `yaml:"output_mapping,omitempty"`

	// used to specify an image artifact from a previous build to be used as the image for a subsequent task container
	ImageArtifactName string `yaml:"image,omitempty"`

	// used by Put to specify params for the subsequent Get
	GetParams Params `yaml:"get_params,omitempty"`

	// used by any step to specify which workers are eligible to run the step
	Tags Tags `yaml:"tags,omitempty"`

	// used by any step to run something when the build is aborted during execution of the step
	Abort *PlanConfig `yaml:"on_abort,omitempty"`

	// used by any step to run something when the step reports a failure
	Failure *PlanConfig `yaml:"on_failure,omitempty"`

	// used on any step to always execute regardless of the step's completed state
	Ensure *PlanConfig `yaml:"ensure,omitempty"`

	// used on any step to execute on successful completion of the step
	Success *PlanConfig `yaml:"on_success,omitempty"`

	// used on any step to swallow failures and errors
	Try *PlanConfig `yaml:"try,omitempty"`

	// used on any step to interrupt the step after a given duration
	Timeout string `yaml:"timeout,omitempty"`

	// not present in yaml
	DependentGet string `yaml:"-" json:"-"`

	// repeat the step up to N times, until it works
	Attempts int `yaml:"attempts,omitempty"`

	Version *VersionConfig `yaml:"version,omitempty"`
}

func (config PlanConfig) Name() string {
	if config.RawName != "" {
		return config.RawName
	}

	if config.Get != "" {
		return config.Get
	}

	if config.Put != "" {
		return config.Put
	}

	if config.Task != "" {
		return config.Task
	}

	return ""
}

func (config PlanConfig) ResourceName() string {
	resourceName := config.Resource
	if resourceName != "" {
		return resourceName
	}

	resourceName = config.Get
	if resourceName != "" {
		return resourceName
	}

	resourceName = config.Put
	if resourceName != "" {
		return resourceName
	}

	panic("no resource name!")
}

func (config PlanConfig) Hooks() Hooks {
	return Hooks{Abort: config.Abort, Failure: config.Failure, Ensure: config.Ensure, Success: config.Success}
}

type ResourceConfigs []ResourceConfig

func (resources ResourceConfigs) Lookup(name string) (ResourceConfig, bool) {
	for _, resource := range resources {
		if resource.Name == name {
			return resource, true
		}
	}

	return ResourceConfig{}, false
}

type JobConfigs []JobConfig

func (jobs JobConfigs) Lookup(name string) (JobConfig, bool) {
	for _, job := range jobs {
		if job.Name == name {
			return job, true
		}
	}

	return JobConfig{}, false
}

func (config Config) JobIsPublic(jobName string) (bool, error) {
	job, found := config.Jobs.Lookup(jobName)
	if !found {
		return false, fmt.Errorf("cannot find job with job name '%s'", jobName)
	}

	return job.Public, nil
}
