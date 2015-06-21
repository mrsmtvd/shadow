package tasks

import (
	"github.com/GeertJohan/go.rice"
	"github.com/kihamo/shadow/service/frontend"
)

func (s *TasksService) GetTemplateBox() *rice.Box {
	return rice.MustFindBox("../tasks/templates")
}

func (s *TasksService) GetFrontendMenu() *frontend.FrontendMenu {
	return &frontend.FrontendMenu{
		Name: "Tasks",
		Url:  "/tasks",
	}
}

func (s *TasksService) SetFrontendHandlers(router *frontend.Router) {
	router.GET(s, "/tasks/stats", &StatsHandler{})
	router.GET(s, "/tasks", &IndexHandler{})
}
