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

async function scrapeTikTok(url) {
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
    let result = { channelName: 'TikTok Error', followersCount: 'N/A' };
    
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
        
        // Navigate to the URL with a longer timeout
        await page.goto(url, { 
            waitUntil: 'networkidle2', 
            timeout: 30000 
        });
        
        // Add random delay to mimic human behavior
        await randomDelay(2000, 4000);
        
        // Wait for the content to load
        await page.waitForSelector('#main-content-others_homepage', { timeout: 20000 });
        
        // Extract data with multiple attempts
        result = { channelName: '', followersCount: '' };
        let attempts = 0;
        const maxAttempts = 3;
        
        while (attempts < maxAttempts && (!result.channelName || !result.followersCount)) {
            attempts++;
            
            try {
                result = await page.evaluate(() => {
                    // Channel name selector (updated)
                    const channelNameElement = document.querySelector('#main-content-others_homepage > div > div.e1457k4r14.css-cooqqt-DivShareLayoutHeader-StyledDivShareLayoutHeaderV2-CreatorPageHeader.e13xij562 > div.css-1o9t6sm-DivShareTitleContainer-CreatorPageHeaderShareContainer.e1457k4r15 > div.css-dozy74-DivUserIdentifierWrapper.e1gnmlil1 > div > div > h1');
                    
                    // Followers count selector
                    const followersElement = document.querySelector('#main-content-others_homepage > div > div.e1457k4r14.css-cooqqt-DivShareLayoutHeader-StyledDivShareLayoutHeaderV2-CreatorPageHeader.e13xij562 > div.css-1o9t6sm-DivShareTitleContainer-CreatorPageHeaderShareContainer.e1457k4r15 > div.css-1ygxkc0-CreatorPageHeaderTextContainer.e1457k4r16 > h3 > div:nth-child(2) > strong');
                    
                    const channelName = channelNameElement ? channelNameElement.innerText.trim() : '';
                    const followersCount = followersElement ? followersElement.innerText.trim() : '';
                    
                    return { channelName, followersCount };
                });
                
                // If we got some data, break out of the loop
                if (result.channelName || result.followersCount) {
                    break;
                }
                
            } catch (error) {
                console.error(`TikTok scraping attempt ${attempts} failed:`, error.message);
            }
            
            // Wait before retrying if not the last attempt
            if (attempts < maxAttempts) {
                await randomDelay(2000, 4000);
            }
        }
        
        // If still no data after all attempts, try alternative selectors
        if (!result.channelName && !result.followersCount) {
            try {
                result = await page.evaluate(() => {
                    // Try alternative selectors
                    let channelName = '';
                    let followersCount = '';
                    
                    // Try different selectors for channel name
                    const channelSelectors = [
                        'h2[data-e2e="user-title"]',
                        '[data-e2e="user-title"]',
                        'h1',
                        'h2',
                        '.tiktok-username'
                    ];
                    
                    for (const selector of channelSelectors) {
                        const element = document.querySelector(selector);
                        if (element && element.innerText.trim()) {
                            channelName = element.innerText.trim();
                            break;
                        }
                    }
                    
                    // Try different selectors for followers count
                    const followersSelectors = [
                        '[data-e2e="followers-count"]',
                        'strong[data-e2e="followers-count"]',
                        '.number[data-e2e="followers-count"]',
                        'div[data-e2e="followers-count"] strong',
                        '.follower-count',
                        'strong'
                    ];
                    
                    for (const selector of followersSelectors) {
                        const element = document.querySelector(selector);
                        if (element && element.innerText.trim() && 
                            (element.innerText.includes('K') || element.innerText.includes('M') || 
                             /^\d+/.test(element.innerText.trim()))) {
                            followersCount = element.innerText.trim();
                            break;
                        }
                    }
                    
                    return { channelName, followersCount };
                });
            } catch (error) {
                console.error('Alternative selector attempt failed:', error.message);
            }
        }
        
        // If still no data, return error
        if (!result.channelName && !result.followersCount) {
            result = { channelName: 'TikTok Scraping Failed', followersCount: 'N/A' };
        }
        
    } catch (error) {
        console.error('TikTok scraping error:', error.message);
        result = { channelName: 'TikTok Error', followersCount: 'N/A' };
    } finally {
        await browser.close();
    }
    
    return result;
}

// Get the URL from command line arguments
const url = process.argv[2];

if (!url) {
    console.error('Please provide a TikTok URL');
    process.exit(1);
}

scrapeTikTok(url).then(result => {
    console.log(JSON.stringify(result));
}).catch(error => {
    console.error('Error:', error);
    console.log(JSON.stringify({ channelName: 'Script Error', followersCount: 'N/A' }));
});