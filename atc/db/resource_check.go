package db

import (
	"time"

	"github.com/concourse/concourse/atc"
)

//go:generate counterfeiter . ResourceCheck

type ResourceCheck interface {
	ID() int
	Resource() (Resource, error)
	Start() error
	Timeout() time.Duration
	FromVersion() atc.Version
	Finish() error
	FinishWithError(message string) error
}

const (
	CheckTypeResource     = "resource"
	CheckTypeResourceType = "resource_type"
)
