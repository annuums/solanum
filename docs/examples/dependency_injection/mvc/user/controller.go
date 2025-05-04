package user

import (
	"github.com/annuums/solanum"
	"net/http"
)

func NewController() *solanum.SolaController {
	ctrl := solanum.NewController()

	ctrl.SetHandlers(
		solanum.SolaService{
			Uri:     "/users",
			Method:  http.MethodGet,
			Handler: retrieveUser,
		},
		solanum.SolaService{
			Uri:     "/users",
			Method:  http.MethodPost,
			Handler: addUser,
		},
	)

	return ctrl
}
