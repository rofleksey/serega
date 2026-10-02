import { useState } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import ViewKanbanOutlinedIcon from '@mui/icons-material/ViewKanbanOutlined';
import LogoutIcon from '@mui/icons-material/Logout';
import { AppBar, Box, IconButton, Stack, Toolbar, Tooltip, Typography } from '@mui/material';
import { useUser } from '@entities/session/auth';
import { RecoverableError } from '@shared/ui/async-states';

export function AppLayout() {
  const navigate = useNavigate();
  const { user, signOut } = useUser();
  const [signOutError, setSignOutError] = useState<unknown>();
  const [signingOut, setSigningOut] = useState(false);
  const leave = async () => {
    setSigningOut(true);
    setSignOutError(undefined);
    try {
      await signOut();
      navigate('/login', { replace: true });
    } catch (reason) {
      setSignOutError(reason);
    } finally {
      setSigningOut(false);
    }
  };
  return <Box sx={{ minHeight: '100dvh', bgcolor: 'background.default' }}>
    <AppBar position="sticky" elevation={0} sx={{ pt: 'env(safe-area-inset-top)', bgcolor: 'background.paper', borderBottom: 1, borderColor: 'divider' }}>
      <Toolbar variant="dense" sx={{ minHeight: 48, px: { xs: 1.5, sm: 2, lg: 2.5 }, width: '100%', gap: 1 }}>
        <Stack direction="row" spacing={1} alignItems="center" sx={{ flexGrow: 1, flexShrink: 0, minWidth: 0 }}><ViewKanbanOutlinedIcon color="primary" fontSize="small" /><Typography component="div" variant="h6">Serega</Typography><Typography variant="caption" color="text.secondary" sx={{ display: { xs: 'none', sm: 'block' } }}>Shared board</Typography></Stack>
        <Typography variant="body2" color="text.secondary" noWrap sx={{ minWidth: 0, maxWidth: { xs: 120, sm: 260 } }}>{user?.username}</Typography>
        <Tooltip title="Sign out"><Box component="span" sx={{ flexShrink: 0 }}><IconButton onClick={() => { void leave(); }} aria-label="Sign out" disabled={signingOut}><LogoutIcon /></IconButton></Box></Tooltip>
      </Toolbar>
    </AppBar>
    <Box component="main" sx={{ px: { xs: 1.5, sm: 2, lg: 2.5 }, py: 2, pb: 'calc(20px + env(safe-area-inset-bottom))', minWidth: 0 }}>
      <Stack spacing={1.5}>{signOutError !== undefined && <RecoverableError title="Sign out failed" error={signOutError} />}<Outlet /></Stack>
    </Box>
  </Box>;
}
