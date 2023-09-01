package server

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/saeede-bellefille/quota-manager/pkg/config"
	"github.com/saeede-bellefille/quota-manager/pkg/quota"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type Job struct {
	Id     uint   `json:"id"`
	UserId uint   `json:"userId"`
	Data   string `json:"data"`
}

type Server struct {
	jobChan  chan Job
	echo     *echo.Echo
	quota    *quota.Quota
	address  string
	poolSize int
}

func New(conf *config.Config) *Server {
	return &Server{
		jobChan:  make(chan Job, conf.QueueSize),
		quota:    quota.New(conf.Redis, conf.UserQuota),
		echo:     echo.New(),
		address:  conf.ListenAddress,
		poolSize: conf.WorkerPoolSize,
	}
}

func (s *Server) Start() error {
	for i := 0; i < s.poolSize; i++ {
		go s.worker()
	}
	s.echo.POST("/run", s.run)
	return s.echo.Start(s.address)
}

func (s *Server) run(c echo.Context) error {
	var input Job
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, Response{Message: "Bad request"})
	}

	if err := s.quota.Check(input.Id, input.UserId, int64(len(input.Data))); err != nil {
		if e, ok := err.(quota.Error); ok {
			return c.JSON(http.StatusPreconditionFailed, Response{Message: e.Error()})
		}
		return c.JSON(http.StatusInternalServerError, Response{Message: "Internal error"})
	}

	s.jobChan <- input

	return c.JSON(http.StatusOK, Response{Success: true})
}

func (s *Server) worker() {
	for job := range s.jobChan {
		s.save(job)
	}
}

func (s *Server) save(j Job) {
	// This function does saving job
	time.Sleep(time.Second)
	log.Print(j.Data)
}
