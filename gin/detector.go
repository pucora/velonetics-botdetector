package gin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	botdetector "github.com/velonetics/velonetics-botdetector/v2"
	velonetics "github.com/velonetics/velonetics-botdetector/v2/velonetics"
	"github.com/velonetics/lura/v2/config"
	"github.com/velonetics/lura/v2/logging"
	"github.com/velonetics/lura/v2/proxy"
	veloneticsgin "github.com/velonetics/lura/v2/router/gin"
)

const logPrefix = "[SERVICE: Gin][Botdetector]"

// Register checks the configuration and, if required, registers a bot detector middleware at the gin engine
func Register(cfg config.ServiceConfig, l logging.Logger, engine *gin.Engine) {
	detectorCfg, err := velonetics.ParseConfig(cfg.ExtraConfig)
	if err == velonetics.ErrNoConfig {
		return
	}
	if err != nil {
		l.Warning(logPrefix, err.Error())
		return
	}
	d, err := botdetector.New(detectorCfg)
	if err != nil {
		l.Warning(logPrefix, "Unable to create the bot detector:", err.Error())
		return
	}

	l.Debug(logPrefix, "The bot detector has been registered successfully")
	engine.Use(middleware(d, l))
}

// New checks the configuration and, if required, wraps the handler factory with a bot detector middleware
func New(hf veloneticsgin.HandlerFactory, l logging.Logger) veloneticsgin.HandlerFactory {
	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		next := hf(cfg, p)
		logPrefix := "[ENDPOINT: " + cfg.Endpoint + "][Botdetector]"

		detectorCfg, err := velonetics.ParseConfig(cfg.ExtraConfig)
		if err == velonetics.ErrNoConfig {
			return next
		}
		if err != nil {
			l.Warning(logPrefix, err.Error())
			return next
		}

		d, err := botdetector.New(detectorCfg)
		if err != nil {
			l.Warning(logPrefix, "Unable to create the bot detector:", err.Error())
			return next
		}

		l.Debug(logPrefix, "The bot detector has been registered successfully")
		return handler(d, next, l)
	}
}

func middleware(f botdetector.DetectorFunc, l logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if f(c.Request) {
			l.Error(logPrefix, errBotRejected)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func handler(f botdetector.DetectorFunc, next gin.HandlerFunc, l logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if f(c.Request) {
			l.Error(logPrefix, errBotRejected)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		next(c)
	}
}

var errBotRejected = errors.New("bot rejected")
