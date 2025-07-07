# Social Scraper 🤖

A tool to extract information from Telegram, Rutube, VK, and Instagram channels from a list of links in an Excel file. This program scrapes channel name and followers count from these channels.

<div align="center">
   <img src="images/Social%20Scraper.png" alt="Social Scraper" width="70%">
</div>

## Usage 🚀

1. Install dependencies:
   ```sh
   go get -u github.com/PuerkitoBio/goquery
   go get -u github.com/xuri/excelize/v2
   npm install puppeteer
   ```

2. Prepare your Excel file 📄:
   - Create an Excel file with links to Telegram, Rutube, VK, and Instagram channels (any format works as long as the links contain `t.me/`, `telegram.me/`, `rutube.ru/`, `vk.com/`, or `instagram.com/`).

3. Update the Puppeteer script 📝:
   - Replace the placeholders for Instagram username and password in `scripts/puppeteer_scraper.js` with your actual Instagram credentials.

4. Run the CLI 💻:
   ```sh
   go run cmd/cli/main.go /path/to/your_excel_file.xlsx
   ```

5. Run the HTTP server 🌐:
   ```sh
   go run cmd/http/main.go
   ```
   - Open your browser and navigate to `http://localhost:3000` to upload a file and download the processed file.

6. Check the output 📊:
   - The program will generate an Excel file in the `results` folder with the extracted information.
   - The output includes channel name, followers count, and the original link.

## Example Result 📈

After running the program, you will get an Excel file with the following format:

| Channel Name     | Followers Count | Original Link             | Platform | Registration Status |
|------------------|----------------:|---------------------------|----------|---------------------|
| Golang News      | 12500           | https://t.me/golang_news  | Telegram | registered          |
| Tech Updates     | 45800           | https://t.me/tech_updates | Telegram | not registered      |
| Programming Tips | 8320            | https://t.me/coding_tips  | Telegram | registered          |

The program provides real-time progress updates in the terminal:
```
Processing: https://t.me/golang_news
Processing: https://t.me/tech_updates
Processing: https://t.me/coding_tips

Success! Results saved to results/unique_id_channels_followers.xlsx
```

## Features ✨

- Automatically scrapes:
   - Telegram
   - Rutube
   - VK
   - Instagram
- Handles rate limiting by implementing delays between requests ⏳
- Well-formatted Excel output with styled headers 📑
- Support for multiple link formats 🔗

## Requirements 📋

- Go 1.13 or higher
- Node.js
- Dependencies:
  - github.com/PuerkitoBio/goquery (HTML parsing)
  - github.com/xuri/excelize/v2 (Excel file handling)
  - puppeteer (Headless browser for VK and Instagram scraping)

## License 📜

This project is licensed under the MIT License.
