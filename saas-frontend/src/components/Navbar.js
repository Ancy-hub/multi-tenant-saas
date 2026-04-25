import React from "react";
import { useNavigate, useLocation } from "react-router-dom";
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Avatar,
  Container,
} from "@mui/material";
import { Dashboard as DashboardIcon, Logout as LogoutIcon } from "@mui/icons-material";
import { useAuth } from "../context/AuthContext";

const Navbar = () => {
  const { user, logout, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [anchorEl, setAnchorEl] = React.useState(null);

  const handleMenu = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    navigate("/login");
    handleClose();
  };

  if (
    !isAuthenticated ||
    location.pathname === "/login" ||
    location.pathname === "/register"
  ) {
    return null;
  }

  return (
    <AppBar position="sticky" sx={{ top: 0, zIndex: 1100, mb: 4 }}>
      <Container maxWidth="xl">
        <Toolbar disableGutters sx={{ justifyContent: 'space-between', py: 1 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <Box 
              sx={{ 
                width: 40, 
                height: 40, 
                borderRadius: 2, 
                background: 'linear-gradient(45deg, #3b82f6, #8b5cf6)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                boxShadow: '0 0 20px rgba(59, 130, 246, 0.4)'
              }}
            >
              <Typography variant="h6" fontWeight="bold" color="white">S</Typography>
            </Box>
            <Typography variant="h5" component="div" sx={{ fontWeight: 700, letterSpacing: '-0.5px' }}>
              SaaS<span style={{ color: '#3b82f6' }}>Core</span>
            </Typography>
          </Box>

          <Box sx={{ display: "flex", alignItems: "center", gap: 3 }}>
            <Button 
              color="inherit" 
              onClick={() => navigate("/dashboard")}
              startIcon={<DashboardIcon />}
              sx={{ 
                opacity: location.pathname === '/dashboard' ? 1 : 0.7,
                '&:hover': { opacity: 1 }
              }}
            >
              Dashboard
            </Button>

            {user && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, pl: 2, borderLeft: '1px solid rgba(255,255,255,0.1)' }}>
                <Box sx={{ textAlign: 'right', display: { xs: 'none', sm: 'block' } }}>
                  <Typography variant="body2" sx={{ fontWeight: 600, color: 'text.primary' }}>
                    {user.name}
                  </Typography>
                  <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                    Admin
                  </Typography>
                </Box>
                <IconButton
                  onClick={handleMenu}
                  sx={{ p: 0, border: '2px solid rgba(59, 130, 246, 0.5)' }}
                >
                  <Avatar sx={{ bgcolor: '#2563eb', width: 36, height: 36 }}>
                    {user.name ? user.name.charAt(0).toUpperCase() : 'U'}
                  </Avatar>
                </IconButton>
                <Menu
                  id="menu-appbar"
                  anchorEl={anchorEl}
                  anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
                  keepMounted
                  transformOrigin={{ vertical: "top", horizontal: "right" }}
                  open={Boolean(anchorEl)}
                  onClose={handleClose}
                  PaperProps={{
                    sx: {
                      mt: 1.5,
                      minWidth: 180,
                      background: 'rgba(18, 24, 38, 0.95)',
                      backdropFilter: 'blur(20px)',
                    }
                  }}
                >
                  <MenuItem onClick={handleLogout} sx={{ color: '#ef4444', gap: 1 }}>
                    <LogoutIcon fontSize="small" /> Logout
                  </MenuItem>
                </Menu>
              </Box>
            )}
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
};

export default Navbar;
