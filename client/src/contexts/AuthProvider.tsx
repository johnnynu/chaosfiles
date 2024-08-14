// src/contexts/AuthContext.tsx
import React, { createContext, useCallback, useEffect, useState } from "react";
import { Hub } from "aws-amplify/utils";
import { fetchAuthSession, getCurrentUser } from "aws-amplify/auth";
import { useAtom } from "jotai";
import { userAtom, errorAtom } from "../atoms";

export const AuthContext = createContext<{
  getAuthToken: () => Promise<string | null>;
  handlePostSignIn: () => Promise<void>;
  isAuth: boolean;
  isLoading: boolean;
}>({
  getAuthToken: async () => null,
  handlePostSignIn: async () => {},
  isAuth: false,
  isLoading: false,
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [, setUser] = useAtom(userAtom);
  const [, setError] = useAtom(errorAtom);
  const [isLoading, setIsLoading] = useState(false);
  const [isAuth, setIsAuth] = useState(false);

  const getAuthToken = useCallback(async (): Promise<string | null> => {
    try {
      const session = await fetchAuthSession();
      const accessToken = session.tokens?.accessToken;
      if (accessToken) {
        return accessToken.toString();
      } else {
        console.error("Id token is undefined");
        return null;
      }
    } catch (error) {
      console.error("Error getting auth token: ", error);
      setError(error);
      return null;
    }
  }, [setError]);

  const handlePostSignIn = useCallback(async (): Promise<void> => {
    console.log("HandlePostSignIn called");
    try {
      const currentUser = await getCurrentUser();
      setUser(currentUser);
      console.log("Current user: ", currentUser);

      const token = await getAuthToken();
      if (!token) {
        throw new Error("Failed to get auth token");
      }
      console.log("Auth token obtained: ", token);

      const res = await fetch(
        "https://4j1h7lzpf5.execute-api.us-east-2.amazonaws.com/dev/chaosfiles-signin",
        {
          method: "POST",
          headers: {
            Authorization: token,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({}),
        }
      );

      if (!res.ok) {
        throw new Error("Failed to sign in user");
      }

      const data = await res.json();
      console.log("Sign in successful:", data);
    } catch (error) {
      console.error("Error in post-sign in process:", error);
      setError(error);
    }
  }, [getAuthToken, setUser, setError]);

  const handleAuthStateChange = useCallback(async () => {
    setIsLoading(true);
    try {
      const currentUser = await getCurrentUser();
      setUser(currentUser);
      await handlePostSignIn();
      setIsAuth(true);
    } catch (error) {
      console.error("Error handling auth state change: ", error);
      setError(error);
      setUser(null);
      setIsAuth(false);
    } finally {
      setIsLoading(false);
    }
  }, [setUser, setError, handlePostSignIn]);

  useEffect(() => {
    const checkAuthState = async () => {
      try {
        await getCurrentUser();
        setIsAuth(true);
      } catch {
        setIsAuth(false);
      } finally {
        setIsLoading(false);
      }
    };

    checkAuthState();

    const unsubscribe = Hub.listen("auth", ({ payload }) => {
      switch (payload.event) {
        case "signInWithRedirect":
          handleAuthStateChange();
          break;
        case "signInWithRedirect_failure":
          setError("An error has occurred during the OAuth flow.");
          break;
        case "signedOut":
          setUser(null);
          break;
      }
    });

    return unsubscribe;
  });

  return (
    <AuthContext.Provider
      value={{ getAuthToken, handlePostSignIn, isAuth, isLoading }}
    >
      {children}
    </AuthContext.Provider>
  );
};
