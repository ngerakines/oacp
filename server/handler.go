package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type handlers struct {
	logger  *zap.Logger
	storage storage
}

func (h handlers) health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func (h handlers) callback(c *gin.Context) {
	ctx := c.Request.Context()

	state := c.Query("state")
	code := c.Query("code")

	if state == "" {
		h.logger.Warn("incoming request without state")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if code == "" {
		h.logger.Warn("incoming request without code")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	location, err := h.storage.GetLocation(ctx, state)
	if err != nil {
		h.logger.
			Warn("unable to retrieve location",
				zap.String("state", state),
				zap.Error(err))
		if errors.Is(err, errLocationNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	u, err := url.Parse(location)
	if err != nil {
		h.logger.
			Warn("unable to parse location",
				zap.String("state", state),
				zap.String("location", location),
				zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	query := u.Query()
	query.Set("state", state)
	query.Set("code", code)
	u.RawQuery = query.Encode()

	c.Redirect(http.StatusFound, u.String())
}

func (h handlers) apiRecordLocation(c *gin.Context) {
	ctx := c.Request.Context()

	state := c.PostForm("state")
	if len(state) == 0 {
		h.logger.
			Warn("state not provided")
		c.String(http.StatusBadRequest, errStateInvalid.Error())
		return
	}

	location := c.PostForm("location")

	if len(location) == 0 {
		h.logger.
			Warn("location not provided")
		c.String(http.StatusBadRequest, errLocationInvalid.Error())
		return
	}

	_, err := url.Parse(location)
	if err != nil {
		h.logger.
			Warn("unable to parse location",
				zap.String("location", location),
				zap.Error(err))
		c.String(http.StatusBadRequest, errLocationInvalid.Error())
		return
	}

	err = h.storage.RecordLocation(ctx, state, location)
	if err != nil {
		h.logger.
			Warn("unable record location",
				zap.String("state", state),
				zap.String("location", location),
				zap.Error(err))
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
