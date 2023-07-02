import { create } from "zustand";

interface VideoSeekStore {
  currentSeek: number;
  setCurrentSeek: (seek: number) => void;
};

export const useVideoSeekStore = create<VideoSeekStore>()((set) => ({
  currentSeek: 0,
  setCurrentSeek: (seek: number) => set(() => ({ currentSeek: seek })),
}));

