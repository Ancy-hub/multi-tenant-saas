import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
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
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  IconButton,
} from "@mui/material";
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
} from "@mui/icons-material";
import { useAuth } from "../context/AuthContext";
import projectService from "../services/projectService";
import taskService from "../services/taskService";

const ProjectDetails = () => {
  const { id } = useParams();
  const { user } = useAuth(); // eslint-disable-line no-unused-vars
  const [project, setProject] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Task creation dialog
  const [createTaskDialogOpen, setCreateTaskDialogOpen] = useState(false);
  const [newTaskTitle, setNewTaskTitle] = useState("");
  const [newTaskDescription, setNewTaskDescription] = useState("");
  const [creatingTask, setCreatingTask] = useState(false);

  // Task editing dialog
  const [editTaskDialogOpen, setEditTaskDialogOpen] = useState(false);
  const [editingTask, setEditingTask] = useState(null);
  const [editTaskTitle, setEditTaskTitle] = useState("");
  const [editTaskDescription, setEditTaskDescription] = useState("");
  const [editTaskStatus, setEditTaskStatus] = useState("todo");
  const [updatingTask, setUpdatingTask] = useState(false);

  useEffect(() => {
    loadProjectData();
  }, [id]); // eslint-disable-line react-hooks/exhaustive-deps

  const loadProjectData = async () => {
    try {
      setLoading(true);
      const [projectData, tasksData] = await Promise.all([
        projectService.getById(id),
        taskService.getByProject(id),
      ]);
      setProject(projectData);
      setTasks(tasksData || []);
    } catch (error) {
      setError("Failed to load project data");
      console.error("Error loading project data:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTask = async () => {
    if (!newTaskTitle.trim()) return;

    try {
      setCreatingTask(true);
      await taskService.create(id, newTaskTitle, newTaskDescription);
      setCreateTaskDialogOpen(false);
      setNewTaskTitle("");
      setNewTaskDescription("");
      loadProjectData(); // Reload tasks
    } catch (error) {
      setError("Failed to create task");
      console.error("Error creating task:", error);
    } finally {
      setCreatingTask(false);
    }
  };

  const handleEditTask = (task) => {
    setEditingTask(task);
    setEditTaskTitle(task.title);
    setEditTaskDescription(task.description || "");
    setEditTaskStatus(task.status);
    setEditTaskDialogOpen(true);
  };

  const handleUpdateTask = async () => {
    if (!editTaskTitle.trim() || !editingTask) return;

    try {
      setUpdatingTask(true);
      await taskService.update(editingTask.id, {
        title: editTaskTitle,
        description: editTaskDescription,
        status: editTaskStatus,
      });
      setEditTaskDialogOpen(false);
      setEditingTask(null);
      loadProjectData(); // Reload tasks
    } catch (error) {
      setError("Failed to update task");
      console.error("Error updating task:", error);
    } finally {
      setUpdatingTask(false);
    }
  };

  const handleDeleteTask = async (taskId) => {
    if (!window.confirm("Are you sure you want to delete this task?")) return;

    try {
      await taskService.delete(taskId);
      loadProjectData(); // Reload tasks
    } catch (error) {
      setError("Failed to delete task");
      console.error("Error deleting task:", error);
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case "todo":
        return "default";
      case "in_progress":
        return "primary";
      case "done":
        return "success";
      default:
        return "default";
    }
  };

  const getStatusLabel = (status) => {
    switch (status) {
      case "todo":
        return "To Do";
      case "in_progress":
        return "In Progress";
      case "done":
        return "Done";
      default:
        return status;
    }
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Typography>Loading...</Typography>
      </Container>
    );
  }

  if (!project) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Alert severity="error">Project not found</Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4 }}>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {project.name}
        </Typography>
        {project.description && (
          <Typography variant="body1" color="text.secondary">
            {project.description}
          </Typography>
        )}
        <Typography variant="body2" sx={{ mt: 1 }}>
          Created: {new Date(project.created_at).toLocaleDateString()}
        </Typography>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Typography variant="h5" gutterBottom sx={{ mt: 4 }}>
        Tasks
      </Typography>

      <Grid container spacing={3}>
        {tasks.map((task) => (
          <Grid item xs={12} sm={6} md={4} key={task.id}>
            <Card
              sx={{ height: "100%", display: "flex", flexDirection: "column" }}
            >
              <CardContent sx={{ flexGrow: 1 }}>
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "flex-start",
                    mb: 1,
                  }}
                >
                  <Typography variant="h6" component="h3" gutterBottom>
                    {task.title}
                  </Typography>
                  <Chip
                    label={getStatusLabel(task.status)}
                    size="small"
                    color={getStatusColor(task.status)}
                  />
                </Box>
                {task.description && (
                  <Typography
                    variant="body2"
                    color="text.secondary"
                    sx={{ mb: 2 }}
                  >
                    {task.description}
                  </Typography>
                )}
                <Typography variant="body2">
                  Created: {new Date(task.created_at).toLocaleDateString()}
                </Typography>
              </CardContent>
              <CardActions>
                <IconButton
                  size="small"
                  onClick={() => handleEditTask(task)}
                  color="primary"
                >
                  <EditIcon />
                </IconButton>
                <IconButton
                  size="small"
                  onClick={() => handleDeleteTask(task.id)}
                  color="error"
                >
                  <DeleteIcon />
                </IconButton>
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
            Create your first task to get started
          </Typography>
        </Box>
      )}

      {/* Floating Action Button for creating new task */}
      <Fab
        color="primary"
        aria-label="add task"
        sx={{ position: "fixed", bottom: 16, right: 16 }}
        onClick={() => setCreateTaskDialogOpen(true)}
      >
        <AddIcon />
      </Fab>

      {/* Create Task Dialog */}
      <Dialog
        open={createTaskDialogOpen}
        onClose={() => setCreateTaskDialogOpen(false)}
      >
        <DialogTitle>Create New Task</DialogTitle>
        <DialogContent>
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
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateTaskDialogOpen(false)}>Cancel</Button>
          <Button
            onClick={handleCreateTask}
            variant="contained"
            disabled={!newTaskTitle.trim() || creatingTask}
          >
            {creatingTask ? "Creating..." : "Create"}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Edit Task Dialog */}
      <Dialog
        open={editTaskDialogOpen}
        onClose={() => setEditTaskDialogOpen(false)}
      >
        <DialogTitle>Edit Task</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Task Title"
            fullWidth
            variant="outlined"
            value={editTaskTitle}
            onChange={(e) => setEditTaskTitle(e.target.value)}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            label="Description (optional)"
            fullWidth
            multiline
            rows={3}
            variant="outlined"
            value={editTaskDescription}
            onChange={(e) => setEditTaskDescription(e.target.value)}
            sx={{ mb: 2 }}
          />
          <FormControl fullWidth margin="dense">
            <InputLabel>Status</InputLabel>
            <Select
              value={editTaskStatus}
              label="Status"
              onChange={(e) => setEditTaskStatus(e.target.value)}
            >
              <MenuItem value="todo">To Do</MenuItem>
              <MenuItem value="in_progress">In Progress</MenuItem>
              <MenuItem value="done">Done</MenuItem>
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditTaskDialogOpen(false)}>Cancel</Button>
          <Button
            onClick={handleUpdateTask}
            variant="contained"
            disabled={!editTaskTitle.trim() || updatingTask}
          >
            {updatingTask ? "Updating..." : "Update"}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default ProjectDetails;
