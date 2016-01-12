package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"gopkg.in/jmcvetta/napping.v3"

	"github.com/Graylog2/nxlog-sidecar/context"
	"github.com/Graylog2/nxlog-sidecar/util"
)

type ResponseCollectorConfiguration struct {
	Inputs  []ResponseCollectorInput  `json:"inputs"`
	Outputs []ResponseCollectorOutput `json:"outputs"`
}

type ResponseCollectorInput struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Properties map[string]string `json:"properties"`
	ForwardTo  string            `json:"forward_to"`
}

type ResponseCollectorOutput struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Properties map[string]string `json:"properties"`
}

type RegistrationRequest struct {
	NodeId      string            `json:"node_id"`
	NodeDetails map[string]string `json:"node_details"`
}

func RequestConfiguration(context *context.Ctx) (ResponseCollectorConfiguration, error) {
	s := napping.Session{}
	url := context.ServerUrl.String() + "/plugins/org.graylog.plugins.collector/" + context.CollectorId
	res := ResponseCollectorConfiguration{}

	resp, err := s.Get(url, nil, &res, nil)
	if err == nil && resp.Status() != 200 {
		logrus.Error("Bad response status from Graylog server: ", resp.Status(), err)
	}

	return res, err
}

func UpdateRegistration(context *context.Ctx) {
	s := napping.Session{}

	registration := RegistrationRequest{}
	registration.NodeId = context.NodeId
	registration.NodeDetails = make(map[string]string)
	registration.NodeDetails["operating_system"] = util.GetSystemName()

	h := http.Header{}
	h.Add("User-Agent", "Graylog Collector v" + util.CollectorVersion)
	h.Add("X-Graylog-Collector-Version", util.CollectorVersion)

	r := napping.Request{
		Url:     context.ServerUrl.String() + "/system/collectors/" + context.CollectorId,
		Method:  "PUT",
		Payload: registration,
		Header:  &h,
	}

	resp, err := s.Send(&r)
	if err == nil && resp.Status() != 202 {
		logrus.Error("Bad response from Graylog server: ", resp.Status(), err)
	}
}