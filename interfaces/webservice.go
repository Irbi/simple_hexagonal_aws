package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"test_aws/usecases"
)

type FilesInteractor interface {
	GetFile(hash string, version int) (file usecases.File, err error)
}

type StatsInteractor interface {
	Stat() (stat usecases.Stat, err error)
	AddValid() error
	AddOut() error
}

type FilesHandler struct {
	FilesInteractor FilesInteractor
	StatsInteractor StatsInteractor
	HashSalt string
}

type StatsHandler struct {
	StatsInteractor StatsInteractor
}

func (handler *FilesHandler) Files(response http.ResponseWriter, request *http.Request) {

	fID := mux.Vars(request)["id"]
	fVer, _ := strconv.Atoi(request.FormValue("version"))
	file, err := handler.FilesInteractor.GetFile(fID, fVer)

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)

	} else if file.Content == nil {
		http.Error(response, "", http.StatusNotFound)

	} else {

		if file.Version != fVer {
			handler.StatsInteractor.AddOut()

			response.WriteHeader(http.StatusOK)
			headerCnt := fmt.Sprintf("attachment; filename=\"%s\"", file.Name)
			response.Header().Set("Content-Disposition", headerCnt)
			response.Write(file.Content)

		} else {
			handler.StatsInteractor.AddValid()

			response.WriteHeader(http.StatusNotModified)
		}
	}
}

func (handler *StatsHandler) Stats(response http.ResponseWriter, request *http.Request) {

	stats, err := handler.StatsInteractor.Stat()

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	} else {
		response.WriteHeader(http.StatusOK)
		response.Header().Set("Content-Type", "application/json")
		json.NewEncoder(response).Encode(stats)
	}
}
