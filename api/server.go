package api

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func InitializeServer() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/repositories", ListRepositoriesHandler).Methods(http.MethodGet)
	router.HandleFunc("/repositories", CreateRepositoryHandler).Methods(http.MethodPost)
	router.HandleFunc("/repositories/{repository}", GetRepositoryInfoHandler).Methods(http.MethodGet)
	router.HandleFunc("/repositories/{repository}", DeleteRepositoryHandler).Methods(http.MethodDelete)
	router.HandleFunc("/repositories/{repository}/branches", ListBranchesHandler).Methods(http.MethodGet)
	router.HandleFunc("/repositories/{repository}/branches/{branch}", GetBranchHandler).Methods(http.MethodGet)
	router.HandleFunc("/repositories/{repository}/commits/{commit}", GetCommitHandler).Methods(http.MethodGet)
	router.HandleFunc("/repositories/{repository}/tree/{tree}", GetTreeHandler).Methods(http.MethodGet)
	router.HandleFunc("/repositories/{repository}/blobs/{blob}", GetBlobHandler).Methods(http.MethodGet)
	router.HandleFunc("/repositories/{repository}/tags", ListTagsHandler).Methods(http.MethodGet)
	router.HandleFunc("/repositories/{repository}/tags/{tag}", GetTagHandler).Methods(http.MethodGet)
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	err := srv.Shutdown(ctx)

	if err != nil {
		log.Println("unable to shutting down the server")
	}

	log.Println("shutting down")
	os.Exit(0)
}
