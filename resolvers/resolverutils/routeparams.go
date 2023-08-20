package resolverutils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/raphael-p/beango/utils/context"
	"github.com/raphael-p/beango/utils/logger"
)

const (
	USERNAME_KEY  = "username"
	CHAT_ID_KEY   = "chatID"
	CHAT_NAME_KEY = "chatName"
)

type RouteParams struct {
	Username string
	ChatID   int64
	ChatName string
}

// TODO: unit test, move some stuff that was meant to test GetRequestContext
func extractRouteParams(r *http.Request, paramKeys ...string) (*RouteParams, *HTTPError) {
	routeParams := new(RouteParams)
	for _, paramKey := range paramKeys {
		value, err := context.GetParam(r, paramKey)
		if err != nil {
			logger.Error(err.Error())
			return nil, &HTTPError{
				http.StatusInternalServerError,
				fmt.Sprint("failed to fetch path parameter: ", paramKey),
			}
		}

		switch paramKey {
		case USERNAME_KEY:
			routeParams.Username = value
		case CHAT_ID_KEY:
			chatID, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, &HTTPError{
					http.StatusBadRequest,
					"chat ID must be an integer",
				}
			}
			routeParams.ChatID = chatID
		case CHAT_NAME_KEY:
			routeParams.ChatName = value
		default:
			message := "invalid route param key: " + paramKey
			logger.Error(message)
			return nil, &HTTPError{http.StatusInternalServerError, message}
		}
	}
	return routeParams, nil
}
