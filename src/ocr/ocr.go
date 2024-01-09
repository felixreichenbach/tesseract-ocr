package vision

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"strings"
	"sync"

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
	// Tessdata path to folder where language files are located
	DataPath string `json:"datapath"`
	// Tesseract configuration parameters see cmd line "tesseract --print-parameters"
	Parameters map[string]string `json:"parameters"`
}

// Validate OCR service configuration and return implicit dependencies
func (cfg *Config) Validate(path string) ([]string, error) {
	return []string{}, nil
}

// The OCR service model
type ocr struct {
	resource.Named
	logger logging.Logger
	mu     sync.Mutex

	// Tesseract client
	tessClient *gosseract.Client
}

// Handle ocr service configuration change
func (ocr *ocr) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	ocr.mu.Lock()
	defer ocr.mu.Unlock()
	if ocr.tessClient == nil {
		ocr.tessClient = gosseract.NewClient()
	}
	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}
	if newConf.DataPath != "" {
		ocr.logger.Infof("BEFORE: Tesseract Data Path: %s", ocr.tessClient.TessdataPrefix)
		if err := ocr.tessClient.SetTessdataPrefix(newConf.DataPath); err != nil {
			return err
		}
	}
	ocr.logger.Infof("AFTER: Tesseract Data Path: %s", ocr.tessClient.TessdataPrefix)

	languages, err := gosseract.GetAvailableLanguages()
	if err != nil {
		return err
	}
	ocr.logger.Infof("Available Languages: %s", strings.Join(languages, " | "))
	ocr.logger.Infof("Configuration Attributes: %s", conf.Attributes)
	for k, v := range newConf.Parameters {
		if err := ocr.tessClient.SetVariable(gosseract.SettableVariable(k), v); err != nil {
			return err
		}
	}
	return nil
}

// Download language files

func downloadLanguages() error {
	/*

		fileUrl := "https://gophercoding.com/img/logo-original.png"

		// Download the file, params:
		// 1) name of file to save as
		// 2) URL to download FROM


		err := DownloadFile("saveas.png", fileUrl)
		if err != nil {
			fmt.Println("Error downloading file: ", err)
			return
		}

		fmt.Println("Downloaded: " + fileUrl)
	*/
	return nil
}

// Process image with OCR
func (ocr *ocr) processOCR(buffer bytes.Buffer) ([]objectdetection.Detection, error) {
	if err := ocr.tessClient.SetImageFromBytes(buffer.Bytes()); err != nil {
		return nil, err
	}
	detections, err := ocr.tessClient.GetBoundingBoxesVerbose()
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
func (ocr *ocr) Close(ctx context.Context) error {
	return ocr.tessClient.Close()
}
