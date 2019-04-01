package atc

type JobConfig struct {
	Name   string `yaml:"name"`
	Public bool   `yaml:"public,omitempty"`

	DisableManualTrigger bool     `yaml:"disable_manual_trigger,omitempty"`
	Serial               bool     `yaml:"serial,omitempty"`
	Interruptible        bool     `yaml:"interruptible,omitempty"`
	SerialGroups         []string `yaml:"serial_groups,omitempty"`
	RawMaxInFlight       int      `yaml:"max_in_flight,omitempty"`
	BuildLogsToRetain    int      `yaml:"build_logs_to_retain,omitempty"`

	Abort   *PlanConfig `yaml:"on_abort,omitempty"`
	Failure *PlanConfig `yaml:"on_failure,omitempty"`
	Ensure  *PlanConfig `yaml:"ensure,omitempty"`
	Success *PlanConfig `yaml:"on_success,omitempty"`

	Plan PlanSequence `yaml:"plan"`
}

func (config JobConfig) Hooks() Hooks {
	return Hooks{
		Abort:   config.Abort,
		Failure: config.Failure,
		Ensure:  config.Ensure,
		Success: config.Success,
	}
}

func (config JobConfig) MaxInFlight() int {
	if config.Serial || len(config.SerialGroups) > 0 {
		return 1
	}

	if config.RawMaxInFlight != 0 {
		return config.RawMaxInFlight
	}

	return 0
}

func (config JobConfig) GetSerialGroups() []string {
	if len(config.SerialGroups) > 0 {
		return config.SerialGroups
	}

	if config.Serial || config.RawMaxInFlight > 0 {
		return []string{config.Name}
	}

	return []string{}
}

func (config JobConfig) Plans() []PlanConfig {
	plan := collectPlans(PlanConfig{
		Do:      &config.Plan,
		Abort:   config.Abort,
		Ensure:  config.Ensure,
		Failure: config.Failure,
		Success: config.Success,
	})

	return plan
}

func collectPlans(plan PlanConfig) []PlanConfig {
	var plans []PlanConfig

	if plan.Abort != nil {
		plans = append(plans, collectPlans(*plan.Abort)...)
	}

	if plan.Success != nil {
		plans = append(plans, collectPlans(*plan.Success)...)
	}

	if plan.Failure != nil {
		plans = append(plans, collectPlans(*plan.Failure)...)
	}

	if plan.Ensure != nil {
		plans = append(plans, collectPlans(*plan.Ensure)...)
	}

	if plan.Try != nil {
		plans = append(plans, collectPlans(*plan.Try)...)
	}

	if plan.Do != nil {
		for _, p := range *plan.Do {
			plans = append(plans, collectPlans(p)...)
		}
	}

	if plan.Aggregate != nil {
		for _, p := range *plan.Aggregate {
			plans = append(plans, collectPlans(p)...)
		}
	}

	return append(plans, plan)
}

func (config JobConfig) InputPlans() []PlanConfig {
	var inputs []PlanConfig

	for _, plan := range config.Plans() {
		if plan.Get != "" {
			inputs = append(inputs, plan)
		}
	}

	return inputs
}

func (config JobConfig) OutputPlans() []PlanConfig {
	var outputs []PlanConfig

	for _, plan := range config.Plans() {
		if plan.Put != "" {
			outputs = append(outputs, plan)
		}
	}

	return outputs
}

func (config JobConfig) Inputs() []JobInput {
	var inputs []JobInput

	for _, plan := range config.Plans() {
		if plan.Get != "" {
			get := plan.Get

			resource := get
			if plan.Resource != "" {
				resource = plan.Resource
			}

			inputs = append(inputs, JobInput{
				Name:     get,
				Resource: resource,
				Passed:   plan.Passed,
				Version:  plan.Version,
				Trigger:  plan.Trigger,
				Params:   plan.Params,
				Tags:     plan.Tags,
			})
		}
	}

	return inputs
}

func (config JobConfig) Outputs() []JobOutput {
	var outputs []JobOutput

	for _, plan := range config.Plans() {
		if plan.Put != "" {
			put := plan.Put

			resource := put
			if plan.Resource != "" {
				resource = plan.Resource
			}

			outputs = append(outputs, JobOutput{
				Name:     put,
				Resource: resource,
			})
		}
	}

	return outputs
}
