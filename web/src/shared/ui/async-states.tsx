import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import InboxOutlinedIcon from '@mui/icons-material/InboxOutlined';
import { Alert, Box, Button, Skeleton, Stack, Typography } from '@mui/material';

export function LoadingState({ label = 'Loading', rows = 3 }: { label?: string; rows?: number }) {
  return <Stack aria-label={label} role="status" spacing={1.25} sx={{ py: 1 }}>
    <Skeleton variant="rounded" height={40} />
    {Array.from({ length: Math.max(0, rows - 1) }, (_, index) => <Skeleton key={index} variant="rounded" height={48} />)}
  </Stack>;
}

export function EmptyState({ title, description, action }: { title: string; description: string; action?: React.ReactNode }) {
  return <Box sx={{ py: 3, px: 2, textAlign: 'center' }}>
    <InboxOutlinedIcon aria-hidden sx={{ color: 'text.secondary', fontSize: 26, mb: 1 }} />
    <Typography component="h2" variant="h6">{title}</Typography>
    <Typography color="text.secondary" sx={{ mt: 0.5, mx: 'auto', maxWidth: 520 }}>{description}</Typography>
    {action && <Box sx={{ mt: 2 }}>{action}</Box>}
  </Box>;
}

export function RecoverableError({ title = 'Something went wrong', error, onRetry, retryLabel = 'Try again' }: { title?: string; error: unknown; onRetry?: () => void; retryLabel?: string }) {
  const message = error instanceof Error ? error.message : 'The request could not be completed.';
  return <Alert severity="error" icon={<ErrorOutlineIcon />} action={onRetry ? <Button color="inherit" onClick={onRetry}>{retryLabel}</Button> : undefined}>
    <Typography fontWeight={700}>{title}</Typography>
    <Typography variant="body2">{message}</Typography>
  </Alert>;
}
