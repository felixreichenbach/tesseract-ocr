package vision

import (
	"context"
	"image"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"
	viz "go.viam.com/rdk/vision"
	"go.viam.com/rdk/vision/classification"
	"go.viam.com/rdk/vision/objectdetection"
)

var Model = resource.NewModel("felixreichenbach", "vision", "ocr")

// Init called upon import, registers this component with the module
func init() {
	resource.RegisterComponent(vision.API, Model, resource.Registration[vision.Service, *Config]{Constructor: newOCR})
}

func newOCR(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (vision.Service, error) {
	ocr := &ocr{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	if err := ocr.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}
	return ocr, nil
}

// Maps JSON component configuration attributes.
type Config struct {
	Setting int `json:"setting"`
}

// Validate ocr service configuration and return implicit dependencies
func (cfg *Config) Validate(path string) ([]string, error) {
	return []string{}, nil
}

// The ocr service model
type ocr struct {
	resource.Named
	logger logging.Logger
}

// Classifications implements vision.Service.
func (*ocr) Classifications(ctx context.Context, img image.Image, n int, extra map[string]interface{}) (classification.Classifications, error) {
	panic("unimplemented")
}

// ClassificationsFromCamera implements vision.Service.
func (*ocr) ClassificationsFromCamera(ctx context.Context, cameraName string, n int, extra map[string]interface{}) (classification.Classifications, error) {
	panic("unimplemented")
}

// Close implements vision.Service.
func (*ocr) Close(ctx context.Context) error {
	panic("unimplemented")
}

// Detections implements vision.Service.
func (*ocr) Detections(ctx context.Context, img image.Image, extra map[string]interface{}) ([]objectdetection.Detection, error) {
	panic("unimplemented")
}

// DetectionsFromCamera implements vision.Service.
func (*ocr) DetectionsFromCamera(ctx context.Context, cameraName string, extra map[string]interface{}) ([]objectdetection.Detection, error) {
	panic("unimplemented")
}

// GetObjectPointClouds implements vision.Service.
func (*ocr) GetObjectPointClouds(ctx context.Context, cameraName string, extra map[string]interface{}) ([]*viz.Object, error) {
	panic("unimplemented")
}

// Handle ocr service configuration change
func (ocr *ocr) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {

	return nil
}
