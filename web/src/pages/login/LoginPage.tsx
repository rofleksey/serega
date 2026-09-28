import { useState, type FormEvent } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Navigate, useLocation, useNavigate } from 'react-router-dom';
import VisibilityOutlinedIcon from '@mui/icons-material/VisibilityOutlined';
import VisibilityOffOutlinedIcon from '@mui/icons-material/VisibilityOffOutlined';
import { Alert, Box, Button, Container, IconButton, InputAdornment, Paper, Stack, TextField, Tooltip, Typography } from '@mui/material';
import { api } from '@shared/api/client';
import { useUser } from '@entities/session/auth';

export function LoginPage() {
  const { user, refresh } = useUser();
  const navigate = useNavigate();
  const location = useLocation();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const routeState = location.state as { from?: string; sessionExpired?: boolean } | null;
  const mutation = useMutation({ mutationFn: () => api.login(username, password), onSuccess: async () => { await refresh(); navigate(routeState?.from ?? '/', { replace: true }); } });
  const submit = (event: FormEvent) => { event.preventDefault(); mutation.mutate(); };
  if (user) return <Navigate to="/" replace />;
  return <Container maxWidth="xs" sx={{ minHeight: '100dvh', display: 'grid', alignItems: 'center', py: 3 }}><Paper component="form" onSubmit={submit} variant="outlined" sx={{ p: { xs: 3, sm: 4 }, boxShadow: '0 18px 60px rgba(0,0,0,.28)' }}><Stack spacing={2.5}><Box><Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 0.75 }}><Box aria-hidden sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: 'primary.main', boxShadow: '0 0 18px rgba(110, 168, 254, .75)' }} /><Typography component="h1" variant="h4">Serega</Typography></Stack><Typography color="text.secondary">Sign in to open the shared board.</Typography></Box>{routeState?.sessionExpired && <Alert severity="warning">Your session expired. Sign in again to continue.</Alert>}{mutation.error instanceof Error && <Alert severity="error">{mutation.error.message}</Alert>}<TextField label="Username" name="username" autoComplete="username" autoCapitalize="none" autoFocus required value={username} onChange={(event) => setUsername(event.target.value)} /><TextField label="Password" name="password" type={showPassword ? 'text' : 'password'} autoComplete="current-password" required value={password} onChange={(event) => setPassword(event.target.value)} slotProps={{ input: { endAdornment: <InputAdornment position="end"><Tooltip title={showPassword ? 'Hide password' : 'Show password'}><IconButton edge="end" aria-label={showPassword ? 'Hide password' : 'Show password'} onClick={() => setShowPassword((value) => !value)}>{showPassword ? <VisibilityOffOutlinedIcon /> : <VisibilityOutlinedIcon />}</IconButton></Tooltip></InputAdornment> } }} /><Button type="submit" variant="contained" disabled={mutation.isPending}>{mutation.isPending ? 'Signing in…' : 'Sign in'}</Button></Stack></Paper></Container>;
}
