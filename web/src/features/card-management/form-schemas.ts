import { z } from 'zod';

export const cardFormSchema = z.object({
  title: z.string().trim().min(1, 'Title is required.').refine((value) => Array.from(value).length <= 200, 'Title must be 200 characters or fewer.'),
  description: z.string().refine((value) => Array.from(value).length <= 5000, 'Description must be 5,000 characters or fewer.'),
  status: z.enum(['todo', 'doing', 'done']),
});

export type CardFormValues = z.infer<typeof cardFormSchema>;
export const cardFormDefaults: CardFormValues = { title: '', description: '', status: 'todo' };
