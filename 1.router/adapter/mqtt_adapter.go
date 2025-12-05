package adapter

import (
	"context"
	"errors"
	"time"

	request "batchLog/0.core/commonResReq/req"
	response "batchLog/0.core/commonResReq/res"
	"batchLog/0.core/logafa"

	jsoniter "github.com/json-iterator/go"
)

type MQTTContext struct {
	payload     string
	jwt         string
	clientID    string
	ip          string
	requestTime time.Time

	ctx    context.Context
	cancel context.CancelFunc
}

func NewMQTTContext(payload, jwt, clientID, ip string, requestTime time.Time) request.RequestContext {
	ctx := context.Background()
	return &MQTTContext{
		payload:     payload,
		jwt:         jwt,
		clientID:    clientID,
		ip:          ip,
		requestTime: requestTime,

		ctx:    ctx,
		cancel: nil,
	}
}

// Create new context
func (m *MQTTContext) GetContext() context.Context {
	return m.ctx
}
func (m *MQTTContext) SetContext(ctx context.Context) {
	m.ctx = ctx
}
func (m *MQTTContext) Cancel() {
	if m.cancel != nil {
		m.cancel()
	}
}
func (m *MQTTContext) SetCancel(c context.CancelFunc) {
	m.cancel = c
}

// BindJSON implements request.RequestContext.
func (m *MQTTContext) BindJSON(obj interface{}) error {
	if m.payload == "" || m.payload == "{}" {
		return errors.New("empty payload")
	}
	return jsoniter.UnmarshalFromString(m.payload, obj)
}

// GetClientID implements request.RequestContext.
func (m *MQTTContext) GetClientID() string {
	return m.clientID
}

// GetClientIP implements request.RequestContext.
func (m *MQTTContext) GetClientIP() string {
	return m.ip
}

// GetJWT implements request.RequestContext.
func (m *MQTTContext) GetJWT() string {
	return m.jwt
}

// GetRequestTime implements request.RequestContext.
func (m *MQTTContext) GetRequestTime() time.Time {
	return m.requestTime
}

// Success implements request.RequestContext.
func (m *MQTTContext) Success(data interface{}) {
	err := m.ctx.Err()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logafa.Debug("Success() 被阻止：Context Timeout")
		} else if errors.Is(err, context.Canceled) {
			logafa.Debug("Success() 被阻止：Context 已取消")
		}
		return
	}

	response.SuccessMqtt(m.getResponseTopic(), m.requestTime, data)
}

// Error implements request.RequestContext.
func (m *MQTTContext) Error(code int, message string) {
	errTopic := "errReq/" + m.clientID
	response.ErrorMqtt(errTopic, code, m.requestTime, message)
}

func (m *MQTTContext) getResponseTopic() string {
	var temp struct {
		SubscribeTo string `json:"subscribeTo"`
	}
	jsoniter.UnmarshalFromString(m.payload, &temp)
	return temp.SubscribeTo
}
