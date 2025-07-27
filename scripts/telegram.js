const puppeteer = require('puppeteer');

// Array of user agents to rotate through
const userAgents = [
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
    'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0'
];

// Function to get a random user agent
function getRandomUserAgent() {
    return userAgents[Math.floor(Math.random() * userAgents.length)];
}

// Function to add random delay
function randomDelay(min = 1000, max = 3000) {
    return new Promise(resolve => {
        const delay = Math.floor(Math.random() * (max - min + 1)) + min;
        setTimeout(resolve, delay);
    });
}

async function scrapeTelegram(url) {
    const browser = await puppeteer.launch({
        headless: true,
        args: [
            '--no-sandbox',
            '--disable-setuid-sandbox',
            '--disable-blink-features=AutomationControlled',
            '--disable-web-security',
            '--disable-features=VizDisplayCompositor'
        ]
    });
    
    const page = await browser.newPage();
    let result = { channelName: 'Telegram Error', followersCount: 'N/A' };
    
    try {
        // Set a random user agent
        await page.setUserAgent(getRandomUserAgent());
        
        // Set viewport
        await page.setViewport({ width: 1366, height: 768 });
        
        // Block unnecessary resources to speed up loading
        await page.setRequestInterception(true);
        page.on('request', (req) => {
            if (req.resourceType() === 'stylesheet' || req.resourceType() === 'font') {
                req.abort();
            } else {
                req.continue();
            }
        });
        
        // Add random delay before navigation
        await randomDelay(500, 1500);
        
        await page.goto(url, { 
            waitUntil: 'networkidle2',
            timeout: 30000 
        });
        
        // Wait for page to load
        await randomDelay(2000, 4000);
        
        // Try to extract channel name
        let channelName = 'Unknown';
        try {
            await page.waitForSelector('div.tgme_page_title', { timeout: 10000 });
            channelName = await page.$eval('div.tgme_page_title', el => el.textContent.trim());
        } catch (error) {
            // Try alternative selectors if main selector fails
            
            // Try alternative selectors
            const titleSelectors = [
                '.tgme_page_title',
                'h1',
                '.page-title',
                '.channel-title'
            ];
            
            for (const selector of titleSelectors) {
                try {
                    const element = await page.$(selector);
                    if (element) {
                        channelName = await page.evaluate(el => el.textContent.trim(), element);
                        if (channelName && channelName !== '') {
                            break;
                        }
                    }
                } catch (e) {
                    // Continue to next selector
                }
            }
        }
        
        // Try to extract followers count
        let followersCount = '0';
        try {
            await page.waitForSelector('div.tgme_page_extra', { timeout: 10000 });
            
            const extraText = await page.$eval('div.tgme_page_extra', el => el.textContent.trim());
            
            // Look for subscriber/member/follower count
            if (extraText.includes('subscriber') || extraText.includes('member') || extraText.includes('follower')) {
                const matches = extraText.match(/[\d\s]+/);
                if (matches && matches[0]) {
                    followersCount = matches[0].replace(/\s/g, '');
                }
            }
        } catch (error) {
            // Try alternative selectors for followers
            
            // Try alternative selectors for followers
            const followerSelectors = [
                '.tgme_page_extra',
                '.subscribers-count',
                '.members-count',
                '.followers-count',
                '[class*="subscriber"]',
                '[class*="member"]'
            ];
            
            for (const selector of followerSelectors) {
                try {
                    const elements = await page.$$(selector);
                    for (const element of elements) {
                        const text = await page.evaluate(el => el.textContent.trim(), element);
                        if (text && (text.includes('subscriber') || text.includes('member') || text.includes('follower'))) {
                            const matches = text.match(/[\d\s]+/);
                            if (matches && matches[0]) {
                                followersCount = matches[0].replace(/\s/g, '');
                                break;
                            }
                        }
                    }
                    if (followersCount !== '0') break;
                } catch (e) {
                    // Continue to next selector
                }
            }
        }
        
        result = {
            channelName: channelName || 'Unknown',
            followersCount: followersCount || '0'
        };
        
    } catch (error) {
        result = { channelName: 'Telegram Error', followersCount: 'N/A' };
    } finally {
        await browser.close();
    }
    
    return result;
}

// Main execution
async function main() {
    const url = process.argv[2];
    
    if (!url) {
        process.exit(1);
    }
    
    try {
        const result = await scrapeTelegram(url);
        console.log(JSON.stringify(result));
    } catch (error) {
        console.log(JSON.stringify({ channelName: 'Telegram Error', followersCount: 'N/A' }));
    }
}

// Run if this script is executed directly
if (require.main === module) {
    main();
}

module.exports = { scrapeTelegram };
