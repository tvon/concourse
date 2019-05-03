package algorithm

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/concourse/concourse/atc/db"
)

// NOTE: we're effectively ignoring check_order here and relying on
// build history - be careful when doing #413 that we don't go
// 'back in time'
//
// QUESTION: is version_md5 worth it? might it be surprising
// that all it takes is (resource name, version identifier)
// to consider something 'passed'?

var ErrLatestVersionNotFound = errors.New("latest version of resource not found")
var ErrVersionNotFound = errors.New("version of resource not found")

type PinnedVersionNotFoundError struct {
	PinnedVersionID int
}

func (e PinnedVersionNotFoundError) Error() string {
	return fmt.Sprintf("pinned version %d not found", e.PinnedVersionID)
}

type NoSatisfiableBuildsForPassedJobError struct {
	JobID int
}

func (e NoSatisfiableBuildsForPassedJobError) Error() string {
	return fmt.Sprintf("passed job %d does not have a build that satisfies the constraints", e.JobID)
}

type version struct {
	ID             int
	VouchedForBy   map[int]bool
	SourceBuildIds []int
	ResolveError   error
}

func newCandidateVersion(id int) *version {
	return &version{
		ID:             id,
		VouchedForBy:   map[int]bool{},
		SourceBuildIds: []int{},
		ResolveError:   nil,
	}
}

func newCandidateError(err error) *version {
	return &version{
		ID:             0,
		VouchedForBy:   map[int]bool{},
		SourceBuildIds: []int{},
		ResolveError:   err,
	}
}

