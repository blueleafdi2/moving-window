package server

import "net/http"

type Service interface {
	CountHandler(w http.ResponseWriter, r *http.Request)
}

var service Service

func InjectService(impl Service) {
	service = impl
}

func RefService() Service {
	return service
}
