/*
 *   Copyright 2022 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package middlewares

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sapcc/go-api-declarations/cadf"
	"github.com/sapcc/go-bits/audittools"
	"go-micro.dev/v4/logger"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/config"
	"github.com/sapcc/andromeda/internal/policy"
)

type auditController struct {
	EventSink    chan<- cadf.Event
	observerUUID string
}

func NewAuditController() *auditController {
	s := make(chan cadf.Event, 20)
	q := auditController{
		EventSink:    s,
		observerUUID: audittools.GenerateUUID(),
	}
	rabbitmqQueueName := config.Global.Audit.QueueName
	transportURL, err := url.Parse(config.Global.Audit.TransportURL)
	if err != nil {
		panic(err)
	}

	go audittools.AuditTrail{
		EventSink: s,
		OnSuccessfulPublish: func() {
			logger.Debug("Notification sent")
		},
		OnFailedPublish: func() {
			logger.Debug("Notification failed")
		},
	}.Commit(*transportURL, rabbitmqQueueName)
	return &q
}

// auditResponseWriter is a wrapper of regular ResponseWriter
type auditResponseWriter struct {
	http.ResponseWriter
	controller *auditController
	request    *http.Request
}

// AuditResource is an audittools.EventRenderer.
type AuditResource struct {
	project     string
	resource    string
	routeParams middleware.RouteParams
}

// Render implements the audittools.EventRenderer interface.
func (a AuditResource) Render() cadf.Resource {
	id := ""
	attachments := []cadf.Attachment{}
	for _, routeParam := range a.routeParams {
		attachments = append(attachments, cadf.Attachment{
			Name:    routeParam.Name,
			Content: routeParam.Value,
		})
		// Last route param is our target id
		id = routeParam.Value
	}
	res := cadf.Resource{
		TypeURI:     fmt.Sprintf("gtm/%s", a.resource),
		ID:          id,
		ProjectID:   a.project,
		Attachments: attachments,
	}

	return res
}

func (arw *auditResponseWriter) WriteHeader(code int) {
	arw.ResponseWriter.WriteHeader(code)

	mr := middleware.MatchedRouteFrom(arw.request)
	resource := strings.Split(policy.RuleFromHTTPRequest(arw.request), ":")[1]
	user, err := auth.UserForRequest(arw.request)
	if err != nil {
		logger.Error(err)
		return
	}

	p := audittools.EventParameters{
		Time:       time.Now(),
		Request:    arw.request,
		User:       user,
		ReasonCode: code,
		Action:     cadf.GetAction(arw.request.Method),
		Target: AuditResource{
			user.ProjectScopeUUID(),
			resource,
			mr.Params,
		},
	}
	p.Observer.TypeURI = "service/gtm"
	p.Observer.Name = "Andromeda"
	p.Observer.ID = arw.controller.observerUUID
	arw.controller.EventSink <- audittools.NewEvent(p)
}

func (ac *auditController) NewAuditResponseWriter(w http.ResponseWriter, r *http.Request) *auditResponseWriter {
	return &auditResponseWriter{w, ac, r}
}

// AuditHandler provides the audit handling.
func (ac *auditController) AuditHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}

		qrw := ac.NewAuditResponseWriter(w, r)
		next.ServeHTTP(qrw, r)
	})
}
