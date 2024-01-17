package http

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/services/socket"
	"github.com/vsheoran/trends/utils"
	"net/http"
)

func SocketHandleFunc(w http.ResponseWriter, r *http.Request) {
	var err error
	var conn *websocket.Conn

	params := mux.Vars(r)
	sasSymbol := params[constants.SasSymbolKey]

	if r.Method == http.MethodGet {
		conn, err = (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)

		uid := uuid.New()

		client := socket.New(logger, conn, sasSymbol, uid.String(), svc.HubService)

		svc.HubService.Register <- client

		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, uid)
		}
	}

	if r.Method == http.MethodPost {
		var st contracts.Stock
		utils.Decode(r.Body, &st)

		err = svc.HubService.UpdateStock(sasSymbol, st)

		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusNoContent)
		}
	}

	if r.Method == http.MethodDelete {
		var ucl *socket.UnregisterClient
		utils.Decode(r.Body, ucl)

		ucl.Ticker = sasSymbol

		svc.HubService.Unregister <- ucl

		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusOK)
	}

	if err != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, ErrorResponse{Error: err.Error()})
	}
}
