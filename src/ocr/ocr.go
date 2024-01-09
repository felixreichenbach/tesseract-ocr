package vision

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
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
	TessDataLocal string `json:"tessdata_local"`
	// Tessdata download url
	TessDataRemote string `json:"tessdata_remote"`
	// Tesseract configuration parameters see cmd line "tesseract --print-parameters"
	Parameters map[string]string `json:"parameters"`
	// Tesseract languages to be used
	Languages []string `json:"languages"`
}

// Validate OCR service configuration and return implicit dependencies
func (cfg *Config) Validate(path string) ([]string, error) {
	if !strings.HasSuffix(cfg.TessDataLocal, "/") {
		return nil, resource.NewConfigValidationError(path, fmt.Errorf("tessdata_local path must end with /"))
	}
	if !strings.HasSuffix(cfg.TessDataRemote, "/") {
		return nil, resource.NewConfigValidationError(path, fmt.Errorf("tessdata_remote path must end with /"))
	}
	return []string{}, nil
}

// The OCR service model
type ocr struct {
	resource.Named
	logger logging.Logger
	mu     sync.Mutex

	// Tesseract client
	tessClient *gosseract.Client
	// Tesseract local folder
	tessDataLocal string
	// Tesseract download url
	tessDataRemote string
	// Tesseract languages to be used
	languages []string
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

	if ocr.tessDataLocal != newConf.TessDataLocal {
		ocr.tessDataLocal = newConf.TessDataLocal
	}

	if ocr.tessDataRemote != newConf.TessDataRemote {
		ocr.tessDataRemote = newConf.TessDataRemote
	}
	ocr.languages = newConf.Languages

	// Download tesseract language data files into the specified folder/path in tessdata_local
	if newConf.TessDataLocal != "" {
		ocr.DownloadTesseractFiles(ocr.tessDataLocal, ocr.tessDataRemote, ocr.languages)
		if err := ocr.tessClient.SetTessdataPrefix(newConf.TessDataLocal); err != nil {
			return err
		}
	}

	// Configure tesseract client ocr language setting
	ocr.tessClient.SetLanguage(ocr.languages...)

	// Apply tesseract parameters
	for k, v := range newConf.Parameters {
		if err := ocr.tessClient.SetVariable(gosseract.SettableVariable(k), v); err != nil {
			return err
		}
	}
	return nil
}

// DownloadFileTesseractFiles
func (ocr *ocr) DownloadTesseractFiles(local_path string, remote_url string, languages []string) error {

	// Create local folder for tesseract configuration and language files
	if err := os.MkdirAll(local_path, os.ModePerm); err != nil {
		ocr.logger.Infof("Local tessdata folder creation failed: %s", err)
	}
	// Download language files
	for _, lang := range languages {
		filename := lang + ".traineddata"
		// Check if file exists
		if _, err := os.Stat(local_path + filename); err != nil {
			ocr.logger.Infof("Downloading file %s from %s!", filename, remote_url)
			// Get the data
			resp, err := http.Get(remote_url + filename)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			// Create the file
			file, err := os.Create(local_path + filename)
			if err != nil {
				return err
			}
			defer file.Close()
			// Write the body to file
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				return err
			}
		}
	}
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
