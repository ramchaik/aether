import { create } from "zustand";

export interface Project {
  id: string;
  name: string;
  slug: string;
  domain: string;
  repositoryUrl: string;
  customDomain: string | null;
  buildCommand: string;
  createdAt: string;
  updatedAt: string;
  userId: string;
  status: "LIVE" | "NOT_LIVE" | "DEPLOYING";
}

interface ProjectStore {
  projects: Project[];
  setProjects: (projects: Project[]) => void;
}

const useProjectStore = create<ProjectStore>((set) => ({
  projects: [],
  setProjects: (projects) => set({ projects }),
}));

export default useProjectStore;
