build:
	go build -o main ./cmd/main.go
pre_encode:
	ffmpeg -i video.mp4 \
		-map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0  \
		-map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0  \
		-b:v:0 350k  -c:v:0 libx264 -filter:v:0 "scale=320:-1"  \
		-b:v:1 1000k -c:v:1 libx264 -filter:v:1 "scale=640:-1"  \
		-b:v:2 3000k -c:v:2 libx264 -filter:v:2 "scale=1280:-1" \
		-b:v:3 245k  -c:v:3 libvpx-vp9 -filter:v:3 "scale=320:-1"  \
		-b:v:4 700k  -c:v:4 libvpx-vp9 -filter:v:4 "scale=640:-1"  \
		-b:v:5 2100k -c:v:5 libvpx-vp9 -filter:v:5 "scale=1280:-1"  \
		-b:v:6 350k  -c:v:6 libx264 -filter:v:6 "scale=320:-1,photosensitivity=30:0.6:20"  \
		-b:v:7 1000k -c:v:7 libx264 -filter:v:7 "scale=640:-1,photosensitivity=30:0.6:20"  \
		-b:v:8 3000k -c:v:8 libx264 -filter:v:8 "scale=1280:-1,photosensitivity=30:0.6:20" \
		-b:v:9 245k  -c:v:9 libvpx-vp9 -filter:v:9 "scale=320:-1,photosensitivity=30:0.6:20"  \
		-b:v:10 700k  -c:v:10 libvpx-vp9 -filter:v:10 "scale=640:-1,photosensitivity=30:0.6:20"  \
		-b:v:11 2100k -c:v:11 libvpx-vp9 -filter:v:11 "scale=1280:-1,photosensitivity=30:0.6:20"  \
		-use_timeline 1 -use_template 1 -window_size 12 -adaptation_sets "id=0,streams=0,2,4,6,8,10 id=1,streams=12,14,16,18,20,22 id=2,streams=a" \
		-hls_playlist true -f dash output/output.mpd