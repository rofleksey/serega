import { Alert, Button, Dialog, DialogActions, DialogContent, DialogTitle, Stack, Typography } from '@mui/material';

export function ConfirmationDialog({ open, title, target, effect, confirmLabel, confirmColor = 'primary', pending = false, safetyNote, onCancel, onConfirm }: { open: boolean; title: string; target: string; effect: string; confirmLabel: string; confirmColor?: 'primary' | 'error' | 'warning'; pending?: boolean; safetyNote?: string; onCancel: () => void; onConfirm: () => void }) {
  return <Dialog open={open} onClose={() => !pending && onCancel()} aria-labelledby="confirmation-dialog-title" fullWidth maxWidth="xs">
    <DialogTitle id="confirmation-dialog-title">{title}</DialogTitle>
    <DialogContent>
      <Stack spacing={1.5}>
        <Typography><strong>Target:</strong> {target}</Typography>
        <Typography>{effect}</Typography>
        {safetyNote && <Alert severity="info">{safetyNote}</Alert>}
      </Stack>
    </DialogContent>
    <DialogActions sx={{ px: 3, pb: 2.5 }}>
      <Button onClick={onCancel} disabled={pending}>Cancel</Button>
      <Button variant="contained" color={confirmColor} onClick={onConfirm} disabled={pending}>{pending ? 'Working…' : confirmLabel}</Button>
    </DialogActions>
  </Dialog>;
}
