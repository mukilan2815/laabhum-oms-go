package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Mukilan-T/laabhum-gateway-go/config"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger"
	"github.com/gorilla/mux"
)

func SetupRoutes(cfg *config.Config, logger *logger.Logger, omsClient *oms.Client) *mux.Router {
    router := mux.NewRouter()

    // SCLAP Routes
    router.HandleFunc("/oms/scalper/order", createScalperOrder(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/order/{parentID}/execute", executeAllChildTrades(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/order/{parentID}/{childID}/execute", executeSpecificChildTrade(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/order/{parentID}/ctc", ctcOrder(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/order/{parentID}/{childID}/ctc", ctcChildOrder(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/order/{orderType}/{parentID}/modify", modifyOrder(logger, omsClient)).Methods(http.MethodPatch)
    router.HandleFunc("/oms/scalper/order/{orderType}/{parentID}/{childID}/modify", modifyChildOrder(logger, omsClient)).Methods(http.MethodPatch)
    router.HandleFunc("/oms/scalper/exit/trade", exitAllTrades(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/trade/{parentID}/exit", exitAllChildTrades(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/trade/{parentID}/{childID}/exit", exitSpecificChildTrade(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/order/{parentID}/cancel", cancelAllChildOrders(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/order/{parentID}/{orderId}/cancel", cancelSpecificChildOrder(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/scalper/trades/{parentID}", getTrades(logger, omsClient)).Methods(http.MethodGet)
    router.HandleFunc("/oms/scalper/order/{parentID}", deleteOrder(logger, omsClient)).Methods(http.MethodDelete)

    // ORDER Routes
    router.HandleFunc("/oms/orders", getOrders(logger, omsClient)).Methods(http.MethodGet)
    router.HandleFunc("/oms/order", createOrder(logger, omsClient)).Methods(http.MethodPut)
    router.HandleFunc("/oms/order/execute", executeOrder(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/order/cancel", cancelOrder(logger, omsClient)).Methods(http.MethodDelete)

    // POSITION Routes
    router.HandleFunc("/oms/positions", getPositions(logger, omsClient)).Methods(http.MethodGet)
    router.HandleFunc("/oms/position/sync", syncPosition(logger, omsClient)).Methods(http.MethodGet)
    router.HandleFunc("/oms/position/convert", convertPosition(logger, omsClient)).Methods(http.MethodPut)
    router.HandleFunc("/oms/position/order", createPositionOrder(logger, omsClient)).Methods(http.MethodPost)
    router.HandleFunc("/oms/position/order", deletePositionOrder(logger, omsClient)).Methods(http.MethodDelete)

    return router
}

// SCLAP Handlers
func createScalperOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var order oms.Order
        if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
            logger.Errorf("Failed to decode scalper order: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        createdOrder, err := omsClient.CreateOrder(order)
        if err != nil {
            logger.Errorf("Failed to create scalper order: %v", err)
            http.Error(w, "Failed to create order", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(createdOrder)
    }
}

func executeAllChildTrades(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]

        err := omsClient.ExecuteAllChildTrades(parentID)
        if err != nil {
            logger.Errorf("Failed to execute child trades for parent ID %s: %v", parentID, err)
            http.Error(w, "Failed to execute child trades", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func executeSpecificChildTrade(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]
        childID := vars["childID"]

        err := omsClient.ExecuteSpecificChildTrade(parentID, childID)
        if err != nil {
            logger.Errorf("Failed to execute child trade %s for parent ID %s: %v", childID, parentID, err)
            http.Error(w, "Failed to execute specific child trade", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func ctcOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]

        err := omsClient.CTCOrder(parentID)
        if err != nil {
            logger.Errorf("Failed to CTC order for parent ID %s: %v", parentID, err)
            http.Error(w, "Failed to CTC order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func ctcChildOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]
        childID := vars["childID"]

        err := omsClient.CTCChildOrder(parentID, childID)
        if err != nil {
            logger.Errorf("Failed to CTC child order %s for parent ID %s: %v", childID, parentID, err)
            http.Error(w, "Failed to CTC child order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func modifyOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        orderType := vars["orderType"]
        parentID := vars["parentID"]

        var order oms.Order
        if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
            logger.Errorf("Failed to decode modify order request: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        err := omsClient.ModifyOrder(orderType, parentID, order)
        if err != nil {
            logger.Errorf("Failed to modify order for parent ID %s: %v", parentID, err)
            http.Error(w, "Failed to modify order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func modifyChildOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        orderType := vars["orderType"]
        parentID := vars["parentID"]

        var order oms.Order
        if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
            logger.Errorf("Failed to decode modify child order request: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        err := omsClient.ModifyChildOrder(orderType, parentID, order)
        if err != nil {
            logger.Errorf("Failed to modify child order for parent ID %s: %v", parentID, err)
            http.Error(w, "Failed to modify child order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func exitAllTrades(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := omsClient.ExitAllTrades()
        if err != nil {
            logger.Errorf("Failed to exit all trades: %v", err)
            http.Error(w, "Failed to exit all trades", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    }
}

func exitAllChildTrades(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]

        err := omsClient.ExitAllChildTrades(parentID)
        if err != nil {
            logger.Errorf("Failed to exit all child trades for parent ID %s: %v", parentID, err)
            http.Error(w, "Failed to exit all child trades", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func exitSpecificChildTrade(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]
        childID := vars["childID"]

        err := omsClient.ExitSpecificChildTrade(parentID, childID)
        if err != nil {
            logger.Errorf("Failed to exit specific child trade %s for parent ID %s: %v", childID, parentID, err)
            http.Error(w, "Failed to exit specific child trade", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func cancelAllChildOrders(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]

        err := omsClient.CancelAllChildOrders(parentID)
        if err != nil {
            logger.Errorf("Failed to cancel all child orders for parent ID %s: %v", parentID, err)
            http.Error(w, "Failed to cancel all child orders", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func cancelSpecificChildOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]
        orderId := vars["orderId"]

        err := omsClient.CancelSpecificChildOrder(parentID, orderId)
        if err != nil {
            logger.Errorf("Failed to cancel child order %s for parent ID %s: %v", orderId, parentID, err)
            http.Error(w, "Failed to cancel specific child order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func getTrades(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]

        trades, err := omsClient.GetTrades(parentID)
        if err != nil {
            logger.Errorf("Failed to get trades for parent ID %s: %v", parentID, err)
            http.Error(w, "Failed to retrieve trades", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(trades)
    }
}

func deleteOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        parentID := vars["parentID"]

        err := omsClient.DeleteOrder(parentID)
        if err != nil {
            logger.Errorf("Failed to delete order for parent ID %s: %v", parentID, err)
            http.Error(w, "Failed to delete order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

// ORDER Handlers
func getOrders(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        orders, err := omsClient.GetOrders()
        if err != nil {
            logger.Errorf("Failed to get orders: %v", err)
            http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(orders)
    }
}

func createOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var order oms.Order
        if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
            logger.Errorf("Failed to decode order: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        createdOrder, err := omsClient.CreateOrder(order)
        if err != nil {
            logger.Errorf("Failed to create order: %v", err)
            http.Error(w, "Failed to create order", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(createdOrder)
    }
}

func executeOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var orderID string
        if err := json.NewDecoder(r.Body).Decode(&orderID); err != nil {
            logger.Errorf("Failed to decode execute order request: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        err := omsClient.ExecuteOrder(orderID)
        if err != nil {
            logger.Errorf("Failed to execute order %s: %v", orderID, err)
            http.Error(w, "Failed to execute order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func cancelOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var orderID string
        if err := json.NewDecoder(r.Body).Decode(&orderID); err != nil {
            logger.Errorf("Failed to decode cancel order request: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        err := omsClient.CancelOrder(orderID)
        if err != nil {
            logger.Errorf("Failed to cancel order %s: %v", orderID, err)
            http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

// POSITION Handlers
func getPositions(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        positions, err := omsClient.GetPositions()
        if err != nil {
            logger.Errorf("Failed to get positions: %v", err)
            http.Error(w, "Failed to retrieve positions", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(positions)
    }
}

func syncPosition(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := omsClient.SyncPositions()
        if err != nil {
            logger.Errorf("Failed to sync positions: %v", err)
            http.Error(w, "Failed to sync positions", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}
func convertPosition(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var position oms.Position // Ensure Position is defined in the oms package
        if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
            logger.Errorf("Failed to decode convert position request: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Assuming ConvertPosition expects a string (like position.ID)
        err := omsClient.ConvertPosition(position.ID) // Change this to the appropriate field
        if err != nil {
            logger.Errorf("Failed to convert position: %v", err)
            http.Error(w, "Failed to convert position", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}
func createPositionOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var positionOrder oms.Order // Use Order if that's what you're working with
        if err := json.NewDecoder(r.Body).Decode(&positionOrder); err != nil {
            logger.Errorf("Failed to decode position order: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        err := omsClient.CreatePositionOrder(positionOrder) // Now it only takes one argument
        if err != nil {
            logger.Errorf("Failed to create position order: %v", err)
            http.Error(w, "Failed to create position order", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated) // No need to encode createdOrder since CreatePositionOrder doesn't return it
    }
}

func deletePositionOrder(logger *logger.Logger, omsClient *oms.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var orderID string
		if err := json.NewDecoder(r.Body).Decode(&orderID); err != nil {
			logger.Errorf("Failed to decode delete position order request: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := omsClient.DeletePositionOrder(orderID)
		if err != nil {
			logger.Errorf("Failed to delete position order %s: %v", orderID, err)
			http.Error(w, "Failed to delete position order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
