import React, { createContext, useCallback, useEffect, useState } from "react";
import { Hub } from "aws-amplify/utils";
import { fetchAuthSession, getCurrentUser, signOut } from "aws-amplify/auth";
import { useAtom } from "jotai";
import { userAtom, errorAtom } from "../atoms";

export const AuthContext = createContext<{
  checkAuthState: () => Promise<void>;
  getAuthToken: () => Promise<string | null>;
  handlePostSignIn: () => Promise<void>;
  handleSignOut: () => Promise<void>;
  isAuth: boolean;
  isLoading: boolean;
}>({
  checkAuthState: async () => {},
  getAuthToken: async () => null,
  handlePostSignIn: async () => {},
  handleSignOut: async () => {},
  isAuth: false,
  isLoading: true,
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [, setUser] = useAtom(userAtom);
  const [, setError] = useAtom(errorAtom);
  const [isLoading, setIsLoading] = useState(true);
  const [isAuth, setIsAuth] = useState(false);

  const getAuthToken = useCallback(async (): Promise<string | null> => {
    try {
      const session = await fetchAuthSession();
      return session.tokens?.accessToken?.toString() ?? null;
    } catch (error) {
      console.error("Error getting auth token: ", error);
      setError(error);
      return null;
    }
  }, [setError]);

  const handlePostSignIn = useCallback(async (): Promise<void> => {
    try {
      // Only call getCurrentUser if we have an active session
      const session = await fetchAuthSession();
      console.log("session: ", session);
      if (session.tokens?.accessToken) {
        const currentUser = await getCurrentUser();
        setUser(currentUser);
        setIsAuth(true);
      } else {
        throw new Error("No active session");
      }
    } catch (error) {
      console.error("Error in post-sign in process:", error);
      setError(error);
      setIsAuth(false);
      setUser(null);
    }
  }, [setUser, setError]);

  const handleSignOut = useCallback(async (): Promise<void> => {
    try {
      await signOut({ global: true });
      setUser(null);
      setIsAuth(false);
    } catch (error) {
      console.error("Error signing out:", error);
      setError(error);
    }
  }, [setUser, setError]);

  const checkAuthState = useCallback(async () => {
    setIsLoading(true);
    try {
      // Check for an active session without throwing an error if not authenticated
      const session = await fetchAuthSession();
      if (session.tokens?.accessToken) {
        const currentUser = await getCurrentUser();
        setUser(currentUser);
        setIsAuth(true);
      } else {
        setIsAuth(false);
        setUser(null);
      }
    } catch (error) {
      console.error("Error checking auth state: ", error);
      setIsAuth(false);
      setUser(null);
      setError(error);
    } finally {
      setIsLoading(false);
    }
  }, [setUser, setError]);

  useEffect(() => {
    checkAuthState();

    const unsubscribe = Hub.listen("auth", ({ payload }) => {
      switch (payload.event) {
        case "signInWithRedirect":
          checkAuthState();
          break;
        case "signInWithRedirect_failure":
          setError("An error has occurred during the OAuth flow.");
          setIsAuth(false);
          setIsLoading(false);
          break;
        case "signedOut":
          setUser(null);
          setIsAuth(false);
          setIsLoading(false);
          break;
      }
    });

    return unsubscribe;
  }, [checkAuthState, setError, setUser]);

  return (
    <AuthContext.Provider
      value={{
        checkAuthState,
        getAuthToken,
        handlePostSignIn,
        handleSignOut,
        isAuth,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
