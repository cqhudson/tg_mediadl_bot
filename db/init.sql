
-- This script initializes all tables needed for the program

CREATE TABLE IF NOT EXISTS accounts (

    -- The Telegram account ID
    id BIGINT PRIMARY KEY NOT NULL,

    -- The Telegram username
    username VARCHAR(32),

    -- Usage metrics associated with each user
    download_attempts INT DEFAULT 0,
    download_successes INT DEFAULT 0,
    download_failures INT DEFAULT 0,
    download_send_failures INT DEFAULT 0,
    downloads_from_youtube INT DEFAULT 0,
    downloads_from_x INT DEFAULT 0,
    downloads_from_instagram INT DEFAULT 0,
    downloads_from_facebook INT DEFAULT 0,
    downloads_from_other_sources INT DEFAULT 0,
    
    -- Used to manage access to the bot
    status TEXT DEFAULT 'unlisted' CHECK(status IN ('admin', 'whitelisted', 'unlisted', 'blacklisted')), 
    
    -- Created when a row is inserted
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Updated whenever a row is updated
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- This trigger should update the updated_at column whenever the associated row is updated
CREATE TRIGGER IF NOT EXISTS update_accounts_updated_at
AFTER UPDATE on accounts
FOR EACH ROW
BEGIN
    UPDATE accounts SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- This Viw is for fetching download metrics for total overall downloads across users
CREATE VIEW IF NOT EXISTS total_downloads AS
SELECT
    SUM(download_attempts) AS total_attempted_downloads,
    SUM(download_successes) AS total_successful_downloads,
    SUM(download_failures) AS total_failed_downloads,
    SUM(download_send_failures) AS total_failed_download_sends,
    SUM(downloads_from_youtube) AS total_youtube,
    SUM(downloads_from_x) AS total_x,
    SUM(downloads_from_instagram) AS total_instagram,
    SUM(downloads_from_facebook) AS total_facebook,
    SUM(downloads_from_other_sources) AS total_other
FROM accounts;

