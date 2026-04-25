import api from "./api";

export const projectService = {
  // Get projects for an organization
  getByOrganization: async (orgId) => {
    const response = await api.get(`/organizations/${orgId}/projects`);
    return response.data;
  },

  // Create new project
  create: async (orgId, name, description) => {
    const response = await api.post(`/organizations/${orgId}/projects`, {
      name,
      description,
    });
    return response.data;
  },

  // Delete project
  delete: async (projectId) => {
    const response = await api.delete(`/projects/${projectId}`);
    return response.data;
  },

  // Get project details
  getById: async (projectId) => {
    const response = await api.get(`/projects/${projectId}`);
    return response.data;
  },
};

export default projectService;
