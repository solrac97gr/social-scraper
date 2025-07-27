deps:
	go get -u github.com/PuerkitoBio/goquery
	go get -u github.com/xuri/excelize/v2
	npm install puppeteer

build:
	go build -o bin/scraper main.go

run:
	infracli run mongo
	go run cmd/http/main.go

cli:
	infracli run mongo
	go run cmd/cli/main.go

docker:
	@echo "🐋 Launching Social Scraper with Docker..."
	@chmod +x docker-start.sh
	@./docker-start.sh

clean:
	rm -rf bin
	rm -rf results
	rm -rf uploads