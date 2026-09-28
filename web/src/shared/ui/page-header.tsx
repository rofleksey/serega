import { Link as RouterLink } from 'react-router-dom';
import { Breadcrumbs, Link, Stack, Typography } from '@mui/material';

export type PageCrumb = { label: string; to?: string };

export function PageHeader({ title, description, crumbs = [], actions }: { title: string; description?: string; crumbs?: PageCrumb[]; actions?: React.ReactNode }) {
  return <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'flex-start' }} justifyContent="space-between">
    <Stack spacing={0.75} sx={{ minWidth: 0 }}>
      {crumbs.length > 0 && <Breadcrumbs aria-label="Breadcrumb" sx={{ color: 'text.secondary', '& .MuiBreadcrumbs-ol': { flexWrap: 'nowrap' } }}>
        {crumbs.map((crumb, index) => crumb.to
          ? <Link key={`${crumb.label}-${index}`} component={RouterLink} to={crumb.to} color="inherit" noWrap>{crumb.label}</Link>
          : <Typography key={`${crumb.label}-${index}`} color="text.primary" variant="body2" noWrap>{crumb.label}</Typography>)}
      </Breadcrumbs>}
      <Typography component="h1" variant="h4" sx={{ overflowWrap: 'anywhere' }}>{title}</Typography>
      {description && <Typography color="text.secondary" sx={{ maxWidth: 760 }}>{description}</Typography>}
    </Stack>
    {actions && <Stack direction="row" spacing={1} useFlexGap flexWrap="wrap" sx={{ flexShrink: 0 }}>{actions}</Stack>}
  </Stack>;
}
