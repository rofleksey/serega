import { describe, expect, it } from 'vitest';
import { formatDateTime, formatRelativeTime } from '@shared/lib/time';

describe('time formatting', () => {
  it('renders valid timestamps for people and rejects invalid values', () => {
    expect(formatDateTime('2026-06-13T14:14:17Z')).toBeTruthy();
    expect(formatDateTime('2026-06-13T14:14:17Z')).not.toContain('T14:14:17Z');
    expect(formatDateTime('not-a-date')).toBeUndefined();
  });

  it('renders relative time from an explicit clock', () => {
    const now = new Date('2026-09-25T12:00:00Z').getTime();
    expect(formatRelativeTime('2026-09-25T11:55:00Z', now)).toMatch(/5 minutes ago/);
    expect(formatRelativeTime('invalid', now)).toBeUndefined();
  });
});
