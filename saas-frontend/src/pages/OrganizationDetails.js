import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Container,
  Typography,
  Box,
  Grid,
  Card,
  CardContent,
  CardActions,
  Button,
  Fab,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  Chip,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  IconButton,
  Tabs,
  Tab,
} from "@mui/material";
import {
  Add as AddIcon,
  Delete as DeleteIcon,
  PersonAdd as PersonAddIcon,
} from "@mui/icons-material";
import { useAuth } from "../context/AuthContext";
import organizationService from "../services/organizationService";
import projectService from "../services/projectService";
import taskService from "../services/taskService";

const OrganizationDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [organization, setOrganization] = useState(null);
  const [members, setMembers] = useState([]);
  const [projects, setProjects] = useState([]);
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [tabValue, setTabValue] = useState(0);

  // Project creation dialog
  const [createProjectDialogOpen, setCreateProjectDialogOpen] = useState(false);
  const [newProjectName, setNewProjectName] = useState("");
  const [newProjectDescription, setNewProjectDescription] = useState("");
  const [creatingProject, setCreatingProject] = useState(false);

  // Member addition dialog
  const [addMemberDialogOpen, setAddMemberDialogOpen] = useState(false);
  const [newMemberEmail, setNewMemberEmail] = useState("");
  const [addingMember, setAddingMember] = useState(false);

  // Task creation dialog
  const [createTaskDialogOpen, setCreateTaskDialogOpen] = useState(false);
  const [newTaskTitle, setNewTaskTitle] = useState("");
  const [newTaskDescription, setNewTaskDescription] = useState("");
  const [newTaskProjectId, setNewTaskProjectId] = useState("");
  const [creatingTask, setCreatingTask] = useState(false);

  useEffect(() => {
    loadOrganizationData();
  }, [id]); // eslint-disable-line react-hooks/exhaustive-deps

  const loadOrganizationData = async () => {
    try {
      setLoading(true);
      const [orgData, membersData, projectsData, tasksData] = await Promise.all([
        organizationService.getById(id),
        organizationService.getMembers(id),
        projectService.getByOrganization(id),
        taskService.getByOrganization(id),
      ]);
      setOrganization(orgData);
      setMembers(membersData || []);
      setProjects(projectsData || []);
      setTasks(tasksData || []);
    } catch (error) {
      setError("Failed to load organization data");
      console.error("Error loading organization data:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateProject = async () => {
    if (!newProjectName.trim()) return;

    try {
      setCreatingProject(true);
      await projectService.create(id, newProjectName, newProjectDescription);
      setCreateProjectDialogOpen(false);
      setNewProjectName("");
      setNewProjectDescription("");
      loadOrganizationData(); // Reload projects
    } catch (error) {
      setError("Failed to create project");
      console.error("Error creating project:", error);
    } finally {
      setCreatingProject(false);
    }
  };

  const handleAddMember = async () => {
    if (!newMemberEmail.trim()) return;

    try {
      setAddingMember(true);
      await organizationService.addMember(id, newMemberEmail);
      setAddMemberDialogOpen(false);
      setNewMemberEmail("");
      loadOrganizationData(); // Reload members
    } catch (error) {
      setError("Failed to add member");
      console.error("Error adding member:", error);
    } finally {
      setAddingMember(false);
    }
  };

  const handleRemoveMember = async (userId) => {
    try {
      await organizationService.removeMember(id, userId);
      loadOrganizationData(); // Reload members
    } catch (error) {
      setError("Failed to remove member");
      console.error("Error removing member:", error);
    }
  };

  const handleCreateTask = async () => {
    if (!newTaskTitle.trim() || !newTaskProjectId) return;

    try {
      setCreatingTask(true);
      await taskService.create(newTaskProjectId, newTaskTitle, newTaskDescription);
      setCreateTaskDialogOpen(false);
      setNewTaskTitle("");
      setNewTaskDescription("");
      setNewTaskProjectId("");
      loadOrganizationData(); // Reload tasks
    } catch (error) {
      setError("Failed to create task");
      console.error("Error creating task:", error);
    } finally {
      setCreatingTask(false);
    }
  };

  const handleProjectClick = (projectId) => {
    navigate(`/projects/${projectId}`);
  };

  const isAdmin =
    members.find((member) => member.user_id === user?.id)?.role === "admin";

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Typography>Loading...</Typography>
      </Container>
    );
  }

  if (!organization) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Alert severity="error">Organization not found</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4 }}>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {organization.name}
        </Typography>
        {organization.description && (
          <Typography variant="body1" color="text.secondary">
            {organization.description}
          </Typography>
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
        <Tabs
          value={tabValue}
          onChange={(e, newValue) => setTabValue(newValue)}
        >
          <Tab label="Projects" />
          <Tab label="Members" />
          <Tab label="All Tasks" />
        </Tabs>
      </Box>

      {/* Projects Tab */}
      {tabValue === 0 && (
        <Box sx={{ mt: 3 }}>
          <Grid container spacing={3}>
            {projects.map((project) => (
              <Grid item xs={12} sm={6} md={4} key={project.id}>
                <Card
                  sx={{
                    height: "100%",
                    display: "flex",
                    flexDirection: "column",
                  }}
                >
                  <CardContent sx={{ flexGrow: 1 }}>
                    <Typography variant="h5" component="h2" gutterBottom>
                      {project.name}
                    </Typography>
                    {project.description && (
                      <Typography variant="body2" color="text.secondary">
                        {project.description}
                      </Typography>
                    )}
                    <Typography variant="body2" sx={{ mt: 1 }}>
                      Created:{" "}
                      {new Date(project.created_at).toLocaleDateString()}
                    </Typography>
                  </CardContent>
                  <CardActions>
                    <Button
                      size="small"
                      color="primary"
                      onClick={() => handleProjectClick(project.id)}
                    >
                      View Project
                    </Button>
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </Grid>

          {projects.length === 0 && (
            <Box sx={{ textAlign: "center", mt: 8 }}>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                No projects yet
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Create your first project to get started
              </Typography>
            </Box>
          )}
        </Box>
      )}

      {/* Members Tab */}
      {tabValue === 1 && (
        <Box sx={{ mt: 3 }}>
          <List>
            {members.map((member) => (
              <ListItem key={member.user_id}>
                <ListItemText
                  primary={member.name || member.email}
                  secondary={
                    <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                      <Typography variant="body2">
                        {member.email}
                      </Typography>
                      <Chip
                        label={member.role}
                        size="small"
                        color={member.role === "admin" ? "primary" : "default"}
                      />
                    </Box>
                  }
                />
                {isAdmin && member.user_id !== user?.id && (
                  <ListItemSecondaryAction>
                    <IconButton
                      edge="end"
                      aria-label="remove member"
                      onClick={() => handleRemoveMember(member.user_id)}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </ListItemSecondaryAction>
                )}
              </ListItem>
            ))}
          </List>
        </Box>
      )}

      {/* Tasks Tab */}
      {tabValue === 2 && (
        <Box sx={{ mt: 3 }}>
          <Grid container spacing={3}>
            {tasks.map((task) => (
              <Grid item xs={12} sm={6} md={4} key={task.id}>
                <Card
                  sx={{
                    height: "100%",
                    display: "flex",
                    flexDirection: "column",
                  }}
                >
                  <CardContent sx={{ flexGrow: 1 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                      <Typography variant="h6" component="h3">
                        {task.title}
                      </Typography>
                      <Chip 
                        label={task.status} 
                        size="small" 
                        color={task.status === "done" ? "success" : task.status === "in_progress" ? "primary" : "default"} 
                      />
                    </Box>
                    {task.description && (
                      <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                        {task.description}
                      </Typography>
                    )}
                    <Typography variant="caption" display="block" color="text.secondary">
                      Created: {new Date(task.created_at).toLocaleDateString()}
                    </Typography>
                  </CardContent>
                  <CardActions>
                    <Button
                      size="small"
                      color="primary"
                      onClick={() => navigate(`/projects/${task.project_id}`)}
                    >
                      View in Project
                    </Button>
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </Grid>
          {tasks.length === 0 && (
            <Box sx={{ textAlign: "center", mt: 8 }}>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                No tasks yet
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Tasks created in this organization's projects will appear here.
              </Typography>
            </Box>
          )}
        </Box>
      )}

      {/* Floating Action Buttons */}
      {isAdmin && (
        <>
          <Fab
            color="primary"
            aria-label="add project"
            sx={{ position: "fixed", bottom: 16, right: 16 }}
            onClick={() => setCreateProjectDialogOpen(true)}
          >
            <AddIcon />
          </Fab>

          {tabValue === 1 && (
            <Fab
              color="secondary"
              aria-label="add member"
              sx={{ position: "fixed", bottom: 16, right: 80 }}
              onClick={() => setAddMemberDialogOpen(true)}
            >
              <PersonAddIcon />
            </Fab>
          )}

          {tabValue === 2 && (
            <Fab
              color="secondary"
              aria-label="add task"
              sx={{ position: "fixed", bottom: 16, right: 80 }}
              onClick={() => {
                if (projects.length > 0 && !newTaskProjectId) {
                  setNewTaskProjectId(projects[0].id);
                }
                setCreateTaskDialogOpen(true);
              }}
            >
              <AddIcon />
            </Fab>
          )}
        </>
      )}

      {/* Create Project Dialog */}
      <Dialog
        open={createProjectDialogOpen}
        onClose={() => setCreateProjectDialogOpen(false)}
      >
        <DialogTitle>Create New Project</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Project Name"
            fullWidth
            variant="outlined"
            value={newProjectName}
            onChange={(e) => setNewProjectName(e.target.value)}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Description (optional)"
            fullWidth
            multiline
            rows={3}
            variant="outlined"
            value={newProjectDescription}
            onChange={(e) => setNewProjectDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateProjectDialogOpen(false)}>
            Cancel
          </Button>
          <Button
            onClick={handleCreateProject}
            variant="contained"
            disabled={!newProjectName.trim() || creatingProject}
          >
            {creatingProject ? "Creating..." : "Create"}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add Member Dialog */}
      <Dialog
        open={addMemberDialogOpen}
        onClose={() => setAddMemberDialogOpen(false)}
      >
        <DialogTitle>Add Member</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="User Email"
            fullWidth
            variant="outlined"
            value={newMemberEmail}
            onChange={(e) => setNewMemberEmail(e.target.value)}
            helperText="Enter the email address of the user to add"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAddMemberDialogOpen(false)}>Cancel</Button>
          <Button
            onClick={handleAddMember}
            variant="contained"
            disabled={!newMemberEmail.trim() || addingMember}
          >
            {addingMember ? "Adding..." : "Add Member"}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Create Task Dialog */}
      <Dialog
        open={createTaskDialogOpen}
        onClose={() => setCreateTaskDialogOpen(false)}
      >
        <DialogTitle>Create New Task</DialogTitle>
        <DialogContent>
          {projects.length === 0 ? (
            <Alert severity="warning" sx={{ mt: 2 }}>
              You need to create a project first before adding tasks.
            </Alert>
          ) : (
            <>
              <TextField
                select
                margin="dense"
                label="Project"
                fullWidth
                variant="outlined"
                value={newTaskProjectId}
                onChange={(e) => setNewTaskProjectId(e.target.value)}
                sx={{ mb: 2, mt: 1 }}
                SelectProps={{
                  native: true,
                }}
              >
                <option value="" disabled>Select a project</option>
                {projects.map((project) => (
                  <option key={project.id} value={project.id}>
                    {project.name}
                  </option>
                ))}
              </TextField>
              <TextField
                autoFocus
                margin="dense"
                label="Task Title"
                fullWidth
                variant="outlined"
                value={newTaskTitle}
                onChange={(e) => setNewTaskTitle(e.target.value)}
                sx={{ mb: 2 }}
              />
              <TextField
                margin="dense"
                label="Description (optional)"
                fullWidth
                multiline
                rows={3}
                variant="outlined"
                value={newTaskDescription}
                onChange={(e) => setNewTaskDescription(e.target.value)}
              />
            </>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateTaskDialogOpen(false)}>
            Cancel
          </Button>
          <Button
            onClick={handleCreateTask}
            variant="contained"
            disabled={projects.length === 0 || !newTaskTitle.trim() || !newTaskProjectId || creatingTask}
          >
            {creatingTask ? "Creating..." : "Create"}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default OrganizationDetails;
