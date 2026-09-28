import { QueryClient } from '@tanstack/react-query';

export function createQueryClient() {
  return new QueryClient({ defaultOptions: {
    queries: { retry: 1, refetchOnWindowFocus: true },
    // This app has no offline queue. Failed writes must return control to the
    // user immediately, keeping their draft available for an explicit retry.
    mutations: { networkMode: 'always' },
  } });
}
