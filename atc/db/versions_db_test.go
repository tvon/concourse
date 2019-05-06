package db_test

import (
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/scheduler/algorithm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Versions DB", func() {
	Describe("OrderPassedJobs", func() {
		var (
			passedJobsPipeline db.Pipeline
			versionsDB         *db.VersionsDB
			jobIDs             []int
			currentJob         db.Job
			orderedJobs        []int
			passedJobs         db.JobSet
		)

		BeforeEach(func() {
			var err error
			passedJobsPipeline, _, err = defaultTeam.SavePipeline("passed-jobs-pipeline", atc.Config{
				Jobs: atc.JobConfigs{
					{
						Name: "current-job",
					},
					{
						Name: "passed-job-1",
					},
					{
						Name: "passed-job-2",
					},
					{
						Name: "passed-job-3",
					},
					{
						Name: "passed-job-4",
					},
					{
						Name: "passed-job-5",
					},
				},
			}, db.ConfigVersion(0), db.PipelineUnpaused)
			Expect(err).NotTo(HaveOccurred())

			var found bool
			currentJob, found, err = passedJobsPipeline.Job("current-job")
			Expect(err).ToNot(HaveOccurred())
			Expect(found).To(BeTrue())
		})

		JustBeforeEach(func() {
			versionsDB = &db.VersionsDB{
				Conn:   dbConn,
				JobIDs: jobIDs,
			}

			var err error
			orderedJobs, err = versionsDB.OrderPassedJobs(currentJob, passedJobs)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when there is a build for the current job", func() {
			BeforeEach(func() {
				build, err := currentJob.CreateBuild()
				Expect(err).ToNot(HaveOccurred())

				scheduled, err := build.Schedule([]db.BuildInput{})
				Expect(err).ToNot(HaveOccurred())
				Expect(scheduled).To(BeTrue())
			})

			Context("when all the passed jobs have build pipes", func() {
				var (
					passedJob1 db.Job
					passedJob2 db.Job
				)

				BeforeEach(func() {
					passedJob1, found, err = passedJobsPipeline.Job("passed-job-1")
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeTrue())

					passedJob2, found, err = passedJobsPipeline.Job("passed-job-2")
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeTrue())

					build1, err := passedJob1.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					build2, err := passedJob2.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					err = currentJob.SaveNextBuildPipes(db.InputMapping{
						db.InputResult{
							PassedBuildIDs: []int{build1.ID(), build2.ID()},
						},
					})
					Expect(err).ToNot(HaveOccurred())

					err = build.AdoptBuildPipes()
					Expect(err).ToNot(HaveOccurred())

					jobIDs = []int{job.ID(), passedJob1.ID(), passedJob2.ID()}

					passedJobs = algorithm.JobSet{passedJob1.ID(): {}, passedJob2.ID(): {}}
				})

				Context("when some passed jobs have the same number of builds", func() {
					It("should order by job id", func() {
						Expect(orderedJobs).To(Equal([]int{passedJob2.ID(), passedJob1.ID()}))
					})
				})

				Context("when the passed jobs have different number of builds", func() {
					BeforeEach(func() {
						build3, err := passedJob2.CreateBuild()
						Expect(err).ToNot(HaveOccurred())
					})

					It("should be ordered by build numbers", func() {
						Expect(orderedJobs).To(Equal([]int{passedJob1.ID(), passedJob2.ID()}))
					})
				})
			})

			Context("when some of the passed jobs have build pipes", func() {
				var (
					passedJob1 db.Job
					passedJob2 db.Job
					passedJob3 db.Job
					passedJob4 db.Job
					passedJob5 db.Job
				)

				BeforeEach(func() {
					passedJob1, found, err = passedJobsPipeline.Job("passed-job-1")
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeTrue())

					passedJob2, found, err = passedJobsPipeline.Job("passed-job-2")
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeTrue())

					passedJob3, found, err = passedJobsPipeline.Job("passed-job-3")
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeTrue())

					passedJob4, found, err = passedJobsPipeline.Job("passed-job-4")
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeTrue())

					passedJob5, found, err = passedJobsPipeline.Job("passed-job-5")
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeTrue())

					_, err = passedJob1.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					_, err = passedJob1.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					build3, err = passedJob2.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					_, err = passedJob4.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					build5, err := passedJob5.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					_, err = passedJob2.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					_, err = passedJob2.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					err = currentJob.SaveNextBuildPipes(db.InputMapping{
						db.InputResult{
							PassedBuildIDs: []int{build3.ID(), build5.ID()},
						},
					})
					Expect(err).ToNot(HaveOccurred())

					err = build.AdoptBuildPipes()
					Expect(err).ToNot(HaveOccurred())

					jobIDs = []int{job.ID(), passedJob1.ID(), passedJob2.ID(), passedJob3.ID(), passedJob4.ID(), passedJob5.ID()}

					passedJobs = algorithm.JobSet{passedJob1.ID(): {}, passedJob2.ID(): {}, passedJob3.ID(): {}, passedJob4.ID(): {}, passedJob5.ID(): {}}
				})

				It("should be ordered first by passed jobs that have build pipes and then by build numbers", func() {
					Expect(orderedJobs).To(Equal([]int{passedJob5.ID(), passedJob2.ID(), passedJob3.ID(), passedJob4.ID(), passedJob1.ID()}))
				})
			})

			Context("when none of the passed jobs have build pipes", func() {
				var passedJob1, passedJob2 db.Job

				BeforeEach(func() {
					passedJob1, found, err = passedJobsPipeline.Job("passed-job-1")
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeTrue())

					passedJob2, found, err = passedJobsPipeline.Job("passed-job-2")
					Expect(err).ToNot(HaveOccurred())
					Expect(found).To(BeTrue())

					_, err = passedJob1.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					_, err = passedJob1.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					_, err = passedJob2.CreateBuild()
					Expect(err).ToNot(HaveOccurred())

					passedJobs = algorithm.JobSet{passedJob1.ID(): {}, passedJob2.ID(): {}}
				})

				It("should be ordered by build numbers", func() {
					Expect(orderedJobs).To(Equal([]int{passedJob2.ID(), passedJob1.ID()}))
				})
			})
		})

		Context("when the current job has no builds", func() {
			var passedJob1, passedJob2 db.Job

			BeforeEach(func() {
				passedJob1, found, err = passedJobsPipeline.Job("passed-job-1")
				Expect(err).ToNot(HaveOccurred())
				Expect(found).To(BeTrue())

				passedJob2, found, err = passedJobsPipeline.Job("passed-job-2")
				Expect(err).ToNot(HaveOccurred())
				Expect(found).To(BeTrue())

				_, err = passedJob1.CreateBuild()
				Expect(err).ToNot(HaveOccurred())

				_, err = passedJob1.CreateBuild()
				Expect(err).ToNot(HaveOccurred())

				_, err = passedJob2.CreateBuild()
				Expect(err).ToNot(HaveOccurred())

				passedJobs = algorithm.JobSet{passedJob1.ID(): {}, passedJob2.ID(): {}}
			})

			It("should be ordered by build numbers", func() {
				Expect(orderedJobs).To(Equal([]int{passedJob2.ID(), passedJob1.ID()}))
			})
		})
	})
})
