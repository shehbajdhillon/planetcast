import { CSSProperties, useEffect, useRef } from 'react';
import videojs from 'video.js';
import Player from 'video.js/dist/types/player';
import 'video.js/dist/video-js.css';

interface VideoJSProps {
  onReady: (player: Player) => void;
  options: Record<string, any>;
  style?: CSSProperties;
}

const VideoJS: React.FC<VideoJSProps> = (props: any) => {
  const videoRef = useRef<HTMLDivElement>(null);
  const playerRef = useRef<Player | null>(null);
  const { options, onReady, style } = props;

  useEffect(() => {
    if (!playerRef.current) {
      const videoElement = document.createElement('video-js');
      videoElement.classList.add('vjs-big-play-centered');
      videoRef.current?.appendChild(videoElement);

      const player = (playerRef.current = videojs(videoElement, options, () => {
        onReady && onReady(player);
      }));
    } else {
      const player = playerRef.current;
      player.autoplay(options.autoplay);
      player.src(options.sources);
    }
  }, [options, videoRef, onReady]);

  useEffect(() => {
    const player = playerRef.current;
    return () => {
      if (player && !player.isDisposed()) {
        player.dispose();
        playerRef.current = null;
      }
    };
  }, [playerRef]);

  return (
    <div data-vjs-player style={{ width: '100%', ...style }}>
      <div ref={videoRef} style={{ ...style }}/>
    </div>
  );
};

interface VideoPlayerProps {
  src: string;
  onTimeUpdate?: (time: number) => void;
  style?: CSSProperties;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ src, onTimeUpdate, style }) => {
  const playerRef = useRef<Player | null>(null);
  const videoJsOptions = {
    autoplay: false,
    controls: true,
    responsive: true,
    fluid: true,
    sources: [
      {
        src: src,
        type: 'video/mp4',
      },
    ],
    playbackRates: [0.5, 1, 1.5, 2],
  };

  const handlePlayerReady = (player: Player) => {
    playerRef.current = player;
    player.aspectRatio('16:9');

    // You can handle player events here, for example:
    player.on('waiting', () => {
      videojs.log('player is waiting');
    });

    player.on('timeupdate', () => {
      onTimeUpdate?.(player.currentTime());
    });

    player.on('dispose', () => {
    });
  };

  return <VideoJS style={style} options={videoJsOptions} onReady={handlePlayerReady} />;
};

export default VideoPlayer;
