package routerx

import (
	"log"

	"github.com/gin-gonic/gin"
)

type APIGroup struct {
	RelativePath string
	Middlewares  []gin.HandlerFunc
	APIs         []*API
	Groups       []*APIGroup
}

type API struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

func NewGroup(relativePath string) *APIGroup {
	return &APIGroup{
		RelativePath: relativePath,
		APIs:         make([]*API, 0),
		Groups:       make([]*APIGroup, 0),
	}
}

func (grp *APIGroup) AddMiddlewares(middlewares ...gin.HandlerFunc) *APIGroup {
	grp.Middlewares = append(grp.Middlewares, middlewares...)
	return grp
}

func (grp *APIGroup) AddAPI(apis ...*API) *APIGroup {
	grp.APIs = append(grp.APIs, apis...)
	return grp
}

func (grp *APIGroup) AddAPIGroup(groups ...*APIGroup) *APIGroup {
	grp.Groups = append(grp.Groups, groups...)
	return grp
}

func (grp *APIGroup) AddToEngine(r *gin.Engine) {
	grp.register(&r.RouterGroup)
}

func (grp *APIGroup) register(r *gin.RouterGroup) {
	rg := r.Group(grp.RelativePath, grp.Middlewares...)

	for _, api := range grp.APIs {
		rg.Handle(api.Method, api.Path, api.Handler)
	}

	for _, grp := range grp.Groups {
		grp.register(rg)
	}

	log.Printf("Register %s group to router path %s\n", rg.BasePath(), r.BasePath())
}
