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

docker-run:
	docker-compose up --build

docker-stop:
	docker-compose down

docker-clean:
	docker-compose down -v --rmi all

clean:
	rm -rf bin
	rm -rf results
	rm -rf uploads
