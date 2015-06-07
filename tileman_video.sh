#!/usr/bin/sh

rm -rf /tmp/tileman_video

if [ -n "$1" ]; then
    region=$1
else
    region=deutschland
fi

./tileman --region=$region --ffmpeg-out=true --dir=/tmp/tileman_video --hours=5 --max-requests=10 && {
    ffmpeg -framerate 30 -i /tmp/tileman_video/%05d.png -c:v libx264 -r 30 -pix_fmt yuv420p /tmp/tileman_video/out.mp4 -loglevel quiet &&
    mpv --loop inf /tmp/tileman_video/out.mp4
}
