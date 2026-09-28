import type { ReactNode } from 'react';
import { Box, CircularProgress } from '@mui/material';
import { Navigate, useLocation } from 'react-router-dom';
import { useUser } from '@entities/session/auth';

export function RequireUser({ children }: { children: ReactNode }) {
  const { user, isLoading } = useUser();
  const location = useLocation();
  if (isLoading) return <Box sx={{ minHeight: '100vh', display: 'grid', placeItems: 'center' }}><CircularProgress aria-label="Loading session" /></Box>;
  return user ? <>{children}</> : <Navigate to="/login" replace state={{ from: location.pathname }} />;
}
