import { useState } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import ViewKanbanOutlinedIcon from '@mui/icons-material/ViewKanbanOutlined';
import LogoutIcon from '@mui/icons-material/Logout';
import { AppBar, Box, Chip, IconButton, Stack, Toolbar, Tooltip, Typography } from '@mui/material';
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
    <AppBar position="sticky" elevation={0} sx={{ pt: 'env(safe-area-inset-top)', bgcolor: 'rgba(17, 21, 27, .94)', backdropFilter: 'blur(14px)', borderBottom: 1, borderColor: 'divider' }}>
      <Toolbar sx={{ minHeight: 64, px: { xs: 2, sm: 3, lg: 4 }, width: '100%', maxWidth: 1600, mx: 'auto', gap: 1 }}>
        <Stack direction="row" spacing={1.25} alignItems="center" sx={{ flexGrow: 1, flexShrink: 0, minWidth: 0 }}><ViewKanbanOutlinedIcon color="primary" /><Typography component="div" variant="h6">Serega</Typography><Chip label="Shared board" size="small" variant="outlined" sx={{ display: { xs: 'none', sm: 'flex' }, ml: 1 }} /></Stack>
        <Typography variant="body2" color="text.secondary" noWrap sx={{ minWidth: 0, maxWidth: { xs: 120, sm: 260 } }}>{user?.username}</Typography>
        <Tooltip title="Sign out"><Box component="span" sx={{ flexShrink: 0 }}><IconButton onClick={() => { void leave(); }} aria-label="Sign out" disabled={signingOut}><LogoutIcon /></IconButton></Box></Tooltip>
      </Toolbar>
    </AppBar>
    <Box component="main" sx={{ px: { xs: 2, sm: 3, lg: 4 }, py: { xs: 3, md: 4 }, pb: 'calc(32px + env(safe-area-inset-bottom))', maxWidth: 1600, mx: 'auto' }}>
      <Stack spacing={3}>{signOutError !== undefined && <RecoverableError title="Sign out failed" error={signOutError} />}<Outlet /></Stack>
    </Box>
  </Box>;
}
