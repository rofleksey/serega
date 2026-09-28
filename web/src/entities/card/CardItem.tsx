import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import EditOutlinedIcon from '@mui/icons-material/EditOutlined';
import { Box, Button, Paper, Stack, TextField, Tooltip, Typography } from '@mui/material';
import { cardStatuses } from '@entities/card/status';
import type { Card, CardStatus } from '@shared/api/client';
import { formatDateTime, formatRelativeTime } from '@shared/lib/time';
import { OverflowActions } from '@shared/ui/overflow-actions';

export function CardItem({ card, pending, onEdit, onMove, onDelete }: { card: Card; pending: boolean; onEdit: () => void; onMove: (status: CardStatus) => void; onDelete: () => void }) {
  return <Paper component="article" aria-labelledby={`card-title-${card.id}`} variant="outlined" sx={{ p: 2, bgcolor: 'background.paper', minWidth: 0, boxShadow: '0 3px 12px rgba(0,0,0,.08)' }}>
    <Stack spacing={1.5}>
      <Stack direction="row" spacing={.5} alignItems="flex-start">
        <Button onClick={onEdit} disabled={pending} aria-label={`Edit ${card.title}`} sx={{ flex: 1, minWidth: 0, p: 0, justifyContent: 'flex-start', textAlign: 'left', color: 'text.primary', '&:hover': { bgcolor: 'transparent', color: 'primary.main' } }}><Typography id={`card-title-${card.id}`} component="h3" variant="subtitle1" fontWeight={700} sx={{ overflowWrap: 'anywhere', lineHeight: 1.4 }}>{card.title}</Typography></Button>
        <Box sx={{ mr: '-10px !important' }}><OverflowActions label={`Actions for ${card.title}`} actions={[
          { label: 'Edit card', icon: <EditOutlinedIcon fontSize="small" />, disabled: pending, onSelect: onEdit },
          { label: 'Delete card', icon: <DeleteOutlineIcon fontSize="small" />, danger: true, disabled: pending, onSelect: onDelete },
        ]} /></Box>
      </Stack>
      {card.description && <Typography variant="body2" color="text.secondary" sx={{ whiteSpace: 'pre-wrap', overflowWrap: 'anywhere', display: '-webkit-box', WebkitLineClamp: 3, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>{card.description}</Typography>}
      <Tooltip title={`Created by ${card.createdBy.username} · ${formatDateTime(card.createdAt)}. Updated by ${card.updatedBy.username} · ${formatDateTime(card.updatedAt)}.`}>
        <Typography variant="caption" color="text.secondary" sx={{ overflowWrap: 'anywhere' }}>Updated by {card.updatedBy.username} · <Box component="time" dateTime={card.updatedAt}>{formatRelativeTime(card.updatedAt)}</Box></Typography>
      </Tooltip>
      <TextField select size="small" id={`status-${card.id}`} label="Status" value={card.status} disabled={pending} onChange={(event) => onMove(event.target.value as CardStatus)} slotProps={{ select: { native: true }, htmlInput: { 'aria-label': `Status for ${card.title}` } }}>
        {cardStatuses.map((status) => <option key={status.value} value={status.value}>{status.label}</option>)}
      </TextField>
    </Stack>
  </Paper>;
}
