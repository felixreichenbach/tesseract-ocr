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
      "name": "license-plates",
      "type": "vision",
      "namespace": "rdk",
      "model": "felixreichenbach:vision:ocr",
      "attributes": {
        "languages": [
          "eng"
        ],
        "parameters": {
          "tessedit_char_blacklist": "*+",
          "tessedit_pageseg_mode": "7"
        },
        "tessdata_local": "./tessdata/",
        "tessdata_remote": "https://github.com/tesseract-ocr/tessdata_fast/raw/main/"
      }
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


### Build on Ubuntu

```
sudo apt install build-essential

sudo apt install pkg-config

# If Unable to find a valid copy of libtoolize or glibtoolize in your PATH!
sudo apt-get install libtool

sudo apt install libjpeg-dev

sudo apt-get install libleptonica-dev

# Install tesseract libraries
git clone https://github.com/tesseract-ocr/tesseract.git
cd tesseract
./autogen.sh
./configure
# make seems not required -> done as part of the next step
sudo make install

go build -o ../bin/ocr .

alternatively:
CGO_ENABLED=1 GOARCH=arm64 go build --ldflags '-extldflags "-fopenmp -L/usr/local/lib/ -Bstatic -ltesseract"' -o ../bin/tesseract-ocr .

# Check dynamically linked libs on file
ldd fileName

# Missing tesseract lib fixed with:
sudo ldconfig


# Build go binary with adding tesseract library statically
go build --ldflags '-extldflags "-fopenmp -L/usr/local/lib/ -Bstatic -ltesseract"' -o ../bin/tesseract-ocr .


```

### Build AppImage

Sample makefile: https://github.com/jeremyrhyde/viam-rplidar/blob/main/Makefile

