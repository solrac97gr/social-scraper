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

async function scrapeRutube(url) {
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
    let result = { channelName: 'Rutube Error', followersCount: 'N/A' };
    
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
            // Primary selector based on the Go extractor
            const titleSelector = 'h1.wdp-feed-banner-module__wdp-feed-banner__title-text';
            await page.waitForSelector(titleSelector, { timeout: 10000 });
            
            const titleElement = await page.$(titleSelector);
            if (titleElement) {
                channelName = await page.evaluate(el => el.getAttribute('title') || el.textContent.trim(), titleElement);
            }
        } catch (error) {
            // Try alternative selectors if main selector fails
            
            // Try alternative selectors
            const titleSelectors = [
                'h1.wdp-feed-banner-module__wdp-feed-banner__title-text',
                '.wdp-feed-banner__title-text',
                'h1[class*="title"]',
                '.channel-title',
                'h1',
                '.page-title'
            ];
            
            for (const selector of titleSelectors) {
                try {
                    const element = await page.$(selector);
                    if (element) {
                        const text = await page.evaluate(el => {
                            return el.getAttribute('title') || el.textContent.trim();
                        }, element);
                        if (text && text !== '') {
                            channelName = text;
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
            // Primary selector based on the Go extractor
            const followersSelector = '.wdp-feed-banner-module__wdp-feed-banner__title p';
            await page.waitForSelector(followersSelector, { timeout: 10000 });
            
            const followersText = await page.$eval(followersSelector, el => el.textContent.trim());
            
            // Clean up the followers string to remove non-digit characters
            const cleanedFollowers = followersText.replace(/\D/g, '');
            if (cleanedFollowers) {
                followersCount = cleanedFollowers;
            }
        } catch (error) {
            // Try alternative selectors for followers
            
            // Try alternative selectors for followers
            const followerSelectors = [
                '.wdp-feed-banner-module__wdp-feed-banner__title p',
                '.wdp-feed-banner__title p',
                '[class*="subscriber"]',
                '[class*="follower"]',
                '.subscribers-count',
                '.followers-count',
                'p[class*="banner"]'
            ];
            
            for (const selector of followerSelectors) {
                try {
                    const elements = await page.$$(selector);
                    for (const element of elements) {
                        const text = await page.evaluate(el => el.textContent.trim(), element);
                        if (text) {
                            const cleanedText = text.replace(/\D/g, '');
                            if (cleanedText && cleanedText !== '') {
                                followersCount = cleanedText;
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
        result = { channelName: 'Rutube Error', followersCount: 'N/A' };
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
        const result = await scrapeRutube(url);
        console.log(JSON.stringify(result));
    } catch (error) {
        console.log(JSON.stringify({ channelName: 'Rutube Error', followersCount: 'N/A' }));
    }
}

// Run if this script is executed directly
if (require.main === module) {
    main();
}

module.exports = { scrapeRutube };
