export interface Team {
  id: number;
  slug: string;
  name: string;
  teamType: 'PERSONAL' | 'TEAM';
  created: string;
  projects: Project[];
};

export interface Project {
  id: number;
  teamId: number;
  title: string;
  sourceMedia: string;
  transformations: Transformation[];
};

export interface Transformation {
  id: number;
  projectId: number;
  targetMedia: string;
  targetLanguage: SupportedLanguage;
  transcript: string;
  isSource: boolean;
  status: string;
  progress: number;
};

export interface Segment {
  id: number;
  start: number;
  end: number;
  text: string;
};

export const SupportedLanguages = [
  "ENGLISH",
  "GERMAN",
  "FRENCH",
  "SPANISH",
  "ITALIAN",
  "PORTUGUESE",
  "POLISH",
] as const;

export type SupportedLanguage = typeof SupportedLanguages[number];

export const UploadOptions = [
  "FILE_UPLOAD",
  "YOUTUBE_LINK",
] as const;

export type UploadOption = typeof UploadOptions[number];

