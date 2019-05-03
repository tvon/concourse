package scheduler

import (
	"time"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/scheduler/algorithm"
	"github.com/concourse/concourse/atc/scheduler/inputmapper"
)

type Scheduler struct {
	Pipeline     db.Pipeline
	InputMapper  inputmapper.InputMapper
	BuildStarter BuildStarter
}

func (s *Scheduler) Schedule(
	logger lager.Logger,
	versions *algorithm.VersionsDB,
	job db.Job,
	resources db.Resources,
	resourceTypes atc.VersionedResourceTypes,
) (map[string]time.Duration, error) {
	jobSchedulingTime := map[string]time.Duration{}

	jStart := time.Now()
	inputMapping, err := s.InputMapper.SaveNextInputMapping(logger, versions, job, resources)
	if err != nil {
		return err
	}

	err := s.ensurePendingBuildExists(logger, versions, job, resources)
	jobSchedulingTime[job.Name()] = time.Since(jStart)

	if err != nil {
		return jobSchedulingTime, err
	}

	err = s.BuildStarter.TryStartPendingBuildsForJob(logger, job, resources, resourceTypes, nextPendingBuilds)
	jobSchedulingTime[job.Name()] = jobSchedulingTime[job.Name()] + time.Since(jStart)

	if err != nil {
		return jobSchedulingTime, err
	}

	return jobSchedulingTime, nil
}

func (s *Scheduler) ensurePendingBuildExists(
	logger lager.Logger,
	versions *algorithm.VersionsDB,
	job db.Job,
	resources db.Resources,
) error {
	inputMapping, found, err := job.GetFullNextBuildInputs()
	if err != nil {
		logger.Error("failed-to-fetch-next-build-inputs", err)
		return err
	}

	if !found {
		// XXX: better info log pls
		logger.Info("next-build-inputs-not-found")
		return nil
	}

	for _, inputConfig := range job.Config().Inputs() {
		inputSource, ok := inputMapping[inputConfig.Name]

		//trigger: true, and the version has not been used
		if ok && inputSource.InputVersion.FirstOccurrence && inputConfig.Trigger {
			err := job.EnsurePendingBuildExists()
			if err != nil {
				logger.Error("failed-to-ensure-pending-build-exists", err)
				return err
			}

			break
		}
	}

	return nil
}
