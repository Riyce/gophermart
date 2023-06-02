package internal

import (
	"fmt"
	"github.com/riyce/gophermart/internal/accrual"
	"github.com/riyce/gophermart/internal/db"
	"github.com/riyce/gophermart/internal/handlers"
	"github.com/riyce/gophermart/internal/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(conf Config) error {
	runLogFile, _ := os.OpenFile(
		conf.LogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)

	var level zerolog.Level
	if conf.Debug {
		level = zerolog.InfoLevel
	} else {
		level = zerolog.WarnLevel
	}
	log.Logger = zerolog.New(multi).With().Timestamp().Logger().Level(level)

	conn, dbErr := db.GetConnection(conf.DBURI)
	if dbErr != nil {
		log.Fatal().Msg("error on create connection")
	} else {
		log.Info().Msg("connection created")
	}

	initErr := db.InitDB(conn)
	if dbErr != nil {
		log.Fatal().Err(initErr).Msg("error on init db")
	} else {
		log.Info().Msg("tables created")
	}

	database := db.NewDB(conn)

	service := services.NewService(database, conf.Key)

	handler := handlers.NewHandler(service)

	orderDB := accrual.NewDBOrderUpdater(conn)

	daemon := accrual.NewClientDaemon(conf.AccrualAddress, orderDB)
	go daemon.RunDaemon()
	log.Info().Msg("daemon started")

	s.httpServer = &http.Server{
		Addr:    conf.Address,
		Handler: handler.InitRoutes(conf.Debug),
	}

	fmt.Printf("Server started at: http://%s !\n", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}
