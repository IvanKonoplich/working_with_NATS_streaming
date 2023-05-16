package httpServerController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *HttpServerController) HandleGetById(c *gin.Context) {

	var inputID string
	inputID = c.Param("id")
	order, err := s.uc.Get(inputID)
	if err != nil {
		if err.Error() == fmt.Sprintf("order: %s does not exist", inputID) {
			NewResponseMessage(c, http.StatusBadRequest, err.Error())
		} else {
			NewResponseMessage(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		result := map[string]string{
			"OrderUid":          "OrderUid: " + order.OrderUid,
			"TrackNumber":       "TrackNumber: " + order.TrackNumber,
			"Entry":             "Entry: " + order.Entry,
			"Delivery":          "Delivery: " + fmt.Sprint(order.Delivery),
			"Payment":           "Payment: " + fmt.Sprint(order.Payment),
			"Items":             "Items: " + fmt.Sprint(order.Items),
			"Locale":            "Locale: " + order.Locale,
			"InternalSignature": "InternalSignature: " + order.InternalSignature,
			"CustomerId":        "CustomerId: " + order.CustomerId,
			"DeliveryService":   "DeliveryService: " + order.DeliveryService,
			"Shardkey":          "Shardkey: " + order.Shardkey,
			"SmId":              "SmId: " + fmt.Sprint(order.SmId),
			"DateCreated":       "DateCreated: " + fmt.Sprint(order.DateCreated),
			"OofShard":          "OofShard: " + order.OofShard,
		}
		c.HTML(http.StatusOK, "index2.tmpl", result)

	}
}
