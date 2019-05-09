package inputconfig

import (
	"errors"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/scheduler/algorithm"
)

//go:generate counterfeiter . InputConfigs

type InputConfigs interface {
	ComputeNextInputs(versionsDB *db.VersionsDB) (db.InputMapping, bool, error)
}

//go:generate counterfeiter . Transformer

type Transformer interface {
	Transform(db *db.VersionsDB, job db.Job, resources db.Resources) (db.InputMapping, bool, error)
}

func NewTransformer() Transformer {
	return &transformer{}
}

type transformer struct{}

func (i *transformer) Transform(
	versions *db.VersionsDB,
	job db.Job,
	resources db.Resources,
) (db.InputMapping, bool, error) {

	inputs := job.Config().Inputs()

	inputConfigs := algorithm.InputConfigs{}
	for i, input := range inputs {
		if input.Version == nil {
			input.Version = &atc.VersionConfig{Latest: true}
		}

		pinnedVersionID := 0
		if input.Version.Pinned != nil {
			resource, found := resources.Lookup(input.Resource)

			if !found {
				continue
			}

			if input.Version != nil && input.Version.Pinned == nil {
				if resource.CurrentPinnedVersion() != nil {
					inputs[i].Version = &atc.VersionConfig{Pinned: resource.CurrentPinnedVersion()}
				}
			}

			id, found, err := resource.ResourceConfigVersionID(input.Version.Pinned)
			if err != nil {
				return nil, err
			}

			if !found {
				continue
			}

			pinnedVersionID = id
		}

		jobs := db.JobSet{}
		for _, passedJobName := range input.Passed {
			jobs[versions.JobIDs[passedJobName]] = true
		}

		inputConfigs = append(inputConfigs, algorithm.InputConfig{
			Name:            input.Name,
			UseEveryVersion: input.Version.Every,
			PinnedVersionID: pinnedVersionID,
			ResourceID:      versions.ResourceIDs[input.Resource],
			Passed:          jobs,
			JobID:           job.ID(),
		})
	}

	mapping := db.InputMapping{}
	versions, err := Resolve(versionsDB, configs)
	if err != nil {
		return nil, false, err
	}

	valid := true
	for i, config := range configs {
		if versions[i] == nil {
			mapping[config.Name] = db.InputResult{
				ResolveError: errors.New("did not finish due to other resource errors"),
			}

			valid = false
		} else if versions[i].ResolveError != nil {
			mapping[config.Name] = db.InputResult{
				ResolveError: versions[i].ResolveError,
			}

			valid = false
		} else {
			firstOccurrence, err := versionsDB.IsVersionFirstOccurrence(versions[i].ID, config.JobID, config.Name)
			if err != nil {
				return nil, false, err
			}

			mapping[config.Name] = db.InputResult{
				Input: db.AlgorithmInput{
					AlgorithmVersion: db.AlgorithmVersion{
						ResourceID: config.ResourceID,
						VersionID:  versions[i].ID,
					},
					FirstOccurrence: firstOccurrence,
				},
				PassedBuildIDs: versions[i].SourceBuildIds,
			}
		}
	}

	return mapping, valid, nil
}
