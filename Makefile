tesseract-ocr:
	# the executable
	sudo apt-get install -y libleptonica-dev
	go build -o $@ -ldflags '-extldflags "-L/usr/local/lib/ -Bstatic -ltesseract"'
	file $@

module.tar.gz: tesseract-ocr
	# the bundled module
	rm -f $@
	tar czf $@ $^
