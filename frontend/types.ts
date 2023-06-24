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
  sourceLanguage: string;
  targetLanguage: string;
  sourceMedia: string;
  targetMedia: string;
};
