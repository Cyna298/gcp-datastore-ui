build:
	@echo "Building the backend service"
	cd backend && go build -o bin/service
# copy the out directory from ../webapp to ..
# so that the webapp is served from the backend
copy-webapp:
	rm -rf $(CURDIR)/backend/out
	echo "Copying webapp/out to backend/out directory"

	cp -r $(CURDIR)/webapp/out $(CURDIR)/backend/

run: build copy-webapp
	cd backend && ./bin/service

.PHONY: build copy-webapp run