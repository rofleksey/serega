import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import AddIcon from '@mui/icons-material/Add';
import RefreshIcon from '@mui/icons-material/Refresh';
import { Alert, Box, Button, Chip, IconButton, Paper, Stack, Tooltip, Typography } from '@mui/material';
import { CardItem } from '@entities/card/CardItem';
import { cardStatuses } from '@entities/card/status';
import { CardDialog } from '@features/card-management/CardDialog';
import { api, ApiRequestError, type Card, type CardStatus } from '@shared/api/client';
import { EmptyState, LoadingState, RecoverableError } from '@shared/ui/async-states';
import { ConfirmationDialog } from '@shared/ui/confirmation-dialog';
import { PageHeader } from '@shared/ui/page-header';

export function KanbanBoard() {
  const client = useQueryClient();
  const query = useQuery({ queryKey: ['cards'], queryFn: ({ signal }) => api.listCards(signal), refetchInterval: 5000, refetchOnWindowFocus: true });
  const [editing, setEditing] = useState<Card | 'new'>();
  const [deleting, setDeleting] = useState<Card>();
  const [notice, setNotice] = useState<string>();
  const move = useMutation({ mutationFn: ({ card, status }: { card: Card; status: CardStatus }) => api.updateCard(card.id, { title: card.title, description: card.description, status, version: card.version }),
    onSuccess: async () => { setNotice(undefined); await client.invalidateQueries({ queryKey: ['cards'] }); },
    onError: (error) => {
      setNotice(error instanceof ApiRequestError && error.code === 'card_conflict' ? 'This card changed while you were moving it. The board has been refreshed; review the card and try again.' : error.message);
      void client.invalidateQueries({ queryKey: ['cards'] });
    },
  });
  const remove = useMutation({ mutationFn: (card: Card) => api.deleteCard(card.id, card.version),
    onSuccess: async () => { setDeleting(undefined); setNotice(undefined); await client.invalidateQueries({ queryKey: ['cards'] }); },
    onError: (error) => {
      setDeleting(undefined);
      setNotice(error instanceof ApiRequestError && error.code === 'card_conflict' ? 'This card changed while you were deleting it. Review the latest card before deleting it again.' : error.message);
      void client.invalidateQueries({ queryKey: ['cards'] });
    },
  });
  const cards = query.data?.cards ?? [];
  return <Stack spacing={3}>
    <PageHeader title="Team board" description="A shared place to plan, make progress, and get things done." actions={<Button variant="contained" startIcon={<AddIcon />} onClick={() => setEditing('new')}>New card</Button>} />
    <Stack direction="row" alignItems="center" justifyContent="space-between" spacing={1}>
      <Typography variant="body2" color="text.secondary">{cards.length} {cards.length === 1 ? 'card' : 'cards'} · Shared with everyone</Typography>
      <Stack direction="row" alignItems="center" spacing={.5}><Typography variant="caption" color="text.secondary">{query.isPaused ? 'Offline' : query.isFetching ? 'Syncing…' : query.isError ? 'Updates paused' : 'Board up to date'}</Typography><Tooltip title="Refresh board"><span><IconButton aria-label="Refresh board" size="small" disabled={query.isFetching} onClick={() => { void query.refetch(); }}><RefreshIcon fontSize="small" /></IconButton></span></Tooltip></Stack>
    </Stack>
    {query.error && <RecoverableError title={query.data ? 'Could not refresh the board' : 'Could not load the board'} error={query.error} onRetry={() => { void query.refetch(); }} />}
    {notice && <Alert severity="warning" onClose={() => setNotice(undefined)}>{notice}</Alert>}
    {query.isPending ? <LoadingState label="Loading board" rows={3} /> : query.data && <>
      {cards.length === 0 && <Paper variant="outlined"><EmptyState title="A fresh start" description="Add the first card. Everyone on the team can pick it up and move it forward." action={<Button startIcon={<AddIcon />} onClick={() => setEditing('new')}>Create your first card</Button>} /></Paper>}
      <Box sx={{ display: 'grid', gridTemplateColumns: { xs: 'minmax(0, 1fr)', md: 'repeat(3, minmax(0, 1fr))' }, gap: 2.5, alignItems: 'start' }}>
        {cardStatuses.map((status) => {
          const column = cards.filter((card) => card.status === status.value);
          return <Paper key={status.value} component="section" aria-labelledby={`column-${status.value}`} variant="outlined" sx={{ p: { xs: 1.5, sm: 2 }, bgcolor: 'rgba(17, 21, 27, .55)', minHeight: { md: 360 } }}>
            <Stack spacing={2}><Stack direction="row" alignItems="center" spacing={1}><Box aria-hidden sx={{ width: 8, height: 8, borderRadius: '50%', bgcolor: status.color }} /><Typography id={`column-${status.value}`} component="h2" variant="subtitle1" fontWeight={700} sx={{ flex: 1 }}>{status.label}</Typography><Chip label={column.length} size="small" sx={{ height: 24, minWidth: 28 }} /></Stack>
              {column.length ? column.map((card) => <CardItem key={card.id} card={card} pending={move.isPending || remove.isPending} onEdit={() => setEditing(card)} onMove={(next) => move.mutate({ card, status: next })} onDelete={() => setDeleting(card)} />) : <Box sx={{ border: '1px dashed', borderColor: 'divider', borderRadius: 1.5, py: 4, px: 2, textAlign: 'center' }}><Typography variant="body2" color="text.secondary">{status.description}</Typography></Box>}
            </Stack>
          </Paper>;
        })}
      </Box>
    </>}
    {editing && <CardDialog key={editing === 'new' ? 'new' : editing.id} card={editing === 'new' ? undefined : editing} onClose={() => setEditing(undefined)} />}
    <ConfirmationDialog open={Boolean(deleting)} title="Delete this card?" target={deleting?.title ?? ''} effect="The card will be permanently removed from the shared board for everyone." confirmLabel="Delete card" confirmColor="error" pending={remove.isPending} onCancel={() => setDeleting(undefined)} onConfirm={() => { if (deleting) remove.mutate(deleting); }} />
  </Stack>;
}
