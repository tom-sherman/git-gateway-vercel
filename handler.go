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
	config := getConfiguration()
	ctx, err := api.WithInstanceConfig(r.Context(), &config, instanceId)

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
		newDummyConnection(),
		"v1",
	)

	handler := getApiHandler(api)
	handler.ServeHTTP(w, r)
}

// This is a hack to get the handler from the API as it's not exported
func getApiHandler(api *api.API) http.Handler {
	rs := reflect.ValueOf(api).Elem().FieldByName("handler")

	return reflect.NewAt(rs.Type(), unsafe.Pointer(rs.UnsafeAddr())).Elem().Interface().(http.Handler)
}

func getConfiguration() conf.Configuration {
	var c conf.Configuration
	envconfig.Process("gitgateway", &c)
	c.ApplyDefaults()
	return c
}

const instanceUuid = "d75d142d-2aca-45db-a42c-f68d6e8376e4"
const instanceId = "ovo-cms"

func getInstance() models.Instance {
	config := getConfiguration()
	return models.Instance{
		ID:            instanceId,
		UUID:          instanceUuid,
		CreatedAt:     time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		DeletedAt:     nil,
		RawBaseConfig: "", // We hope this is not used
		BaseConfig:    &config,
	}
}

// It's like a database connection but it's fake
// We don't care about being multi-tenant, so we can just always return the same instance
type dummyConnection struct{}

func (*dummyConnection) Close() error {
	return nil
}

func (*dummyConnection) Automigrate() error {
	return errors.New("migration not supported")
}

func (*dummyConnection) GetInstanceByUUID(uuid string) (*models.Instance, error) {
	if uuid == instanceUuid {
		i := getInstance()
		return &i, nil
	}
	return nil, errors.New("instance not found")
}

func (*dummyConnection) GetInstance(instanceID string) (*models.Instance, error) {
	if instanceID == instanceId {
		i := getInstance()
		return &i, nil
	}
	return nil, errors.New("instance not found")
}

func (*dummyConnection) CreateInstance(instance *models.Instance) error {
	if instance.ID == instanceId || instance.UUID == instanceUuid {
		return errors.New("instance already exists")
	}
	return nil
}

func (*dummyConnection) DeleteInstance(instance *models.Instance) error {
	return nil
}

func (*dummyConnection) UpdateInstance(instance *models.Instance) error {
	return nil
}

func newDummyConnection() *dummyConnection {
	return &dummyConnection{}
}
