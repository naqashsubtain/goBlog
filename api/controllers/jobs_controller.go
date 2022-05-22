package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/naqash/goBlog/api/auth"
	"github.com/naqash/goBlog/api/models"
	"github.com/naqash/goBlog/api/responses"
	"github.com/naqash/goBlog/api/utils/formaterror"
)

type queryParams struct {
	Distance float64 `json:"distance"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

func (server *Server) CreateJob(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	Job := models.Job{}
	err = json.Unmarshal(body, &Job)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	Job.Prepare()
	err = Job.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != Job.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	JobCreated, err := Job.SaveJob(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, JobCreated.ID))
	responses.JSON(w, http.StatusCreated, JobCreated)
}

func (server *Server) GetJobs(w http.ResponseWriter, r *http.Request) {

	Job := models.Job{}

	Jobs, err := Job.FindAllJobs(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, Jobs)
}

func (server *Server) GetJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	Job := models.Job{}

	JobReceived, err := Job.FindJobByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, JobReceived)
}

func (server *Server) UpdateJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the Job id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized token not found"))
		return
	}

	// Check if the Job exist
	Job := models.Job{}
	err = server.DB.Debug().Model(models.Job{}).Where("id = ?", pid).Take(&Job).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("job not found"))
		return
	}

	// If a user attempt to update a Job not belonging to him
	if uid != Job.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized job is not for user"))
		return
	}
	// Read the data Job
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	JobUpdate := models.Job{}
	err = json.Unmarshal(body, &JobUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != JobUpdate.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	JobUpdate.Prepare()
	err = JobUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	JobUpdate.ID = Job.ID //this is important to tell the model the Job id to update, the other update field are set above

	JobUpdated, err := JobUpdate.UpdateAJob(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, JobUpdated)
}

func (server *Server) DeleteJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid Job id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Job exist
	Job := models.Job{}
	err = server.DB.Debug().Model(models.Job{}).Where("id = ?", pid).Take(&Job).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this Job?
	if uid != Job.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = Job.DeleteAJob(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}

func (server *Server) GetJobsbyDistance(w http.ResponseWriter, r *http.Request) {

	Job := models.Job{}
	// vars := mux.Vars(r)
	params := queryParams{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	err = json.Unmarshal(body, &params)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fmt.Println("quer params: ")
	fmt.Print(params.Distance)
	fmt.Print(params.Lat)
	fmt.Println(params.Lon)
	Jobs, err := Job.FindJobsByDistance(server.DB, params.Distance, params.Lat, params.Lon)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, Jobs)
}
