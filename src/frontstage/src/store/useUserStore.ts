import { create } from "zustand";
import { currentUser } from "@clerk/nextjs/server";

type UserState = {
  user: any;
  fetchCurrentUser: () => Promise<void>;
};

const useUserStore = create<UserState>((set) => ({
  user: null, // Initial state
  fetchCurrentUser: async () => {
    const user = await currentUser();
    set({ user });
  },
}));

export default useUserStore;
