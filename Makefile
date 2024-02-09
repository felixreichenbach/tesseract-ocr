ocr: *.go */*.go go.*
	# the executable
	go build -o $@ --ldflags '-extldflags "-L/usr/local/lib/ -Bstatic -ltesseract"'
	# go build -o $@ -ldflags "-s -w" -tags osusergo,netgo
	file $@

module.tar.gz: ocr
	# the bundled module
	rm -f $@
	tar czf $@ $^