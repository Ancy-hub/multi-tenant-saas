import api from "./api";

export const authService = {
  // Login user
  login: async (email, password) => {
    const response = await api.post(`/login?t=${Date.now()}`, {
      email,
      password,
    });
    return response.data;
  },

  // Register new user
  register: async (name, email, password) => {
    const response = await api.post(`/users?t=${Date.now()}`, {
      name,
      email,
      password,
    });
    return response.data;
  },

  // Get current user info (if needed)
  getCurrentUser: async () => {
    // This would require a backend endpoint to get current user
    // For now, we'll decode the token on frontend
    return null;
  },
};

export default authService;
