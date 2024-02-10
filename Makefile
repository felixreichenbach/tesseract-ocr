tesseract-ocr:
	# the executable
	#sudo apt update
	apt-get -y install libleptonica-dev
	go build -o $@ -ldflags '-extldflags "-L/usr/local/lib/ -Bstatic -ltesseract"'
	file $@

module.tar.gz: tesseract-ocr
	# the bundled module
	rm -f $@
	tar czf $@ $^
