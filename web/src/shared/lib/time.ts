const dateTimeFormatter = new Intl.DateTimeFormat('en-GB', { dateStyle: 'medium', timeStyle: 'short' });
const relativeFormatter = new Intl.RelativeTimeFormat('en-GB', { numeric: 'always' });

export function formatDateTime(value?: string): string | undefined {
  if (!value) return undefined;
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return undefined;
  return dateTimeFormatter.format(parsed);
}

export function formatRelativeTime(value: string, now = Date.now()): string | undefined {
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return undefined;
  const seconds = Math.round((parsed.getTime() - now) / 1000);
  const absolute = Math.abs(seconds);
  if (absolute < 60) return relativeFormatter.format(seconds, 'second');
  if (absolute < 3600) return relativeFormatter.format(Math.round(seconds / 60), 'minute');
  if (absolute < 86400) return relativeFormatter.format(Math.round(seconds / 3600), 'hour');
  if (absolute < 2_592_000) return relativeFormatter.format(Math.round(seconds / 86400), 'day');
  if (absolute < 31_536_000) return relativeFormatter.format(Math.round(seconds / 2_592_000), 'month');
  return relativeFormatter.format(Math.round(seconds / 31_536_000), 'year');
}
