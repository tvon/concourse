package inputmapper

import (
	"code.cloudfoundry.org/lager"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/scheduler/inputmapper/inputconfig"
)

//go:generate counterfeiter . InputMapper

type InputMapper interface {
	SaveNextInputMapping(
		logger lager.Logger,
		versions *db.VersionsDB,
		job db.Job,
		resources db.Resources,
	) error
}

func NewInputMapper(transformer inputconfig.Transformer) InputMapper {
	return &inputMapper{transformer: transformer}
}

type inputMapper struct {
	transformer algorithm.Transformer
}

func (i *inputMapper) SaveNextInputMapping(
	logger lager.Logger,
	versions *db.VersionsDB,
	job db.Job,
	resources db.Resources,
) error {
	logger = logger.Session("save-next-input-mapping")

	mapping, ok, err := i.transformer.Transform(versions, job, resources)
	if err != nil {
		logger.Error("failed-to-compute-input-mappings", err)
		return err
	}

	err = job.SaveNextInputMapping(mapping, ok)
	if err != nil {
		logger.Error("failed-to-save-next-input-mapping", err)
		return err
	}

	return nil
}
