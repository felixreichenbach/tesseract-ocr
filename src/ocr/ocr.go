package vision

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"
	viz "go.viam.com/rdk/vision"
	"go.viam.com/rdk/vision/classification"
	"go.viam.com/rdk/vision/objectdetection"

	"github.com/otiai10/gosseract/v2"
)

var Model = resource.NewModel("felixreichenbach", "vision", "ocr")

// Init called upon import, registers this OCR service with the module
func init() {
	resource.RegisterDefaultService(vision.API, Model, resource.Registration[vision.Service, *Config]{Constructor: newOCR})
}

// Instantiates an OCR vision service
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

// OCR vision service configuration attributes
type Config struct {
}

// Validate OCR service configuration and return implicit dependencies
func (cfg *Config) Validate(path string) ([]string, error) {
	return []string{}, nil
}

// The OCR service model
type ocr struct {
	resource.Named
	logger logging.Logger
}

// Handle ocr service configuration change
func (ocr *ocr) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	return nil
}

// Process image with OCR
func processOCR(buffer bytes.Buffer) ([]objectdetection.Detection, error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetPageSegMode(gosseract.PageSegMode(3))
	client.SetLanguage()
	client.SetImageFromBytes(buffer.Bytes())

	detections, err := client.GetBoundingBoxesVerbose()
	if err != nil {
		return nil, err
	}
	result := []objectdetection.Detection{}
	for _, detection := range detections {
		newDetection := objectdetection.NewDetection(detection.Box, detection.Confidence, detection.Word)
		result = append(result, newDetection)
	}
	return result, nil
}

// Detections implements vision.Service.
func (ocr *ocr) Detections(ctx context.Context, img image.Image, extra map[string]interface{}) ([]objectdetection.Detection, error) {
	image_buf := new(bytes.Buffer)
	if err := jpeg.Encode(image_buf, img, nil); err != nil {
		return nil, err
	}
	result, _ := processOCR(*image_buf)
	return result, nil
}

// DetectionsFromCamera implements vision.Service.
func (ocr *ocr) DetectionsFromCamera(ctx context.Context, cameraName string, extra map[string]interface{}) ([]objectdetection.Detection, error) {
	// @TODO How to get the robot or camera directly without adding dependencies?
	panic("unimplemented")
}

// GetObjectPointClouds implements vision.Service.
func (ocr *ocr) GetObjectPointClouds(ctx context.Context, cameraName string, extra map[string]interface{}) ([]*viz.Object, error) {
	panic("unimplemented")
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
