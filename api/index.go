package handler

import (
	"errors"
	"net/http"
	"os"
	"reflect"
	"time"
	"unsafe"

	"github.com/kelseyhightower/envconfig"
	"github.com/netlify/git-gateway/api"
	"github.com/netlify/git-gateway/conf"
	"github.com/netlify/git-gateway/models"
	"github.com/sirupsen/logrus"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	config := GetConfiguration()
	ctx, err := api.WithInstanceConfig(r.Context(), &config, INSTANCE_ID)

	// Can't set info level in production because git gateway logs the access token :-/
	logrus.SetLevel(logrus.InfoLevel)
	if env := os.Getenv("VERCEL_ENV"); env != "" && env != "production" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api := api.NewAPIWithVersion(
		ctx,
		&conf.GlobalConfiguration{},
		NewDummyConnection(),
		"v1",
	)

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
	c.ApplyDefaults()
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
