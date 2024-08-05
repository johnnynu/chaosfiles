// src/contexts/AuthContext.tsx
import React, { useCallback, useEffect } from "react";
import { Hub } from "aws-amplify/utils";
import { getCurrentUser } from "aws-amplify/auth";
import { useAtom } from "jotai";
import { userAtom, errorAtom } from "../atoms";

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [, setUser] = useAtom(userAtom);
  const [, setError] = useAtom(errorAtom);

  const getUser = useCallback(async (): Promise<void> => {
    try {
      const currentUser = await getCurrentUser();
      setUser(currentUser);
    } catch (error) {
      console.error(error);
      setUser(null);
    }
  }, [setUser]);

  useEffect(() => {
    const unsubscribe = Hub.listen("auth", ({ payload }) => {
      switch (payload.event) {
        case "signInWithRedirect":
          getUser();
          break;
        case "signInWithRedirect_failure":
          setError("An error has occurred during the OAuth flow.");
          break;
        case "signedOut":
          setUser(null);
          break;
      }
    });

    getUser();
    return unsubscribe;
  }, [getUser, setError, setUser]);

  return <>{children}</>;
};
