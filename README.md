# Viam OCR Vision Service Module

This Viam vision service module uses [tesseract-ocr](https://github.com/tesseract-ocr/tesseract) through the [gosseract](https://pkg.go.dev/github.com/otiai10/gosseract/v2) wrapper and allows you to process an image extract text information from it. An example could be to extract license plate information to automatically open gates etc.

## Add the Module (local deploy)

We are are going to keep it very simple and deploy Configure a local and will look into using the Viam Registry in in a later part. Deploying local module is straight forward through the web user interface directly or by adding directly to the RAW JSON configuration.

```
{
  "name": "a-module-name",
  "executable_path": "<-- Path to the sensor binary including the binary -->",
  "type": "local"
}
```

## Configure Component

Add this sample configuration to the smart machine "components" part either in RAW JSON mode or through the we user interface by choosing "local service" in the menu.

```
    {
      "type": "vision",
      "namespace": "rdk",
      "model": "felixreichenbach:vision:ocr",
      "attributes": {
        "tessedit_pageseg_mode": "7",
        "tessedit_char_blacklist": "*+"
      },
      "name": "license"
    }
```
You can find a table of all possible tesseract configuration attributes [here](tesseract-config-params.md).

## Build the Module

From within the "src" directory run:

```go build -o ../bin/ocr .```

### BUILD INSTRUCTIONS MAC

To be able to successfully build the module, the following libraries are required.
I also uninstalled leptonica and tesseract with brew ignoring dependencies as mentioned here: https://github.com/otiai10/gosseract/issues/234#issuecomment-1707339205

```
wget http://www.leptonica.org/source/leptonica-1.78.0.tar.gz or https://github.com/DanBloomberg/leptonica/releases/tag/1.84.1
tar -xzvf leptonica-1.78.0.tar.gz
cd leptonica-1.78.0
./configure
make && sudo make install
```

```
brew install automake

git clone https://github.com/tesseract-ocr/tesseract.git
cd tesseract
./autogen.sh
./configure
make
sudo make install
```
