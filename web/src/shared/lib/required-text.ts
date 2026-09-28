import { z } from 'zod';

export const requiredText = (label: string) => z.string().refine((value) => Boolean(value.trim()), `${label} is required.`);
