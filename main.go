package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"test_aws/infrastructure"
	"test_aws/interfaces"
	"test_aws/usecases"
	"time"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		exitErrorf("Can not loca .env file")
	}

	remoteDBHandler := infrastructure.NewS3Handler(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_REGION"), os.Getenv("AWS_BUCKET"))
	localDBHandler := infrastructure.NewSqliteHandler(os.Getenv("LOCAL_DB_NAME"))
	storageHandler := infrastructure.NewStorageHandler(os.Getenv("LOCAL_FILE_PATH"))

	storageInteractor := new(usecases.StorageInterfactor)
	storageInteractor.Repository = interfaces.NewStorageRepo(storageHandler)

	remoteInteractor := new(usecases.RemoteInteractor)
	remoteInteractor.Repository = interfaces.NewRemoteFileRepo(remoteDBHandler)

	localInteractor := new(usecases.LocalInteractor)
	localInteractor.Repository = interfaces.NewLocalFileRepo(localDBHandler)

	statInteractor := new(usecases.StatInteractor)
	statInteractor.Repository = interfaces.NewLocalStatRepo(localDBHandler)

	filesInteractor := new(usecases.FileInteractor)
	filesInteractor.Local = localInteractor
	filesInteractor.Storage = storageInteractor
	filesInteractor.Remote = remoteInteractor

	feedHandler := new(interfaces.FeedHandler)
	feedHandler.FilesInteractor = filesInteractor

	statHandler := new(interfaces.StatsHandler)
	statHandler.StatsInteractor = statInteractor

	filesHandler := new(interfaces.FilesHandler)
	filesHandler.FilesInteractor = filesInteractor
	filesHandler.StatsInteractor = statInteractor

	r := mux.NewRouter()

	tick, _ := strconv.Atoi(os.Getenv("AWS_FEED_TICK"))
	feedHandler.Feed(time.Duration(tick))

	//r.HandleFunc("/feed", func(res http.ResponseWriter, req *http.Request) {
	//	tick, _ := strconv.Atoi(os.Getenv("AWS_FEED_TICK"))
	//	feedHandler.Feed(time.Duration(tick))
	//})

	r.HandleFunc("/files/{id}", func(res http.ResponseWriter, req *http.Request) {
		filesHandler.Files(res, req)
	})

	r.HandleFunc("/stats", func(res http.ResponseWriter, req *http.Request) {
		statHandler.Stats(res, req)
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}
