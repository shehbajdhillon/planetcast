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
};

export const SupportedLanguages = ["ENGLISH", "HINDI"] as const;

export type SupportedLanguage = typeof SupportedLanguages[number];

