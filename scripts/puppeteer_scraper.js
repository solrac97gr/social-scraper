const puppeteer = require('puppeteer');

async function scrapeVK(url) {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    await page.goto(url, { waitUntil: 'networkidle2' });

    // Extract channel name and followers count
    const result = await page.evaluate(() => {
        const channelName = document.querySelector('.page_block.redesigned-cover-block .redesigned-group-info .redesigned-group-info__main .page_top h1')?.innerText.trim() || '';
        const followersText = document.querySelector('#page_subscribers > div > span')?.innerText.trim() || '';
        return { channelName, followersText };
    });

    await browser.close();
    return result;
}

// Get the URL from command line arguments
const url = process.argv[2];
scrapeVK(url).then(result => {
    console.log(JSON.stringify(result));
}).catch(error => {
    console.error('Error:', error);
});
