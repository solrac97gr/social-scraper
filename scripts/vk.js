const puppeteer = require('puppeteer');

// Array of user agents to rotate through
const userAgents = [
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
    'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0',
    'Mozilla/5.0 (X11; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36 Edg/92.0.902.55'
];

// Function to get a random user agent
function getRandomUserAgent() {
    return userAgents[Math.floor(Math.random() * userAgents.length)];
}

// Function to add random delay
function randomDelay(min = 0, max = 2000) {
    return new Promise(resolve => {
        const delay = Math.floor(Math.random() * (max - min + 1)) + min;
        setTimeout(resolve, delay);
    });
}

async function scrapeVK(url) {
    const browser = await puppeteer.launch(
        { headless: true, args: ['--no-sandbox', '--disable-setuid-sandbox'] }
    );
    const page = await browser.newPage();
    
    // Set a random user agent
    await page.setUserAgent(getRandomUserAgent());
    
    // Navigate to the URL
    await page.goto(url, { waitUntil: 'networkidle2' });

    // Handle multiple challenge pages in a loop
    let challengeAttempts = 0;
    const maxChallengeAttempts = 5; // Maximum number of challenge attempts
    
    while (challengeAttempts < maxChallengeAttempts) {
        const currentUrl = page.url();
        
        if (currentUrl.includes('/challenge.html')) {
           
            challengeAttempts++;
            
            // Add random delay before clicking the button (1-3 seconds)
            await randomDelay(1000, 3000);
            
            try {
                // Wait for the challenge button to be available
                await page.waitForSelector('body > div > button.start', { timeout: 15000 });
                
                // Add another small delay before clicking
                await randomDelay(500, 1500);
                
                // Click the anti-bot button
                await page.click('body > div > button.start');
                
                // Wait for navigation with a longer timeout
                await page.waitForNavigation({ 
                    waitUntil: 'networkidle2', 
                    timeout: 20000 
                });
                
                
                // Add delay after successful navigation
                await randomDelay(2000, 4000);
                
            } catch (error) {
                
                if (challengeAttempts >= maxChallengeAttempts) {
                    await browser.close();
                    return { channelName: 'Challenge Error - Max attempts reached', followersText: 'N/A' };
                }
                
                // Wait before retrying
                await randomDelay(3000, 5000);
            }
        } else {
            // No challenge page detected, break out of the loop
            break;
        }
    }
    
    // If we've reached max attempts and still on challenge page
    if (challengeAttempts >= maxChallengeAttempts && page.url().includes('/challenge.html')) {
        await browser.close();
        return { channelName: 'Challenge Error - Too many challenges', followersText: 'N/A' };
    }

    
    // Verify we're on the correct page before scraping
    const finalUrl = page.url();
    if (finalUrl.includes('/challenge.html')) {
        await browser.close();
        return { channelName: 'Still on challenge page', followersText: 'N/A' };
    }

    // Extract channel name and followers count with retry logic
    let scrapingAttempts = 0;
    const maxScrapingAttempts = 3;
    let result = { channelName: '', followersText: '' };
    
    while (scrapingAttempts < maxScrapingAttempts && (!result.channelName || !result.followersText)) {
        scrapingAttempts++;
        
        try {
            result = await page.evaluate(() => {
                const channelName = document.querySelector('.page_block.redesigned-cover-block .redesigned-group-info .redesigned-group-info__main .page_top h1')?.innerText.trim() || '';
                const followersText = document.querySelector('#page_subscribers > div > span')?.innerText.trim() || '';
                return { channelName, followersText };
            });
            
            // If we got some data, break out of the loop
            if (result.channelName || result.followersText) {
                break;
            }
            
        } catch (error) {
            console.error(`Scraping attempt ${scrapingAttempts} failed:`, error.message);
        }
        
        // Wait before retrying if not the last attempt
        if (scrapingAttempts < maxScrapingAttempts) {
            await page.waitForTimeout(3000);
        }
    }
    
    // If still no data after all attempts, return error
    if (!result.channelName && !result.followersText) {
        result = { channelName: 'Scraping Failed', followersText: 'N/A' };
    }

    await browser.close();
    return result;
}

// Get the URL from command line arguments
const url = process.argv[2];

if (!url) {
    console.error('Please provide a VK URL');
    process.exit(1);
}

scrapeVK(url).then(result => {
    console.log(JSON.stringify(result));
}).catch(error => {
    console.error('Error:', error);
    console.log(JSON.stringify({ channelName: 'VK Script Error', followersText: 'N/A' }));
});
