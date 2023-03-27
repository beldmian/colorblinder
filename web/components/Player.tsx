import { Replay } from 'vimond-replay';
import 'vimond-replay/index.css';
import ShakaVideoStreamer from 'vimond-replay/video-streamer/shaka-player';

type PlayerProps = {
    url: string;
};

export default function Player({ url }: PlayerProps) {
    return (
        <>
            <Replay
                source={url}
                initialPlaybackProps={{ isPaused: false }}
                options={{controls:{liveDisplayMode: 'live-offset'}}}
            >
                <ShakaVideoStreamer />
            </Replay>
        </>
    )
}