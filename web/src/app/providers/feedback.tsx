import { createContext, useCallback, useContext, useMemo, useState } from 'react';
import { Alert, Snackbar } from '@mui/material';

type Feedback = { message: string; severity: 'success' | 'error' };
type FeedbackContextValue = { showSuccess: (message: string) => void; showError: (message: string) => void };

const FeedbackContext = createContext<FeedbackContextValue | undefined>(undefined);

export function FeedbackProvider({ children }: { children: React.ReactNode }) {
  const [feedback, setFeedback] = useState<Feedback>();
  const showSuccess = useCallback((message: string) => setFeedback({ message, severity: 'success' }), []);
  const showError = useCallback((message: string) => setFeedback({ message, severity: 'error' }), []);
  const value = useMemo(() => ({ showSuccess, showError }), [showSuccess, showError]);
  return <FeedbackContext.Provider value={value}>
    {children}
    <Snackbar open={Boolean(feedback)} autoHideDuration={5000} onClose={() => setFeedback(undefined)} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
      {feedback ? <Alert severity={feedback.severity} variant="filled" onClose={() => setFeedback(undefined)} sx={{ width: '100%' }}>{feedback.message}</Alert> : undefined}
    </Snackbar>
  </FeedbackContext.Provider>;
}

export function useFeedback() {
  const value = useContext(FeedbackContext);
  if (!value) throw new Error('useFeedback must be used inside FeedbackProvider');
  return value;
}
