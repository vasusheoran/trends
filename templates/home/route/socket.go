package route

import (
	"github.com/vsheoran/trends/pkg/transport"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/services/socket"
	"github.com/vsheoran/trends/utils"
)

func SocketHandleFunc(w http.ResponseWriter, r *http.Request) {
	logger.Log("msg", "SocketHandleFunc")

	var err error

	params := mux.Vars(r)
	sasSymbol := params[constants.SasSymbolKey]

	if r.Method == http.MethodGet {
		conn := &websocket.Conn{}
		conn, err = (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(
			w,
			r,
			nil,
		)

		uid := uuid.New()

		client := socket.New(logger, conn, sasSymbol, uid.String(), svc.HubService)

		svc.HubService.Register <- client

		if err == nil {
			w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
			w.WriteHeader(http.StatusOK)
			utils.Encode(w, uid)
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
		utils.Encode(w, transport.ErrorResponse{Error: err.Error()})
	}
}
