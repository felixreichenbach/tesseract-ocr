package vision

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"strconv"
	"sync"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"
	viz "go.viam.com/rdk/vision"
	"go.viam.com/rdk/vision/classification"
	"go.viam.com/rdk/vision/objectdetection"
	"go.viam.com/utils"

	"github.com/otiai10/gosseract/v2"
	"github.com/pkg/errors"
)

var Model = resource.NewModel("felixreichenbach", "vision", "ocr")

// Init called upon import, registers this OCR service with the module
func init() {
	resource.RegisterService(vision.API, Model, resource.Registration[vision.Service, *Config]{Constructor: newOCR})
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
	// The tessdata prefix path for the trained data
	TessdataPrefix string `json:"tessdataprefix"`
	// The page segmentation mode "PSM"
	PSM int `json:"psm"`
	// The languages to use
	Languages []string `json:"langugages"`
}

// Validate OCR service configuration and return implicit dependencies
func (cfg *Config) Validate(path string) ([]string, error) {
	if !((cfg.PSM >= 0) && (cfg.PSM <= 13)) {
		return nil, utils.NewConfigValidationError(path, errors.Errorf("PSM must be in the range of 0-13 integer."))
	}
	return []string{}, nil
}

// The OCR service model
type ocr struct {
	resource.Named
	logger logging.Logger
	mu     sync.Mutex
	// Path to tessdata folder containing traineddata
	tessdataPrefix string
	// Page segmentation mode setting
	psm int
	// Languages
	languages []string
}

// Handle ocr service configuration change
func (ocr *ocr) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	ocr.mu.Lock()
	defer ocr.mu.Unlock()
	// Set TessdataPrefix path
	ocr.tessdataPrefix = conf.Attributes.String("tessdataprefix")
	// Set the configured psm value else default to 3 which is tesseract's default psm value
	ocr.psm = conf.Attributes.Int("psm", 3)
	// Set language models to use
	ocr.languages = conf.Attributes.StringSlice("languages")
	return nil
}

// Process image with OCR
func (ocr *ocr) processOCR(buffer bytes.Buffer) ([]objectdetection.Detection, error) {
	client := gosseract.NewClient()
	defer client.Close()

	// TODO: Not tested yet
	if ocr.tessdataPrefix != "" {
		client.TessdataPrefix = ocr.tessdataPrefix
	}
	// Set the Page Segmentation Mode "PSM"
	ocr.logger.Infof("OCR:PSM set to: %s", ocr.psm)
	if err := client.SetVariable("tessedit_pageseg_mode", strconv.Itoa(ocr.psm)); err != nil {
		return nil, err
	}
	// TODO: Not tested yet
	if len(ocr.languages) != 0 {
		if err := client.SetLanguage(ocr.languages...); err != nil {
			return nil, err
		}
	}
	if err := client.SetImageFromBytes(buffer.Bytes()); err != nil {
		return nil, err
	}
	client.Text()
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
	result, _ := ocr.processOCR(*image_buf)
	return result, nil
}

// DetectionsFromCamera implements vision.Service.
func (ocr *ocr) DetectionsFromCamera(ctx context.Context, cameraName string, extra map[string]interface{}) ([]objectdetection.Detection, error) {
	// TODO: Add cameras as dependencies and then use the one provided to choose out of them

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
