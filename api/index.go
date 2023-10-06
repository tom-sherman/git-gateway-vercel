package handler

import (
	"net/http"
	"reflect"
	"time"
	"unsafe"

	"github.com/netlify/git-gateway/api"
	"github.com/netlify/git-gateway/conf"
	"github.com/netlify/git-gateway/models"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	api := api.NewAPI(
		&conf.GlobalConfiguration{
			API: struct {
				Host     string
				Port     int "envconfig:\"PORT\" default:\"8081\""
				Endpoint string
			}{
				Host:     "localhost",
				Port:     8081,
				Endpoint: "api",
			},
			DB: conf.DBConfiguration{
				Dialect:     "dummy",
				Driver:      "dummy",
				URL:         "dummy",
				Automigrate: false,
				Namespace:   "dummy",
			},
			Logging: conf.LoggingConfig{
				Level:            "info",
				File:             "dummy",
				DisableColors:    false,
				QuoteEmptyFields: true,
				TSFormat:         time.RFC3339Nano,
				Fields:           map[string]interface{}{},
			},
		},
		NewDummyConnection(),
	)

	handler := GetApiHandler(api)
	handler.ServeHTTP(w, r)
}

func GetApiHandler(api *api.API) http.Handler {
	rs := reflect.ValueOf(api).Elem().FieldByName("handler")

	return reflect.NewAt(rs.Type(), unsafe.Pointer(rs.UnsafeAddr())).Elem().Interface().(http.Handler)
}

type DummyConnection struct{}

func (d *DummyConnection) Close() error {
	return nil
}

func (d *DummyConnection) Automigrate() error {
	return nil
}

func (d *DummyConnection) GetInstanceByUUID(uuid string) (*models.Instance, error) {
	return nil, nil
}

func (d *DummyConnection) GetInstance(instanceID string) (*models.Instance, error) {
	return nil, nil
}

func (d *DummyConnection) CreateInstance(instance *models.Instance) error {
	return nil
}

func (d *DummyConnection) DeleteInstance(instance *models.Instance) error {
	return nil
}

func (d *DummyConnection) UpdateInstance(instance *models.Instance) error {
	return nil
}

func NewDummyConnection() *DummyConnection {
	return &DummyConnection{}
}
