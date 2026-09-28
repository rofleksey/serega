import { lazy, Suspense } from 'react';
import { QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { CssBaseline, ThemeProvider } from '@mui/material';
import { UserProvider } from '@entities/session/auth';
import { FeedbackProvider } from '@app/providers/feedback';
import { createQueryClient } from '@app/providers/query-client';
import { AppLayout } from '@app/layout/AppLayout';
import { RequireUser } from '@app/routing/RequireUser';
import { theme } from '@app/theme';
import { LoginPage } from '@pages/login/LoginPage';
import { LoadingState } from '@shared/ui/async-states';

const BoardPage = lazy(() => import('@pages/board/BoardPage').then(({ BoardPage }) => ({ default: BoardPage })));
const queryClient = createQueryClient();

export function App() {
  return <QueryClientProvider client={queryClient}>
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <FeedbackProvider>
        <BrowserRouter>
          <UserProvider>
            <Suspense fallback={<LoadingState label="Loading page" />}>
              <Routes>
                <Route path="/login" element={<LoginPage />} />
                <Route element={<RequireUser><AppLayout /></RequireUser>}>
                  <Route index element={<BoardPage />} />
                </Route>
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </Suspense>
          </UserProvider>
        </BrowserRouter>
      </FeedbackProvider>
    </ThemeProvider>
  </QueryClientProvider>;
}
