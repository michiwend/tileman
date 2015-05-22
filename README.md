# tileman
Download weather radar images from kachelmannwetter.com

## Usage
type `tileman -help` to get the following list

Flag        | Description                                      | Default
------------|--------------------------------------------------|----------------
-dir        | Directory for saving the results.                | ./tileman_out
-end-date   | End date in the form "2006-01-20"                | today
-end-time   | End time in the form "15:04"                     | 15 minutes ago
-region     | Which region map to use?                         | germany
-res        | Time resolution. Use a multiple of 5, minimum 5! | 5
-start-date | Start date in the form "2006-01-20"              | today
-start-time | Start time in the form "15:04"                   | 2 hours ago
-ffmpeg-out | Generate files in the form 00001.png             | false

### Create a video with ffmpeg
If you want to create a video, you can use `-ffmpeg-out=true` to have
numbered filenames and afterwards let ffmpeg do the job:
```
ffmpeg -framerate 30 -i %05d.png -c:v libx264 -r 30 -pix_fmt yuv420p out.mp4
```


