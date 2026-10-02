import { useEffect, useState } from 'react';
import { zodResolver } from '@hookform/resolvers/zod';
import { Controller, useForm } from 'react-hook-form';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Dialog, DialogActions, DialogContent, DialogTitle, MenuItem, Paper, Stack, TextField, Typography } from '@mui/material';
import { cardStatuses } from '@entities/card/status';
import { captureCardCache, reconcileSavedCard } from '@entities/card/cache';
import { cardFormDefaults, cardFormSchema, type CardFormValues } from '@features/card-management/form-schemas';
import { api, apiFormErrors, ApiRequestError, type Card } from '@shared/api/client';
import { formatDateTime } from '@shared/lib/time';

export function CardDialog({ card, onClose, onSaved }: { card?: Card; onClose: () => void; onSaved?: (card: Card) => void }) {
  const client = useQueryClient();
  const [base, setBase] = useState(card);
  const [reviewed, setReviewed] = useState(false);
  const [reviewing, setReviewing] = useState(false);
  const [reviewError, setReviewError] = useState<unknown>();
  const form = useForm<CardFormValues>({ resolver: zodResolver(cardFormSchema), defaultValues: card ? { title: card.title, description: card.description, status: card.status } : cardFormDefaults, mode: 'onChange' });
  const save = useMutation({ mutationFn: (values: CardFormValues) => base
    ? api.updateCard(base.id, { ...values, version: base.version })
    : api.createCard({ title: values.title, description: values.description }),
  onMutate: () => captureCardCache(client),
  onSuccess: async (saved, _variables, scope) => { const latest = await reconcileSavedCard(client, saved, scope); if (!latest) return; onSaved?.(latest); onClose(); },
  onError: (error) => { if (error instanceof ApiRequestError && (error.code === 'card_conflict' || error.code === 'not_found')) void client.invalidateQueries({ queryKey: ['cards'] }); },
  });
  const conflict = save.error instanceof ApiRequestError && save.error.code === 'card_conflict';
  const removed = save.error instanceof ApiRequestError && save.error.code === 'not_found';
  useEffect(() => {
    if (!save.error) return;
    const details = apiFormErrors(save.error);
    for (const [field, message] of Object.entries(details.fields)) if (field in cardFormDefaults) form.setError(field as keyof CardFormValues, { type: 'server', message });
  }, [save.error, form]);
  const reviewLatest = async () => {
    if (!base) return;
    setReviewing(true);
    setReviewError(undefined);
    try {
      const result = await client.fetchQuery({ queryKey: ['cards'], queryFn: ({ signal }) => api.listCards(signal), staleTime: 0 });
      const latest = result.cards.find((item) => item.id === base.id);
      if (!latest) throw new Error('This card was deleted. Your draft is still here so you can copy it before closing.');
      setBase(latest);
      setReviewed(true);
      save.reset();
    } catch (error) {
      setReviewError(error);
    } finally {
      setReviewing(false);
    }
  };
  const errors = form.formState.errors;
  const close = () => { if (!save.isPending && !reviewing) onClose(); };
  return <Dialog open onClose={close} fullWidth maxWidth="sm" aria-labelledby="card-dialog-title">
    <DialogTitle id="card-dialog-title">{card ? 'Edit card' : 'New card'}</DialogTitle>
    <DialogContent><Stack component="form" id="card-form" onSubmit={form.handleSubmit((values) => { if (!save.isPending && !conflict && !removed && !reviewing) save.mutate(values); })} noValidate spacing={1.5} sx={{ pt: 1 }}>
      {conflict && <Alert severity="warning" action={<Button color="inherit" disabled={reviewing} onClick={() => { void reviewLatest(); }}>{reviewing ? 'Loading…' : 'Review latest'}</Button>}>Someone changed this card. Your draft is kept below. Review the latest version before saving again.</Alert>}
      {removed && <Alert severity="warning">This card was deleted. Your draft is kept below so you can copy it before closing.</Alert>}
      {save.error && !conflict && !removed && <Alert severity="error">{apiFormErrors(save.error).general}</Alert>}
      {reviewError instanceof Error && <Alert severity="error">{reviewError.message}</Alert>}
      {reviewed && base && <Paper variant="outlined" sx={{ p: 1.5 }}><Stack spacing={1}><Typography variant="subtitle2">Latest saved version</Typography><Typography sx={{ overflowWrap: 'anywhere' }}>{base.title}</Typography>{base.description && <Typography variant="body2" color="text.secondary" sx={{ whiteSpace: 'pre-wrap', overflowWrap: 'anywhere', maxHeight: 180, overflow: 'auto' }}>{base.description}</Typography>}<Typography variant="caption" color="text.secondary">{cardStatuses.find((status) => status.value === base.status)?.label} · {base.updatedBy.username} · {formatDateTime(base.updatedAt)}</Typography><Typography variant="body2">Your draft is below. Saving applies it over this version.</Typography></Stack></Paper>}
      <TextField autoFocus required label="Title" {...form.register('title')} error={Boolean(errors.title)} helperText={errors.title?.message} disabled={save.isPending} />
      <TextField label="Description (optional)" multiline minRows={3} maxRows={12} {...form.register('description')} error={Boolean(errors.description)} helperText={errors.description?.message} disabled={save.isPending} />
      {card && <Controller control={form.control} name="status" render={({ field }) => <TextField select label="Status" {...field} disabled={save.isPending} error={Boolean(errors.status)} helperText={errors.status?.message}>{cardStatuses.map((status) => <MenuItem key={status.value} value={status.value}>{status.label}</MenuItem>)}</TextField>} />}
      {card && <Typography variant="caption" color="text.secondary">Created by {card.createdBy.username} · {formatDateTime(card.createdAt)}</Typography>}
    </Stack></DialogContent>
    <DialogActions sx={{ px: 3, pb: 2.5 }}><Button onClick={close} disabled={save.isPending || reviewing}>Cancel</Button><Button form="card-form" type="submit" variant="contained" disabled={save.isPending || conflict || removed || reviewing}>{save.isPending ? 'Saving…' : card ? 'Save changes' : 'Create card'}</Button></DialogActions>
  </Dialog>;
}
