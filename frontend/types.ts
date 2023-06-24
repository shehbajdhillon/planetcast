export interface Team {
  id: number;
  slug: string;
  name: string;
  teamType: 'PERSONAL' | 'TEAM';
  created: string;
};

export interface Project {
  id: number;
  teamId: number;
  title: string;
  sourceLanguage: SupportedLanguage;
  targetLanguage: string;
  sourceMedia: string;
  targetMedia: string;
};

export const SupportedLanguages = ["ENGLISH", "HINDI"] as const;

export type SupportedLanguage = typeof SupportedLanguages[number];

