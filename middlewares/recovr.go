// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apex/log"

	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/models"
)

// Recovr is a http middleware that recovers from panics and logs (or prints) the error.
func Recovr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					err = fmt.Errorf("%v", rec)
				}
				log.WithError(err).Error("Recovered from HTTP panic")

				var errMsg string
				if config.Global.Default.Debug {
					errMsg = err.Error()
				} else {
					// In production, we do not expose the error message to the user
					errMsg = http.StatusText(http.StatusInternalServerError)
				}

				switch r.Header.Get("Accept") {
				case "application/json", "application/json; charset=utf-8", "*/*":
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusInternalServerError)
					if err := json.NewEncoder(w).Encode(models.Error{
						Code:    500,
						Message: errMsg,
					}); err != nil {
						http.Error(w, errMsg, http.StatusInternalServerError)
					}
				default:
					// For non-JSON requests, we return a plain text error message
					http.Error(w, errMsg, http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
