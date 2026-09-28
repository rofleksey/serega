import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { api, ApiRequestError, type User } from '@shared/api/client';
import { SESSION_EXPIRED_EVENT } from '@shared/api/session-expiry';
import { csrf } from '@shared/api/transport';

type UserContextValue = {
  user?: User;
  isLoading: boolean;
  refresh: () => Promise<unknown>;
  signOut: () => Promise<void>;
};

const UserContext = createContext<UserContextValue | undefined>(undefined);

export function UserProvider({ children }: { children: React.ReactNode }) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const sessionExpired = useRef(false);
  const [signedOut, setSignedOut] = useState(false);
  const query = useQuery({ queryKey: ['current-user'], queryFn: api.getCurrentUser, retry: false, enabled: !signedOut });
  const clearSession = useCallback(() => {
    setSignedOut(true);
    csrf.clear();
    void queryClient.cancelQueries();
    queryClient.clear();
  }, [queryClient]);
  const handleSessionExpired = useCallback(() => {
    if (sessionExpired.current) return;
    sessionExpired.current = true;
    const from = `${window.location.pathname}${window.location.search}${window.location.hash}`;
    clearSession();
    navigate('/login', { replace: true, state: { from: from === '/login' ? '/' : from, sessionExpired: true } });
  }, [clearSession, navigate]);
  useEffect(() => {
    window.addEventListener(SESSION_EXPIRED_EVENT, handleSessionExpired);
    return () => window.removeEventListener(SESSION_EXPIRED_EVENT, handleSessionExpired);
  }, [handleSessionExpired]);
  useEffect(() => {
    // Query retains successful data after a refetch error. A known session must
    // be discarded when a later /me request reports the cookie is missing.
    if (query.data && query.error instanceof ApiRequestError && query.error.code === 'unauthorized') handleSessionExpired();
  }, [query.data, query.error, handleSessionExpired]);
  const value = useMemo<UserContextValue>(() => ({
    user: signedOut ? undefined : query.data,
    isLoading: query.isLoading,
    refresh: async () => {
      const result = await query.refetch();
      if (result.error) throw result.error;
      if (result.data) {
        sessionExpired.current = false;
        setSignedOut(false);
      }
      return result;
    },
    signOut: async () => {
      await api.logout();
      clearSession();
    },
  }), [clearSession, query.data, query.isLoading, query.refetch, signedOut]);
  return <UserContext.Provider value={value}>{children}</UserContext.Provider>;
}

export function useUser() {
  const value = useContext(UserContext);
  if (!value) throw new Error('useUser must be used inside UserProvider');
  return value;
}