func Resolve(db *db.VersionsDB, inputConfigs InputConfigs) ([]*version, error) {
	versions := make([]*version, len(inputConfigs))

	_, err := tryResolve(0, db, inputConfigs, versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func tryResolve(depth int, db *db.VersionsDB, inputConfigs InputConfigs, candidates []*version) (bool, error) {
	// NOTE: this is probably made most efficient by doing it in order of inputs
	// with jobs that have the broadest output sets, so that we can pin the most
	// at once
	//
	// NOTE 3: maybe also select distinct build outputs so we don't waste time on
	// the same thing (i.e. constantly re-triggered build)
	//
	// NOTE : make sure everything is deterministically ordered

	for i, inputConfig := range inputConfigs {
		debug := func(messages ...interface{}) {
			log.Println(
				append(
					[]interface{}{
						strings.Repeat("-", depth) + fmt.Sprintf("[%s]", inputConfig.Name),
					},
					messages...,
				)...,
			)
		}

		// coming from recursive call; already gave up on input
		if candidates[i] != nil && candidates[i].ResolveError != nil {
			continue
		}

		if len(inputConfig.Passed) == 0 {
			fmt.Println("PASSED")
			// coming from recursive call; already set to the latest version
			if candidates[i] != nil {
				fmt.Println("RECURSIVE ", candidates[i])
				continue
			}

			var versionID int
			if inputConfig.PinnedVersionID != 0 {
				fmt.Println("PPPPPPPPPPP")
				// pinned
				exists, err := db.FindVersionOfResource(inputConfig.PinnedVersionID)
				if err != nil {
					return false, err
				}

				if !exists {
					candidates[i] = newCandidateError(PinnedVersionNotFoundError{inputConfig.PinnedVersionID})
					return false, nil
				}

				versionID = inputConfig.PinnedVersionID
				debug("setting candidate", i, "to unconstrained version", versionID)
			} else if inputConfig.UseEveryVersion {
				fmt.Println("EEEEEEEEEEEEEE")
				buildID, found, err := db.LatestBuildID(inputConfig.JobID)
				if err != nil {
					return false, err
				}

				if found {
					versionID, found, err = db.NextEveryVersion(buildID, inputConfig.ResourceID)
					if err != nil {
						return false, err
					}

					if !found {
						candidates[i] = newCandidateError(ErrVersionNotFound)
						return false, nil
					}
				} else {
					versionID, found, err = db.LatestVersionOfResource(inputConfig.ResourceID)
					if err != nil {
						return false, err
					}

					if !found {
						candidates[i] = newCandidateError(ErrLatestVersionNotFound)
						return false, nil
					}
				}

				debug("setting candidate", i, "to version for version every", versionID)
			} else {
				fmt.Println("ELSEEEEEEEEEEEEEEEEEEEEEEEE")
				// there are no passed constraints, so just take the latest version
				var err error
				var found bool
				versionID, found, err = db.LatestVersionOfResource(inputConfig.ResourceID)
				if err != nil {
					return false, nil
				}

				if !found {
					fmt.Println("UUUUUUUUU")
					candidates[i] = newCandidateError(ErrLatestVersionNotFound)
					return false, nil
				}

				debug("setting candidate", i, "to version for latest", versionID)
			}

			candidates[i] = newCandidateVersion(versionID)
			continue
		}

		orderedJobs := []int{}
		if len(inputConfig.Passed) != 0 {
			var err error
			orderedJobs, err = db.OrderPassedJobs(inputConfig.JobID, inputConfig.Passed)
			if err != nil {
				return false, err
			}
		}

		fmt.Println("OOOOOOOOOOOOORDERRR: ", orderedJobs)
		for _, jobID := range orderedJobs {
			if candidates[i] != nil {
				debug(i, "has a candidate")

				// coming from recursive call; we've already got a candidate
				if candidates[i].VouchedForBy[jobID] {
					debug("job", jobID, i, "already vouched for", candidates[i].ID)
					// we've already been here; continue to the next job
					continue
				} else {
					debug("job", jobID, i, "has not vouched for", candidates[i].ID)
				}
			} else {
				debug(i, "has no candidate yet")
			}

			// loop over previous output sets, latest first
			var builds []int

			if inputConfig.UseEveryVersion {
				buildID, found, err := db.LatestBuildID(inputConfig.JobID)
				if err != nil {
					return false, err
				}

				if found {
					constraintBuildID, found, err := db.LatestConstraintBuildID(buildID, jobID)
					if err != nil {
						return false, err
					}

					if found {
						builds, err = db.UnusedBuilds(constraintBuildID, jobID)
						if err != nil {
							return false, err
						}
					}
				}
			}

			var err error
			if len(builds) == 0 {
				builds, err = db.SuccessfulBuilds(jobID)
				if err != nil {
					return false, err
				}
			}

			for _, buildID := range builds {
				outputs, err := db.BuildOutputs(buildID)
				if err != nil {
					return false, err
				}

				debug("job", jobID, "trying build", jobID, buildID)

				restore := map[int]*version{}

				var mismatch bool

				// loop over the resource versions that came out of this build set
			outputs:
				for _, output := range outputs {
					debug("build", buildID, "output", output.ResourceID, output.VersionID)

					// try to pin each candidate to the versions from this build
					for c, candidate := range candidates {
						if inputConfigs[c].ResourceID != output.ResourceID {
							// unrelated to this output
							continue
						}

						if !inputConfigs[c].Passed[jobID] {
							// this candidate is unaffected by the current job
							debug("independent", inputConfigs[c].Passed, jobID)
							continue
						}

						if db.DisabledVersionIDs[output.VersionID] {
							mismatch = true
							break outputs
						}

						if candidate != nil && candidate.ID != output.VersionID {
							// don't return here! just try the next output set. it's possible
							// we just need to use an older output set.
							debug("mismatch")
							mismatch = true
							break outputs
						}

						// if this doesn't work out, restore it to either nil or the
						// candidate *without* the job vouching for it
						if candidate == nil {
							restore[c] = nil

							debug("setting candidate", c, "to", output.VersionID)
							candidates[c] = newCandidateVersion(output.VersionID)
						}

						debug("job", jobID, "vouching for", output.ResourceID, "version", output.VersionID)
						candidates[c].VouchedForBy[jobID] = true
						candidates[c].SourceBuildIds = append(candidates[c].SourceBuildIds, buildID)
					}
				}

				// we found a candidate for ourselves and the rest are OK too - recurse
				if candidates[i] != nil && candidates[i].VouchedForBy[jobID] && !mismatch {
					debug("recursing")

					resolved, err := tryResolve(depth+1, db, inputConfigs, candidates)
					if err != nil {
						return false, err
					}

					if resolved {
						// we've attempted to resolve all of the inputs!
						return true, nil
					}
				}

				debug("restoring")

				for c, version := range restore {
					// either there was a mismatch or resolving didn't work; go on to the
					// next output set
					debug("restoring candidate", c, "to", version)
					candidates[c] = version
				}
			}

			// we've exhausted all the builds and never found a matching input set;
			// give up on this input
			candidates[i] = newCandidateError(NoSatisfiableBuildsForPassedJobError{jobID})
			return false, nil
		}
	}

	// go to the end of all the inputs
	return true, nil
}
