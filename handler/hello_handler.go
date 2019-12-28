package handler

import (
	"github.com/ccchieh/ginHelper"
	"github.com/ccchieh/ginHelper_example/middleware"
	"github.com/ccchieh/ginHelper_example/service"
	"github.com/gin-gonic/gin"
)

func (h *Helper) HelloHandler() (r *ginHelper.Router) {
	return &ginHelper.Router{
		Param:  new(service.Hello),
		Path:   "/",
		Method: "GET",
	}
}

func (h *Helper) AdminHandler() (r *ginHelper.Router) {
	return &ginHelper.Router{
		Param:  new(service.Hello),
		Path:   "/admin",
		Method: "GET",
		Handlers: []gin.HandlerFunc{
			middleware.AdminMiddleware(),
			ginHelper.GenHandlerFunc,
		},
	}
}

func (h *Helper) UnAdminHandler() (r *ginHelper.Router) {
	return &ginHelper.Router{
		Param:  new(service.Hello),
		Path:   "/user",
		Method: "GET",
		Handlers: []gin.HandlerFunc{
			middleware.UnAdminMiddleware(),
		},
	}
}
