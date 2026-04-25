import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  Container, Typography, Box, Grid, Card, CardContent, CardActions,
  Button, Fab, Dialog, DialogTitle, DialogContent, DialogActions,
  TextField, Alert, Skeleton, Chip, Divider, IconButton
} from "@mui/material";
import { 
  Add as AddIcon, 
  Business as BusinessIcon,
  ArrowForward as ArrowForwardIcon,
  Settings as SettingsIcon
} from "@mui/icons-material";
import { useAuth } from "../context/AuthContext";
import organizationService from "../services/organizationService";

const Dashboard = () => {
  const [organizations, setOrganizations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [newOrgName, setNewOrgName] = useState("");
  const [newOrgDescription, setNewOrgDescription] = useState("");
  const [creating, setCreating] = useState(false);
  const { user } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    loadOrganizations();
  }, []);

  const loadOrganizations = async () => {
    try {
      setLoading(true);
      const data = await organizationService.getAll();
      setOrganizations(data || []);
    } catch (error) {
      setError("Failed to load organizations");
      console.error("Error loading organizations:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateOrganization = async () => {
    if (!newOrgName.trim()) return;

    try {
      setCreating(true);
      await organizationService.create(newOrgName, newOrgDescription);
      setCreateDialogOpen(false);
      setNewOrgName("");
      setNewOrgDescription("");
      loadOrganizations(); // Reload organizations
    } catch (error) {
      setError("Failed to create organization");
      console.error("Error creating organization:", error);
    } finally {
      setCreating(false);
    }
  };

  const handleOrganizationClick = (orgId) => {
    navigate(`/organizations/${orgId}`);
  };

  return (
    <Container maxWidth="xl" sx={{ mt: 2, mb: 8 }}>
      <Box sx={{ mb: 6, display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end' }}>
        <Box>
          <Typography variant="h3" component="h1" gutterBottom sx={{ fontWeight: 700, background: 'linear-gradient(45deg, #60a5fa, #a78bfa)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
            Welcome back, {user?.name}!
          </Typography>
          <Typography variant="h6" color="text.secondary" sx={{ fontWeight: 400 }}>
            Manage your organizations and projects
          </Typography>
        </Box>
        <Button 
          variant="contained" 
          startIcon={<AddIcon />}
          onClick={() => setCreateDialogOpen(true)}
          sx={{ display: { xs: 'none', md: 'flex' } }}
        >
          New Organization
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 4, borderRadius: 2 }}>
          {error}
        </Alert>
      )}

      <Grid container spacing={4}>
        {loading ? (
          // Skeleton Loaders
          [1, 2, 3].map((item) => (
            <Grid item xs={12} sm={6} md={4} key={item}>
              <Card sx={{ height: 220, display: "flex", flexDirection: "column" }}>
                <CardContent sx={{ flexGrow: 1 }}>
                  <Skeleton animation="wave" height={40} width="80%" sx={{ mb: 2 }} />
                  <Skeleton animation="wave" height={20} width="100%" />
                  <Skeleton animation="wave" height={20} width="60%" />
                </CardContent>
              </Card>
            </Grid>
          ))
        ) : (
          organizations.map((org) => (
            <Grid item xs={12} sm={6} md={4} key={org.id}>
              <Card sx={{ height: "100%", display: "flex", flexDirection: "column", position: 'relative', overflow: 'hidden' }}>
                <Box sx={{ position: 'absolute', top: -20, right: -20, opacity: 0.05, transform: 'rotate(15deg)' }}>
                  <BusinessIcon sx={{ fontSize: 140 }} />
                </Box>
                <CardContent sx={{ flexGrow: 1, zIndex: 1 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                      <Box sx={{ p: 1, borderRadius: 2, bgcolor: 'rgba(59, 130, 246, 0.1)', color: '#60a5fa' }}>
                        <BusinessIcon />
                      </Box>
                      <Typography variant="h5" component="h2" sx={{ fontWeight: 600 }}>
                        {org.name}
                      </Typography>
                    </Box>
                    <IconButton size="small" sx={{ color: 'text.secondary' }}>
                      <SettingsIcon fontSize="small" />
                    </IconButton>
                  </Box>
                  <Typography variant="body1" color="text.secondary" sx={{ mb: 3, minHeight: 48, display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>
                    {org.description || "No description provided."}
                  </Typography>
                  <Divider sx={{ mb: 2, borderColor: 'rgba(255,255,255,0.05)' }} />
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Chip 
                      label="Active" 
                      size="small" 
                      sx={{ bgcolor: 'rgba(16, 185, 129, 0.1)', color: '#34d399', fontWeight: 600 }} 
                    />
                    <Typography variant="caption" color="text.secondary">
                      Created: {new Date(org.created_at).toLocaleDateString()}
                    </Typography>
                  </Box>
                </CardContent>
                <CardActions sx={{ p: 2, pt: 0, zIndex: 1 }}>
                  <Button
                    fullWidth
                    variant="outlined"
                    color="primary"
                    endIcon={<ArrowForwardIcon />}
                    onClick={() => handleOrganizationClick(org.id)}
                    sx={{ borderRadius: 2, borderColor: 'rgba(59, 130, 246, 0.3)', '&:hover': { borderColor: '#3b82f6', bgcolor: 'rgba(59, 130, 246, 0.05)' } }}
                  >
                    View Workspace
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          ))
        )}
      </Grid>

      {organizations.length === 0 && !loading && (
        <Box sx={{ textAlign: "center", mt: 10, p: 6, borderRadius: 4, bgcolor: 'rgba(255,255,255,0.02)', border: '1px dashed rgba(255,255,255,0.1)' }}>
          <BusinessIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2, opacity: 0.5 }} />
          <Typography variant="h5" gutterBottom sx={{ fontWeight: 600 }}>
            No organizations yet
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 4, maxWidth: 400, mx: 'auto' }}>
            Get started by creating your first organization. Organizations help you group projects and manage team access.
          </Typography>
          <Button 
            variant="contained" 
            size="large"
            startIcon={<AddIcon />}
            onClick={() => setCreateDialogOpen(true)}
          >
            Create Organization
          </Button>
        </Box>
      )}

      {/* Floating Action Button for mobile */}
      <Fab
        color="primary"
        aria-label="add organization"
        sx={{ position: "fixed", bottom: 24, right: 24, display: { md: 'none' } }}
        onClick={() => setCreateDialogOpen(true)}
      >
        <AddIcon />
      </Fab>

      {/* Create Organization Dialog */}
      <Dialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle sx={{ fontWeight: 600, pb: 1 }}>Create New Organization</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            An organization is a top-level entity that contains projects, teams, and settings.
          </Typography>
          <TextField
            autoFocus
            margin="dense"
            label="Organization Name"
            fullWidth
            variant="outlined"
            value={newOrgName}
            onChange={(e) => setNewOrgName(e.target.value)}
            sx={{ mb: 3 }}
          />
          <TextField
            margin="dense"
            label="Description (optional)"
            fullWidth
            multiline
            rows={3}
            variant="outlined"
            value={newOrgDescription}
            onChange={(e) => setNewOrgDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions sx={{ p: 3, pt: 0 }}>
          <Button onClick={() => setCreateDialogOpen(false)} sx={{ color: 'text.secondary' }}>Cancel</Button>
          <Button
            onClick={handleCreateOrganization}
            variant="contained"
            disabled={!newOrgName.trim() || creating}
            sx={{ px: 4 }}
          >
            {creating ? "Creating..." : "Create Organization"}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Dashboard;
