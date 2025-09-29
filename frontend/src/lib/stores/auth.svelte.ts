// Global reactive auth store using Svelte 5 runes
import type { User } from "$types/auth";

class AuthStore {
  user = $state<User | null>(null);
  isAuthenticated = $derived(this.user !== null);
  isLoading = $state(true);

  setUser(user: User | null) {
    this.user = user;
  }

  setLoading(loading: boolean) {
    this.isLoading = loading;
  }

  logout() {
    this.user = null;
  }
}

export const authStore = new AuthStore();
