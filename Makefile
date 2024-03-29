tesseract-ocr:
	# the executable
	#sudo apt update
	apt-get -y install libleptonica-dev
	apt-get -y install libtool
	git clone https://github.com/tesseract-ocr/tesseract.git
	cd tesseract && ./autogen.sh && ./configure && make && make install
	go build -o $@ -ldflags '-extldflags "-L/usr/local/lib/ -Bstatic -ltesseract"'
	file $@

module.tar.gz: tesseract-ocr
	# the bundled module
	rm -f $@
	tar czf $@ $^
