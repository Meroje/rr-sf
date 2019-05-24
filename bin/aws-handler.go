package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fujiwara/ridge"
	"github.com/otiai10/copy"
	"github.com/spiral/roadrunner"
	"github.com/spiral/roadrunner/service/http"
)

var srv *roadrunner.Server

func init() {
	opcachePath := os.TempDir() + "/opcache"
	if err := os.MkdirAll(opcachePath, os.ModePerm); err != nil {
		panic(err)
	}

	err := copy.Copy(os.Getenv("LAMBDA_TASK_ROOT")+"/var/cache/prod", os.TempDir()+"/cache")
	if err != nil {
		panic(err)
	}

	os.Setenv("PATH", os.Getenv("PATH")+":"+os.Getenv("LAMBDA_TASK_ROOT")+"/bin")
	os.Setenv("LD_LIBRARY_PATH", os.Getenv("LD_LIBRARY_PATH")+":"+os.Getenv("LAMBDA_TASK_ROOT")+"/bin/lib")

	srv = roadrunner.NewServer(
		&roadrunner.ServerConfig{
			Command: "php -c=config/php.ini -d opcache.file_cache=" + opcachePath + " public/worker.php",
			Relay:   "pipes",
			Pool: &roadrunner.Config{
				NumWorkers:      1,
				MaxJobs:         100,
				AllocateTimeout: time.Second,
				DestroyTimeout:  time.Second,
			},
		},
	)
}

func main() {
	if err := srv.Start(); err != nil {
		panic(err)
	}
	defer srv.Stop()

	lambda.Start(handle)
}

func handle(event json.RawMessage) (interface{}, error) {
	r, err := ridge.NewRequest(event)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(r, &http.UploadsConfig{
		Dir:    os.TempDir(),
		Forbid: []string{".php", ".exe", ".bat"},
	})
	if err != nil {
		return nil, err
	}

	payload, err := request.Payload()
	if err != nil {
		return nil, err
	}

	res, err := srv.Exec(payload)
	if err != nil {
		return nil, err
	}

	response, err := http.NewResponse(res)
	if err != nil {
		return nil, err
	}

	w := ridge.NewResponseWriter()

	err = response.Write(w)
	if err != nil {
		return nil, err
	}

	return w.Response(), nil
}
