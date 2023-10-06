package handler

import (
	"errors"
	"net/http"
	"reflect"
	"time"
	"unsafe"

	"github.com/kelseyhightower/envconfig"
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
				Host:     "localhost", // Should be unused
				Port:     8081,        // Should be unused
				Endpoint: "api",
			},
			// None of this should be used as we're passing a dummy connection
			DB: conf.DBConfiguration{
				Dialect:     "dummy",
				Driver:      "dummy",
				URL:         "dummy",
				Automigrate: false, // We don't have a database so we don't need to migrate
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

	// Remove /api from the path
	r.URL.Path = r.URL.Path[4:]

	handler := GetApiHandler(api)
	handler.ServeHTTP(w, r)
}

// This is a hack to get the handler from the API as it's not exported
func GetApiHandler(api *api.API) http.Handler {
	rs := reflect.ValueOf(api).Elem().FieldByName("handler")

	return reflect.NewAt(rs.Type(), unsafe.Pointer(rs.UnsafeAddr())).Elem().Interface().(http.Handler)
}

func GetConfiguration() conf.Configuration {
	var c conf.Configuration
	envconfig.Process("gitgateway", &c)
	return c
}

const INSTANCE_UUID = "d75d142d-2aca-45db-a42c-f68d6e8376e4"
const INSTANCE_ID = "ovo-cms"

func GetInstance() models.Instance {
	config := GetConfiguration()
	return models.Instance{
		ID:            INSTANCE_ID,
		UUID:          INSTANCE_UUID,
		CreatedAt:     time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:     nil,
		RawBaseConfig: "", // We hope this is not used
		BaseConfig:    &config,
	}
}

// It's like a database connection but it's fake
// We don't care about being multi-tenant, so we can just always return the same instance
type DummyConnection struct{}

func (*DummyConnection) Close() error {
	return nil
}

func (*DummyConnection) Automigrate() error {
	return errors.New("migration not supported")
}

func (*DummyConnection) GetInstanceByUUID(uuid string) (*models.Instance, error) {
	if uuid == INSTANCE_UUID {
		i := GetInstance()
		return &i, nil
	}
	return nil, errors.New("instance not found")
}

func (*DummyConnection) GetInstance(instanceID string) (*models.Instance, error) {
	if instanceID == INSTANCE_ID {
		i := GetInstance()
		return &i, nil
	}
	return nil, errors.New("instance not found")
}

func (*DummyConnection) CreateInstance(instance *models.Instance) error {
	if instance.ID == INSTANCE_ID || instance.UUID == INSTANCE_UUID {
		return errors.New("instance already exists")
	}
	return nil
}

func (*DummyConnection) DeleteInstance(instance *models.Instance) error {
	return nil
}

func (*DummyConnection) UpdateInstance(instance *models.Instance) error {
	return nil
}

func NewDummyConnection() *DummyConnection {
	return &DummyConnection{}
}
