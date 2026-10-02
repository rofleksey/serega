import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import AddIcon from '@mui/icons-material/Add';
import RefreshIcon from '@mui/icons-material/Refresh';
import { Alert, Box, Button, IconButton, Paper, Stack, Tab, Tabs, Tooltip, Typography, useMediaQuery, useTheme } from '@mui/material';
import { CardItem } from '@entities/card/CardItem';
import { captureCardCache, reconcileSavedCard } from '@entities/card/cache';
import { cardStatuses } from '@entities/card/status';
import { CardDialog } from '@features/card-management/CardDialog';
import { api, ApiRequestError, type Card, type CardStatus } from '@shared/api/client';
import { LoadingState, RecoverableError } from '@shared/ui/async-states';
import { ConfirmationDialog } from '@shared/ui/confirmation-dialog';

export function KanbanBoard() {
  const client = useQueryClient();
  const mobile = useMediaQuery(useTheme().breakpoints.down('md'));
  const query = useQuery({ queryKey: ['cards'], queryFn: ({ signal }) => api.listCards(signal), refetchInterval: 5000, refetchOnWindowFocus: true });
  const [selectedStatus, setSelectedStatus] = useState<CardStatus>('todo');
  const [editing, setEditing] = useState<Card | 'new'>();
  const [deleting, setDeleting] = useState<Card>();
  const [notice, setNotice] = useState<string>();
  const selectColumn = (status: CardStatus) => {
    setSelectedStatus(status);
    window.scrollTo({ top: 0, behavior: 'auto' });
  };
  const move = useMutation({ mutationFn: ({ card, status }: { card: Card; status: CardStatus }) => api.updateCard(card.id, { title: card.title, description: card.description, status, version: card.version }),
    onMutate: () => captureCardCache(client),
    onSuccess: async (card, _variables, scope) => { const latest = await reconcileSavedCard(client, card, scope); if (!latest) return; setNotice(undefined); if (mobile) selectColumn(latest.status); },
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
  const columns = cardStatuses.map((status) => ({ ...status, cards: cards.filter((card) => card.status === status.value) }));
  const syncState = query.isPaused ? 'Offline' : query.isFetching ? 'Syncing…' : query.isError ? 'Updates paused' : 'Board up to date';
  return <Stack spacing={1.5}>
    <Stack component="header" direction="row" spacing={1} alignItems="center" justifyContent="space-between">
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={{ xs: .25, sm: 1.5 }} alignItems={{ sm: 'center' }} sx={{ minWidth: 0 }}><Typography component="h1" variant="h4">Team board</Typography><Typography variant="caption" color={query.isPaused || query.isError ? 'warning.main' : 'text.secondary'}>{query.data && `${cards.length} ${cards.length === 1 ? 'card' : 'cards'} · `}{syncState}</Typography></Stack>
      <Stack direction="row" spacing={.5} alignItems="center" sx={{ flexShrink: 0 }}><Tooltip title="Refresh board"><span><IconButton aria-label="Refresh board" disabled={query.isFetching} onClick={() => { void query.refetch(); }}><RefreshIcon fontSize="small" /></IconButton></span></Tooltip><Button variant="contained" startIcon={<AddIcon />} onClick={() => setEditing('new')}>New card</Button></Stack>
    </Stack>
    {query.error && <RecoverableError title={query.data ? 'Could not refresh the board' : 'Could not load the board'} error={query.error} onRetry={() => { void query.refetch(); }} />}
    {notice && <Alert severity="warning" onClose={() => setNotice(undefined)}>{notice}</Alert>}
    {query.isPending ? <LoadingState label="Loading board" rows={3} /> : query.data && <>
      {mobile && <Tabs value={selectedStatus} onChange={(_event, next: CardStatus) => selectColumn(next)} variant="fullWidth" aria-label="Board columns" sx={{ position: 'sticky', top: 'calc(48px + env(safe-area-inset-top))', bgcolor: 'background.default', zIndex: (theme) => theme.zIndex.appBar - 1 }}>{columns.map((status) => <Tab key={status.value} id={`board-tab-${status.value}`} aria-controls={`board-column-${status.value}`} value={status.value} label={`${status.label} (${status.cards.length})`} sx={{ px: .5 }} />)}</Tabs>}
      <Box sx={{ display: 'grid', gridTemplateColumns: { xs: 'minmax(0, 1fr)', md: 'repeat(3, minmax(0, 1fr))' }, gap: 1.5, alignItems: 'start' }}>
        {columns.filter((status) => !mobile || status.value === selectedStatus).map((status) => {
          return <Paper key={status.value} component="section" role={mobile ? 'tabpanel' : undefined} id={`board-column-${status.value}`} aria-labelledby={mobile ? `board-tab-${status.value}` : `column-${status.value}`} variant="outlined" sx={{ p: { xs: 0, md: 1 }, borderWidth: { xs: 0, md: 1 }, bgcolor: 'transparent', minWidth: 0 }}>
            <Stack spacing={1}>{!mobile && <Stack direction="row" alignItems="center" spacing={.75} sx={{ px: .5, py: .25 }}><Box aria-hidden sx={{ width: 6, height: 6, borderRadius: '50%', bgcolor: status.color }} /><Typography id={`column-${status.value}`} component="h2" variant="subtitle2" sx={{ flex: 1 }}>{status.label}</Typography><Typography variant="caption" color="text.secondary" sx={{ fontVariantNumeric: 'tabular-nums' }}>{status.cards.length}</Typography></Stack>}
              {status.cards.length ? status.cards.map((card) => <CardItem key={card.id} card={card} pending={move.isPending || remove.isPending} onEdit={() => setEditing(card)} onMove={(next) => move.mutate({ card, status: next })} onDelete={() => setDeleting(card)} />) : <Box sx={{ py: 2, px: 1.5, textAlign: 'center' }}><Typography variant="body2" color="text.secondary">No cards in {status.label}.</Typography></Box>}
            </Stack>
          </Paper>;
        })}
      </Box>
    </>}
    {editing && <CardDialog key={editing === 'new' ? 'new' : editing.id} card={editing === 'new' ? undefined : editing} onSaved={(card) => { if (mobile) selectColumn(card.status); }} onClose={() => setEditing(undefined)} />}
    <ConfirmationDialog open={Boolean(deleting)} title="Delete this card?" target={deleting?.title ?? ''} effect="The card will be permanently removed from the shared board for everyone." confirmLabel="Delete card" confirmColor="error" pending={remove.isPending} onCancel={() => setDeleting(undefined)} onConfirm={() => { if (deleting) remove.mutate(deleting); }} />
  </Stack>;
}
