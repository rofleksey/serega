import type { CardStatus } from '@shared/api/client';

export const cardStatuses: Array<{ value: CardStatus; label: string; color: string; description: string }> = [
  { value: 'todo', label: 'To do', color: 'primary.main', description: 'Ready when you are' },
  { value: 'doing', label: 'In progress', color: 'warning.main', description: 'Getting things done' },
  { value: 'done', label: 'Done', color: 'success.main', description: 'A little more accomplished' },
];
