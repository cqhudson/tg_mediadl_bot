# Telegram Media Downloader Bot - [@tg_mediadl_bot](https://t.me/tg_mediadl_bot) 

A Telegram bot written in Go. This project was made for the purpose of learning the Go programming language by developing a tool I could use daily. 

Sending a video link to this bot will initiate a download, then the bot will send you the video in telegram. 

## dependencies
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) (used for handling the downloads). Stored in `exec/yt-dlp`
- [ffmpeg](https://ffmpeg.org/) (a dependency of yt-dlp)

## limitations
50mb file size upload limit from bots according to [Telegram Bot API](https://core.telegram.org/bots/faq)

<img width="808" height="84" alt="image" src="https://github.com/user-attachments/assets/bb304170-2131-4585-97b9-8f82924e6542" />


## How are downloads managed?
- YouTube `download\yt\video_id.mp4`
 
YouTube videos will be named after the video ID in the youtube url. For example, given the link `https://www.youtube.com/watch?v=n8-wN0lc5qk&list=RDn8-wN0lc5qk&start_radio=1`, the video id is the string of characters after `?v=`. So for this specific URL the video ID is `n8-wN0lc5qk` and the video will be saved as `n8-wN0lc5qk.mp4`. This will allow us to search for existing downloads, rather than redownloading the same video.

Each video should be deleted after a number of days to save on disk space. This might be implemented later.
