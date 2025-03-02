deps:
	go get -u github.com/PuerkitoBio/goquery
	go get -u github.com/xuri/excelize/v2
	npm install puppeteer

build:
	go build -o bin/scraper main.go

run:
	go run cmd/http/main.go

cli:
	go run cmd/cli/main.go

clean:
	rm -rf bin