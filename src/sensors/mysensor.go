// This module implements a simple sensor returning the value of the setting attribute in the sensor's configuration.
// It extends the built-in resource subtype Sensor and implements its methods to handle resource construction, attribute configuration, and reconfiguration.

package mysensor

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

// Init called upon import, registers this component with the module
func init() {
	resource.RegisterComponent(sensor.API, Model, resource.Registration[sensor.Sensor, *Config]{Constructor: newSensor})
}

// To be used for functions which are not meant to be implemented in your component
var errUnimplemented = errors.New("unimplemented")

// Your model's colon-delimited-triplet (acme:demo:mybase). acme = namespace, demo = repo-name, mybase = model name
// If you plan to upload this module to the Viam registry, "acme" must match your Viam registry namespace.
var Model = resource.NewModel("viam-soleng", "go-resources", "gosensor")

// Maps JSON component configuration attributes.
type Config struct {
	Setting int `json:"setting"`
}

// Implement component configuration validation and and return implicit dependencies.
func (cfg *Config) Validate(path string) ([]string, error) {

	if cfg.Setting == 0 {
		return nil, fmt.Errorf(`expected "setting" attribute bigger than 0 for mysensor %q`, path)
	}
	return []string{}, nil
}

// Your sensor model type
type mySensor struct {
	resource.Named
	logger logging.Logger
	value  int // the value to be returned by the sensor.readings() method
}

// Sensor type constructor.
// Called upon sensor instantiation when a sensor model is added to the machine configuration
func newSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	s := &mySensor{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	if err := s.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return s, nil
}

// Reconfigure reconfigures with new settings.
func (s *mySensor) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {

	// This takes the generic resource.Config passed down from the parent and converts it to the
	// model-specific (aka "native") Config structure defined, above making it easier to directly access attributes.
	sensorConfig, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}
	s.logger.Infof(`Reconfiguring sensor value from %v to %v`, s.value, sensorConfig.Setting)
	s.value = sensorConfig.Setting
	return err
}

// Get sensor reading
func (s *mySensor) Readings(ctx context.Context, _ map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{"setting": s.value}, nil
}

// DoCommand can be implemented to extend sensor functionality but returns unimplemented in this example.
func (s *mySensor) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, errUnimplemented
}

// The close method is executed when the component is shut down
func (s *mySensor) Close(ctx context.Context) error {
	return nil
}
