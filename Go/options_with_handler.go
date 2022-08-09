package main

import (
	"fmt"

	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/valyala/fasthttp"
)

// OptionsWithHandler is a struct with a mux and shared credentials
type OptionsWithHandler struct {
	Options *options
	Client  *twilio.RestClient
}

// NewMOptionsWithHandler returns a OptionsWithHandler for http requests
// with shared credentials
func NewMOptionsWithHandler(o *options) OptionsWithHandler {
	return OptionsWithHandler{
		o,
		twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: o.AccountSid,
			Password: o.AuthToken,
		}),
	}
}

// HandleFastHTTP is the router function
func (m OptionsWithHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	} else {
		if string(ctx.Request.Header.Peek("Content-Type")) != "application/json" {
			ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		} else {
			switch string(ctx.Path()) {
			case "/":
				m.ping(ctx)
			case "/sms":
				m.smsRequest(ctx)
			case "/call":
				m.callRequest(ctx)
			case "/callandsms":
				m.callRequest(ctx)
				m.smsRequest(ctx)
			default:
				ctx.Error("Not found", fasthttp.StatusNotFound)
			}
		}
	}
}

func (m OptionsWithHandler) ping(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "ping")
}

func (m OptionsWithHandler) smsRequest(ctx *fasthttp.RequestCtx) {
	receiver := m.findReciver(ctx)
	if receiver == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Error("Bad request: receiver not specified")
		return
	}

	status, message, err := m.findMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Error("Bad request: Error parsing json")
		return
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(receiver)
	params.SetFrom(m.Options.Sender)
	params.SetBody(status + ": " + message)

	resp, err := m.Client.Api.CreateMessage(params)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Error("Bad request: unable to send SMS to this number")
		return
	}

	fmt.Println("Message Sid: " + *resp.Sid)
}

func (m OptionsWithHandler) callRequest(ctx *fasthttp.RequestCtx) {
	receiver := m.findReciver(ctx)
	if receiver == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Error("Bad request: receiver not specified")
		return
	}

	_, message, err := m.findMessage(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Error("Bad request: Error parsing json")
		return
	}

	params := &openapi.CreateCallParams{}
	params.SetTo(receiver)
	params.SetFrom(m.Options.Sender)
	params.SetTwiml("<response><say>" + message + "</say></response>")

	resp, err := m.Client.Api.CreateCall(params)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Error("Bad request: unable to call this number")
		return
	} else {
		fmt.Println("Call Status: " + *resp.Status)
		fmt.Println("Call Sid: " + *resp.Sid)
		fmt.Println("Call Direction: " + *resp.Direction)
	}
}

func (m OptionsWithHandler) findReciver(ctx *fasthttp.RequestCtx) string {
	sendOptions := new(options)
	*sendOptions = *m.Options
	const rcvKey = "receiver"
	args := ctx.QueryArgs()
	if nil != args && args.Has(rcvKey) {
		rcv := string(args.Peek(rcvKey))
		sendOptions.Receiver = rcv
	}
	return sendOptions.Receiver
}

// TODO change from hard-coded message to message retrieved from JSON WebHook
func (m OptionsWithHandler) findMessage(ctx *fasthttp.RequestCtx) (status string, alert string, err error) {
	body := ctx.PostBody()
	statu, _ := jsonparser.GetString(body, "status")
	fmt.Println(statu)
	return "firing", "Hello World", nil
}
