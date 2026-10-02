import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import EditOutlinedIcon from '@mui/icons-material/EditOutlined';
import { Box, Button, Paper, Stack, TextField, Tooltip, Typography } from '@mui/material';
import { cardStatuses } from '@entities/card/status';
import type { Card, CardStatus } from '@shared/api/client';
import { formatDateTime, formatRelativeTime } from '@shared/lib/time';
import { OverflowActions } from '@shared/ui/overflow-actions';

export function CardItem({ card, pending, onEdit, onMove, onDelete }: { card: Card; pending: boolean; onEdit: () => void; onMove: (status: CardStatus) => void; onDelete: () => void }) {
  return <Paper component="article" aria-labelledby={`card-title-${card.id}`} variant="outlined" sx={{ p: 1.25, bgcolor: 'background.paper', minWidth: 0 }}>
    <Stack spacing={.75}>
      <Stack direction="row" spacing={.5} alignItems="flex-start">
        <Button onClick={onEdit} disabled={pending} aria-label={`Edit ${card.title}`} sx={{ flex: 1, minWidth: 0, p: 0, justifyContent: 'flex-start', textAlign: 'left', color: 'text.primary', '&:hover': { bgcolor: 'transparent', color: 'primary.main' } }}><Typography id={`card-title-${card.id}`} component="h3" variant="body2" fontWeight={650} sx={{ overflowWrap: 'anywhere', lineHeight: 1.4 }}>{card.title}</Typography></Button>
        <Box sx={{ flexShrink: 0, mr: '-6px !important' }}><OverflowActions label={`Actions for ${card.title}`} actions={[
          { label: 'Edit card', icon: <EditOutlinedIcon fontSize="small" />, disabled: pending, onSelect: onEdit },
          { label: 'Delete card', icon: <DeleteOutlineIcon fontSize="small" />, danger: true, disabled: pending, onSelect: onDelete },
        ]} /></Box>
      </Stack>
      {card.description && <Typography variant="body2" color="text.secondary" sx={{ whiteSpace: 'pre-wrap', overflowWrap: 'anywhere', display: '-webkit-box', WebkitLineClamp: 3, WebkitBoxOrient: 'vertical', overflow: 'hidden' }}>{card.description}</Typography>}
      <Stack direction="row" spacing={1} alignItems="center" justifyContent="space-between">
        <Tooltip title={`Created by ${card.createdBy.username} · ${formatDateTime(card.createdAt)}. Updated by ${card.updatedBy.username} · ${formatDateTime(card.updatedAt)}.`}>
          <Typography variant="caption" color="text.secondary" sx={{ minWidth: 0, overflowWrap: 'anywhere' }}>{card.updatedBy.username} · <Box component="time" dateTime={card.updatedAt}>{formatRelativeTime(card.updatedAt)}</Box></Typography>
        </Tooltip>
        <TextField select size="small" variant="standard" id={`status-${card.id}`} value={card.status} disabled={pending} onChange={(event) => onMove(event.target.value as CardStatus)} slotProps={{ input: { disableUnderline: true }, select: { native: true }, htmlInput: { 'aria-label': `Status for ${card.title}` } }} sx={{ width: 138, flexShrink: 0, '& .MuiInputBase-root': { minHeight: 32, px: 1, borderRadius: 1, bgcolor: 'action.hover', '@media (pointer: coarse), (max-width: 899.95px)': { minHeight: 44 } }, '& .MuiNativeSelect-select': { py: .5 } }}>
          {cardStatuses.map((status) => <option key={status.value} value={status.value}>{status.label}</option>)}
        </TextField>
      </Stack>
    </Stack>
  </Paper>;
}
