package aws

import (
	"github.com/GeertJohan/go.rice"
	"github.com/kihamo/shadow/service/frontend"
)

func (s *AwsService) GetTemplateBox() *rice.Box {
	return rice.MustFindBox("../aws/templates")
}

func (s *AwsService) GetFrontendMenu() *frontend.FrontendMenu {
	return &frontend.FrontendMenu{
		Name: "Aws",
		Url:  "/aws",
	}
}

func (s *AwsService) SetFrontendHandlers(router *frontend.Router) {
	router.GET(s, "/aws", &IndexHandler{})
}
