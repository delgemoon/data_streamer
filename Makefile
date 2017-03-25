build:
	go build

stream:
	# go build; ffmpeg -y -i input.mkv -vcodec libx264 -b:v 0.03M -vf scale=202:-1 -r 15 -f h264 - < /dev/null | ./data_streamer write connection_id
	ffmpeg -y -i input.mkv -vcodec libx264 -b:v 0.03M -vf scale=202:-1 -r 15 -f h264 - | ./data_streamer write connection_id

read:
	./data_streamer read connection_id | ffmpeg -y -re -i pipe:0 frame_%03d.ppm
	# ./data_streamer read connection_id | ffmpeg -y -i pipe:0 -map 0 -flags:v +global_header -vcodec libx264  works.mp4
	# ./data_streamer read connection_id > finalout.mp4;
	# ./data_streamer read connection_id

debug:
	# ffmpeg -i works.mp4 frame%03d.jpg
	# THIS WORKS
	# ffmpeg -y -i input.mkv -vcodec libx264 -b:v 0.03M -vf scale=202:-1 -r 15 -f h264 - | ffmpeg -y -i pipe:0 -map 0 -flags:v +global_header -vcodec libx264 works.mp4
	#
	# ffmpeg -y -i input.mkv -vcodec libx264 -b:v 0.03M -vf scale=202:-1 -r 15 -f h264 - | ffmpeg -y -i pipe:0 frame_%03d.ppm
	# ffmpeg -i input.mkv frame_%03d.ppm
	ffmpeg -i input.mkv -f h264 - | ffmpeg -i pipe:0 frame_%03d.ppm
