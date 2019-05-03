package algorithm

import (
	"errors"
	"fmt"

	"github.com/concourse/concourse/atc/db"
)

type InputConfigs []InputConfig

type InputConfig struct {
	Name            string
	JobName         string
	Passed          db.JobSet
	UseEveryVersion bool
	PinnedVersionID int
	ResourceID      int
	JobID           int
}

func (configs InputConfigs) ComputeNextInputs(versionsDB *db.VersionsDB) (db.InputMapping, bool, error) {
	mapping := db.InputMapping{}

	versions, err := Resolve(versionsDB, configs)
	if err != nil {
		return nil, false, err
	}

	fmt.Println("################", versions[0])
	fmt.Println("$$$$$$$$$$$$$$$$", configs)
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

			fmt.Println("####################", versions[i])
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

// func (configs InputConfigs) Resolve(db *VersionsDB) (InputMapping, bool, error) {
// 	jobs := JobSet{}
// 	inputCandidates := InputCandidates{}

// 	for _, inputConfig := range configs {
// 		var versionCandidates VersionCandidates
// 		var err error

// 		if len(inputConfig.Passed) == 0 {
// 			if inputConfig.PinnedVersionID != 0 {
// 				versionCandidates, err = db.FindVersionOfResource(inputConfig.ResourceID, inputConfig.PinnedVersionID)
// 			} else {
// 				versionCandidates, err = db.AllVersionsOfResource(inputConfig.ResourceID)
// 			}
// 		} else {
// 			jobs = jobs.Union(inputConfig.Passed)

// 			versionCandidates, err = db.VersionsOfResourcePassedJobs(
// 				inputConfig.ResourceID,
// 				inputConfig.Passed,
// 			)
// 		}
// 		if err != nil {
// 			return nil, false, err
// 		}

// 		if inputConfig.UseEveryVersion {
// 			versionCandidates, err = versionCandidates.ConsecutiveVersions(inputConfig.JobID, inputConfig.ResourceID)
// 			if err != nil {
// 				return nil, false, err
// 			}
// 		}

// 		inputCandidates = append(inputCandidates, InputVersionCandidates{
// 			Input:             inputConfig.Name,
// 			Passed:            inputConfig.Passed,
// 			UseEveryVersion:   inputConfig.UseEveryVersion,
// 			PinnedVersionID:   inputConfig.PinnedVersionID,
// 			VersionCandidates: versionCandidates,
// 		})
// 	}

// 	basicMapping, ok, err := inputCandidates.Reduce(0, jobs)
// 	if err != nil {
// 		return nil, false, err
// 	}

// 	if !ok {
// 		return nil, false, nil
// 	}

// 	mapping := InputMapping{}
// 	for _, inputConfig := range configs {
// 		inputName := inputConfig.Name
// 		inputVersionID := basicMapping[inputName]
// 		firstOccurrence, err := db.IsVersionFirstOccurrence(inputVersionID, inputConfig.JobID, inputName)
// 		if err != nil {
// 			return nil, false, err
// 		}

// 		mapping[inputName] = InputVersion{
// 			ResourceID:      inputConfig.ResourceID,
// 			VersionID:       inputVersionID,
// 			FirstOccurrence: firstOccurrence,
// 		}
// 	}

// 	return mapping, true, nil
// }
