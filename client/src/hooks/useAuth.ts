import { useAtom } from "jotai";
import { userAtom, errorAtom } from "../atoms";
import { signInWithRedirect } from "aws-amplify/auth";
import { AuthContext } from "../contexts/AuthProvider";
import { useContext, useCallback } from "react";

export const useAuth = () => {
  const [user] = useAtom(userAtom);
  const [error] = useAtom(errorAtom);
  const {
    getAuthToken,
    handlePostSignIn,
    handleSignOut,
    checkAuthState,
    isAuth,
    isLoading,
  } = useContext(AuthContext);

  const signIn = useCallback(async () => {
    try {
      await signInWithRedirect({ provider: "Google" });
    } catch (error) {
      console.error("Error signing in:", error);
    }
  }, []);

  return {
    user,
    checkAuthState,
    signIn,
    signOut: handleSignOut,
    getAuthToken,
    handlePostSignIn,
    isAuth,
    isLoading,
    error,
  };
};
