import api from "./api";

export const organizationService = {
  // Get all organizations
  getAll: async () => {
    const response = await api.get("/organizations");
    return response.data;
  },

  // Get organization by ID
  getById: async (id) => {
    const response = await api.get(`/organizations/${id}`);
    return response.data;
  },

  // Create new organization
  create: async (name, description) => {
    const response = await api.post("/organizations", { name, description });
    return response.data;
  },

  // Update organization
  update: async (id, name) => {
    const response = await api.patch(`/organizations/${id}`, { name });
    return response.data;
  },

  // Get organization members
  getMembers: async (orgId) => {
    const response = await api.get(`/organizations/${orgId}/members`);
    return response.data;
  },

  // Add user to organization
  addMember: async (orgId, userId, role) => {
    const response = await api.post(`/organizations/${orgId}/members`, {
      user_id: userId,
      role: role,
    });
    return response.data;
  },

  // Remove member from organization
  removeMember: async (orgId, userId) => {
    const response = await api.delete(
      `/organizations/${orgId}/members/${userId}`,
    );
    return response.data;
  },

  // Update member role
  updateMemberRole: async (orgId, userId, role) => {
    const response = await api.patch(
      `/organizations/${orgId}/members/${userId}`,
      { role },
    );
    return response.data;
  },
};

export default organizationService;
