package abortbuild

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"

	"github.com/concourse/atc/config"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/web/routes"
)

type handler struct {
	logger lager.Logger

	jobs       config.Jobs
	db         db.DB
	httpClient *http.Client
}

func NewHandler(logger lager.Logger, jobs config.Jobs, db db.DB) http.Handler {
	return &handler{
		logger: logger,

		jobs: jobs,
		db:   db,

		httpClient: &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: 5 * time.Minute,
			},
		},
	}
}

func (handler *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buildID, err := strconv.Atoi(r.FormValue(":build_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log := handler.logger.Session("abort", lager.Data{
		"build": buildID,
	})

	abortURL, err := handler.db.AbortBuild(buildID)
	if err != nil {
		log.Error("failed-to-set-aborted", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if abortURL != "" {
		resp, err := handler.httpClient.Post(abortURL, "", nil)
		if err != nil {
			log.Error("failed-to-abort-build", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp.Body.Close()
	}

	redirectPath, err := routes.Routes.CreatePathForRoute(routes.GetBuild, rata.Params{
		"build_id": fmt.Sprintf("%d", buildID),
	})
	if err != nil {
		log.Fatal("failed-to-create-redirect-uri", err)
	}

	http.Redirect(w, r, redirectPath, 302)
}
