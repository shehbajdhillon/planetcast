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
  sourceLanguage: SupportedLanguage;
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
};

export interface Segment {
  id: number;
  start: number;
  end: number;
  text: string;
};

export const SupportedLanguages = [
  "ENGLISH",
  "HINDI",
  "SPANISH",

  "GERMAN",
  "POLISH",
  "ITALIAN",

  "FRENCH",
  "PORTUGUESE",
] as const;

export type SupportedLanguage = typeof SupportedLanguages[number];

