package transport

import (
	"context"
	"fmt"
	"github.com/vsheoran/trends/pkg/constants"
	"github.com/vsheoran/trends/pkg/contracts"
	"github.com/vsheoran/trends/templates/home"
	"github.com/vsheoran/trends/utils"
	"net/http"
)

func InitTicker(key string, svc Services, w http.ResponseWriter, r *http.Request) {
	_, err0 := svc.TickerService.Init(key, "")
	if err0 != nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(
			w,
			ErrorResponse{Error: fmt.Sprintf("failed to initilize ticker: %s", err0.Error())},
		)
		return
	}

	data := svc.TickerService.GetAllSummary()
	if data == nil {
		w.Header().Add(constants.HeaderContentTypeKey, constants.HeaderContentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, ErrorResponse{Error: "no ticker data found"})
		return
	}
	component := home.Home(contracts.HTMXData{SummaryMap: data})
	component.Render(context.Background(), w)
}
