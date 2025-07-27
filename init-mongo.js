// MongoDB initialization script for social-scraper
// This script initializes the database and collections

db = db.getSiblingDB('social_scraper');

// Create initial collections (they will be created automatically when first used)
// Just ensuring the database exists
db.createCollection('influencer-analysis');
db.createCollection('users');

print('Database social_scraper initialized successfully');

print('Database initialization completed for social-scraper');
