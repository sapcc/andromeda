// SPDX-FileCopyrightText: Copyright 2025 SAP SE or an SAP affiliate company
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"encoding/json"
	"time"

	"github.com/apex/log"
	"github.com/go-openapi/runtime/middleware"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/sapcc/andromeda/internal/auth"
	"github.com/sapcc/andromeda/internal/driver"
	"github.com/sapcc/andromeda/internal/utils"
	"github.com/sapcc/andromeda/models"
	"github.com/sapcc/andromeda/restapi/operations/administrative"
)

type ServiceController struct {
	CommonController
}

// GetServices GET /services
func (c ServiceController) GetServices(params administrative.GetServicesParams) middleware.Responder {
	if _, err := auth.Authenticate(params.HTTPRequest, nil); err != nil {
		return administrative.NewGetServicesDefault(403).WithPayload(utils.PolicyForbidden)
	}

	//goland:noinspection GoPreferNilSlice
	var responseServices = []*models.Service{}
	var response micro.Info

	replyTo := nats.NewInbox()
	sub, err := c.nc.SubscribeSync(replyTo)
	if err != nil {
		log.Fatal(err.Error())
	}
	_ = c.nc.Flush()

	// Send the request
	subject, _ := micro.ControlSubject(micro.InfoVerb, "", "")
	if err = c.nc.PublishRequest(subject, replyTo, nil); err != nil {
		log.Fatal(err.Error())
	}

	timeout := 100 * time.Millisecond
	start := time.Now()
	for time.Since(start) < timeout {
		var msg *nats.Msg
		msg, err = sub.NextMsg(1 * time.Second)
		if err != nil {
			break
		}

		if err = json.Unmarshal(msg.Data, &response); err != nil {
			log.Fatal(err.Error())
		}

		responseServices = append(responseServices, &models.Service{
			ID:       response.ID,
			Version:  response.Version,
			Type:     response.Name,
			Provider: driver.GetServiceType(response.Name),
			Metadata: response.Metadata,
			// Todo: add metadata support to stormRPC
		})
	}
	_ = sub.Unsubscribe()

	return administrative.NewGetServicesOK().WithPayload(&administrative.GetServicesOKBody{Services: responseServices})
}
