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

async function scrapeInstagram(url) {
    const browser = await puppeteer.launch({ headless: true });
    const page = await browser.newPage();

    // Log in to Instagram
    await page.goto('https://www.instagram.com/accounts/login/', { waitUntil: 'networkidle2' });

    // Accept cookies if the modal appears
    try {
        await page.waitForSelector('body > div.x1n2onr6.xzkaem6 > div.x9f619.x1n2onr6.x1ja2u2z > div > div.x1uvtmcs.x4k7w5x.x1h91t0o.x1beo9mf.xaigb6o.x12ejxvf.x3igimt.xarpa2k.xedcshv.x1lytzrv.x1t2pt76.x7ja8zs.x1n2onr6.x1qrby5j.x1jfb8zj > div > div > div > div > div.x7r02ix.xf1ldfh.x131esax.xdajt7p.xxfnqb6.xb88tzc.xw2csxc.x1odjw0f.x5fp0pe.x5yr21d.x19onx9a > div > button._a9--._ap36._a9_0', { timeout: 100 });
        const acceptCookiesButton = await page.$('body > div.x1n2onr6.xzkaem6 > div.x9f619.x1n2onr6.x1ja2u2z > div > div.x1uvtmcs.x4k7w5x.x1h91t0o.x1beo9mf.xaigb6o.x12ejxvf.x3igimt.xarpa2k.xedcshv.x1lytzrv.x1t2pt76.x7ja8zs.x1n2onr6.x1qrby5j.x1jfb8zj > div > div > div > div > div.x7r02ix.xf1ldfh.x131esax.xdajt7p.xxfnqb6.xb88tzc.xw2csxc.x1odjw0f.x5fp0pe.x5yr21d.x19onx9a > div > button._a9--._ap36._a9_0');
        if (acceptCookiesButton) {
            await acceptCookiesButton.click();
            await page.waitForNavigation({ waitUntil: 'networkidle2' });
        }
    } catch (e) {
    }

    await page.type('input[name="username"]', 'your-username'); // Replace with your Instagram username
    await page.type('input[name="password"]', 'your-password'); // Replace with your Instagram password
    await page.click('button[type="submit"]');
    await page.waitForNavigation({ waitUntil: 'networkidle2' });

    // Navigate to the target Instagram page
    await page.goto(url, { waitUntil: 'networkidle2' });

    // Extract channel name and followers count
    const result = await page.evaluate(() => {
        let channelName = document.querySelector('meta[property="og:title"]')?.getAttribute('content').trim() || '';
        const followersText = document.querySelector('meta[name="description"]')?.getAttribute('content').trim() || '';
        const followersMatch = followersText.match(/(\d+(?:\.\d+)?[KM]?) Followers/);
        let followersCount = followersMatch ? followersMatch[1] : '';

        // Convert followers count to a number
        followersCount = followersCount.toUpperCase();
        if (followersCount.includes('K')) {
            followersCount = followersCount.replace('K', '000');
        } else if (followersCount.includes('M')) {
            followersCount = followersCount.replace('M', '000000');
        }

        // Extract the username from the channel name
        if (channelName.includes('•')) {
            channelName = channelName.split('•')[0].trim();
            channelName = channelName.split(' ').pop();
            channelName = channelName.replace('(', '');
            channelName = channelName.replace(')', '');
            channelName = channelName.replace('@', '');
        }

        return { channelName, followersCount };
    });

    await browser.close();
    return result;
}

// Get the URL from command line arguments
const url = process.argv[2];
const domain = new URL(url).hostname;

if (domain.includes('vk.com')) {
    scrapeVK(url).then(result => {
        console.log(JSON.stringify(result));
    }).catch(error => {
        console.error('Error:', error);
    });
} else if (domain.includes('instagram.com')) {
    scrapeInstagram(url).then(result => {
        console.log(JSON.stringify(result));
    }).catch(error => {
        console.error('Error:', error);
    });
} else {
    console.error('Unsupported domain:', domain);
}
