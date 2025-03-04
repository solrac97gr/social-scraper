const puppeteer = require('puppeteer');

async function checkRegistrationStatus(url) {
  const browser = await puppeteer.launch({
    headless: true,
    args: [
      '--no-sandbox',
      '--disable-setuid-sandbox',
      '--disable-dev-shm-usage',
      '--disable-accelerated-2d-canvas',
      '--no-first-run',
      '--no-zygote',
      '--single-process',
      '--disable-gpu'
    ]
  });
  const page = await browser.newPage();
  await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36');
  await page.goto('https://www.gosuslugi.ru/snet', { waitUntil: 'networkidle2' });

  // Fill out the form and submit
  try {
    await page.waitForSelector('input[role="textbox"]');
    await page.type('input[role="textbox"]', url);
    await page.click('#tabpanel-link > lib-button > div > button');
  } catch (error) {
    await browser.close();
    return false;
  }

  // Wait for the result
  try {
    await page.waitForSelector('#print-page > app-snet > section:nth-child(3) > div > div > app-status-snet > div > div.flex-container.align-items-center > p');
    const result = await page.evaluate(() => {
      return document.querySelector('#print-page > app-snet > section:nth-child(3) > div > div > app-status-snet > div > div.flex-container.align-items-center > p').innerText;
    });
    await browser.close();
    return result.includes('Ресурс включён в перечень Роскомнадзора');
  } catch (error) {
    await browser.close();
    return false;
  }
}


const url = process.argv[2];

checkRegistrationStatus(url).then(isRegistered => {
  console.log(JSON.stringify({ isRegistered }));
}).catch(error => {
  // save error to file
  fs.writeFileSync('error.txt', error);
  console.log(JSON.stringify({ isRegistered }));
});
