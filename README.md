# Telegram Followers Checker

A tool to extract information about Telegram channels from a list of links in an Excel file. This program scrapes channel name and followers count from Telegram channels.

## Usage

1. Install dependencies:
   ```sh
   go get -u github.com/PuerkitoBio/goquery
   go get -u github.com/xuri/excelize/v2
   ```

2. Prepare your Excel file:
   - Create an Excel file with links to Telegram channels (any format works as long as the links contain `t.me/` or `telegram.me/`).

3. Run the program:
   ```sh
   go run main.go /path/to/your_excel_file.xlsx
   ```

4. Check the output:
   - The program will generate an Excel file named `telegram_channels_followers.xlsx` with the extracted information.
   - The output includes channel name, followers count, and the original link.

## Example Result

After running the program, you will get an Excel file with the following format:

| Channel Name     | Followers Count | Original Link             |
|------------------|----------------:|---------------------------|
| Golang News      | 12500           | https://t.me/golang_news  |
| Tech Updates     | 45800           | https://t.me/tech_updates |
| Programming Tips | 8320            | https://t.me/coding_tips  |

The program provides real-time progress updates in the terminal:
```
Processing: https://t.me/golang_news
Processing: https://t.me/tech_updates
Processing: https://t.me/coding_tips

Success! Results saved to telegram_channels_followers.xlsx
```

## Features

- Automatically scrapes Telegram channel information
- Handles rate limiting by implementing delays between requests
- Well-formatted Excel output with styled headers
- Support for multiple Telegram link formats

## Requirements

- Go 1.13 or higher
- Dependencies:
  - github.com/PuerkitoBio/goquery (HTML parsing)
  - github.com/xuri/excelize/v2 (Excel file handling)

## License

This project is licensed under the MIT License.
