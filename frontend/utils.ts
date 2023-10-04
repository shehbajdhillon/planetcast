export function formatTime(seconds: number): string {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainingSeconds = Math.floor(seconds % 60);

  const formattedHours = hours.toString().padStart(2, '0');
  const formattedMinutes = minutes.toString().padStart(2, '0');
  const formattedSeconds = remainingSeconds.toString().padStart(2, '0');

  if (hours == 0) {
    return `${formattedMinutes}:${formattedSeconds}`;
  }
  return `${formattedHours}:${formattedMinutes}:${formattedSeconds}`;
}

export function matchYoutubeUrl(url: string): boolean {
  const p = /^(?:https?:\/\/)?(?:m\.|www\.)?(?:youtu\.be\/|youtube\.com\/(?:embed\/|v\/|watch\?v=|watch\?.+&v=))((\w|-){11})(?:\S+)?$/;
  if (url.match(p)) {
    return true;
  }
  return false;
}

const videoRegexpList: RegExp[] = [
  /(?:v|embed|shorts|watch\?v)(?:=|\/)([^"&?\/=%]{11})/,
  /(?:=|\/)([^"&?\/=%]{11})/,
  /([^"&?/=%]{11})/,
];

// ExtractVideoID extracts the videoID from the given string
export function extractVideoID(videoID: string): boolean {
  if (videoID.includes("youtu") || /["?&/<%]/.test(videoID)) {
    for (const re of videoRegexpList) {
      if (re.test(videoID)) {
        const subs = videoID.match(re) || [];
        videoID = subs[1] || videoID;
      }
    }
  }

  if (/["?&/<%]/.test(videoID)) {
    return false
  }

  if (videoID.length < 10) {
    return false
  }

  return true;
}

