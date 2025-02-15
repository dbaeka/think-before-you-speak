package http

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/swaggo/echo-swagger"
	"golang.org/x/sync/errgroup"
	"net/http"
	"think-before-you-speak/docs"
	"think-before-you-speak/pkg/config"
	"think-before-you-speak/pkg/handler"
	"think-before-you-speak/pkg/log"
	"think-before-you-speak/pkg/service"
)

type Server struct {
	engine *echo.Echo
	config *config.Config

	handlers []handler.Handler
}

func NewEcho(conf *config.Config) (*Server, error) {
	chatService := service.NewChatService(conf.Pricing)

	chatHandler := handler.NewChatHandler(chatService)
	systemHandler := handler.NewSystemHandler()

	handlers := []handler.Handler{chatHandler, systemHandler}

	e := echo.New()
	e.HideBanner = true

	useMiddlewares(e, conf)
	e.HTTPErrorHandler = HandleError

	s := Server{
		engine:   e,
		config:   conf,
		handlers: handlers,
	}

	return &s, nil
}

func HandleError(err error, c echo.Context) {
	log.FromContext(c.Request().Context()).Err(err).Msg("HTTP error")

	httpCode := http.StatusInternalServerError
	msg := any("Internal server error")

	httpErr := &echo.HTTPError{}
	if errors.As(err, &httpErr) {
		httpCode = httpErr.Code
		msg = httpErr.Message
	}

	jsonErr := c.JSON(
		httpCode,
		map[string]any{
			"error": msg,
		},
	)

	if jsonErr != nil {
		panic(err)
	}
}

func (s *Server) Run() error {
	l := log.Get()
	defer s.Close()

	s.initRouter()

	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	l.Info().Msgf("Start server on: %s", addr)

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		err := s.engine.Start(addr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		return s.engine.Shutdown(ctx)
	})

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Close() {
	l := log.Get()
	l.Info().Msg("Closing server")
}

func (s *Server) initRouter() {
	l := log.Get()
	root := s.engine

	// register non-resource routers
	root.GET("/swagger/*", echoSwagger.WrapHandler)

	api := root.Group(s.config.Server.Prefix)
	// port int to string
	port := fmt.Sprintf("%d", s.config.Server.Port)
	docs.SwaggerInfo.Host = s.config.Server.Host + ":" + port

	docs.SwaggerInfo.BasePath = s.config.Server.Prefix
	handlers := make([]string, 0, len(s.handlers))
	for _, router := range s.handlers {
		router.RegisterRoutes(api)
		handlers = append(handlers, router.Name())
	}
	l.Info().Msgf("server enabled handlers: %v", handlers)
}
