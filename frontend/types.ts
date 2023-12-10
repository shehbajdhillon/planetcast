export interface Team {
  id: number;
  slug: string;
  name: string;
  teamType: 'PERSONAL' | 'TEAM';
  created: string;
  projects: Project[];
  subscriptionPlans: SubscriptionPlan[];
  members: TeamMember[];
  invitees: TeamInvite[];
};

export interface TeamMember {
  membershipType: string;
  user: Userinfo;
};

export interface TeamInvite {
  inviteeEmail: string;
  slug: string;
};

interface Userinfo {
  id: number;
  email: string;
  fullName: string;
}

export interface Project {
  id: number;
  teamId: number;
  title: string;
  sourceMedia: string;
  dubbingCreditsRequired?: number;
  transformations: Transformation[];
};

export interface SubscriptionPlan {
  id: number;
  teamId: number;
  stripeSubscriptionId?: string;
  remainingCredits: number;
  subscriptionData?: SubscriptionData;
}

export interface SubscriptionData {
  currentPeriodStart: string,
  currentPeriodEnd: string,
  status: string,
  interval: string,
  planName: string,
  costInUsd: number,
  lastFourCardDigits: string,
}

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
  "JAPANESE",
  "CHINESE",
  "GERMAN",

  "HINDI",
  "FRENCH",
  "KOREAN",
  "PORTUGUESE",

  "ITALIAN",
  "SPANISH",
  "INDONESIAN",
  "DUTCH",

  "TURKISH",
  "FILIPINO",
  "POLISH",
  "SWEDISH",

  "BULGARIAN",
  "ROMANIAN",
  "ARABIC",
  "CZECH",

  "GREEK",
  "FINNISH",
  "CROATIAN",
  "MALAY",

  "SLOVAK",
  "DANISH",
  "TAMIL",
  "UKRAINIAN",
] as const;

export type SupportedLanguage = typeof SupportedLanguages[number];

export const UploadOptions = [
  "FILE_UPLOAD",
  "YOUTUBE_LINK",
] as const;

export type UploadOption = typeof UploadOptions[number];

