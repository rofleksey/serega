import FiberManualRecordIcon from '@mui/icons-material/FiberManualRecord';
import { Stack, Typography } from '@mui/material';

type Tone = 'neutral' | 'primary' | 'success' | 'warning' | 'error';

const colors: Record<Tone, string> = {
  neutral: 'text.secondary',
  primary: 'primary.main',
  success: 'success.main',
  warning: 'warning.main',
  error: 'error.main',
};

export function StatusIndicator({ label, tone = 'neutral', detail }: { label: string; tone?: Tone; detail?: string }) {
  return <Stack component="span" direction="row" spacing={0.75} alignItems="center" sx={{ minWidth: 0 }}>
    <FiberManualRecordIcon aria-hidden fontSize="inherit" sx={{ color: colors[tone], flexShrink: 0 }} />
    <Typography component="span" variant="body2" fontWeight={650} noWrap>{label}</Typography>
    {detail && <Typography component="span" variant="body2" color="text.secondary" noWrap>{detail}</Typography>}
  </Stack>;
}
