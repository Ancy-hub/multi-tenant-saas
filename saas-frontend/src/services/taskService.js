import api from "./api";

export const taskService = {
  // Get tasks for a project
  getByProject: async (projectId) => {
    const response = await api.get(`/projects/${projectId}/tasks`);
    return response.data;
  },

  // Get tasks for an organization
  getByOrganization: async (orgId) => {
    const response = await api.get(`/organizations/${orgId}/tasks`);
    return response.data;
  },

  // Create new task
  create: async (projectId, title, description, status = "todo") => {
    const response = await api.post(`/projects/${projectId}/tasks`, {
      title,
      description,
      status,
    });
    return response.data;
  },

  // Update task
  update: async (taskId, updates) => {
    const response = await api.patch(`/tasks/${taskId}`, updates);
    return response.data;
  },

  // Delete task
  delete: async (taskId) => {
    const response = await api.delete(`/tasks/${taskId}`);
    return response.data;
  },
};

export default taskService;
