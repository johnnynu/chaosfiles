import { useAtom } from "jotai";
import { userAtom, errorAtom } from "../atoms";
import { signInWithRedirect, signOut } from "aws-amplify/auth";
import { AuthContext } from "../contexts/AuthProvider";
import { useContext } from "react";

export const useAuth = () => {
  const [user] = useAtom(userAtom);
  const [error] = useAtom(errorAtom);
  const { getAuthToken, handlePostSignIn, isAuth, isLoading } =
    useContext(AuthContext);

  const signIn = async () => {
    try {
      await signInWithRedirect({ provider: "Google" });
    } catch (error) {
      console.error("Error signing in:", error);
    }
  };

  const handleSignOut = async () => {
    try {
      await signOut({ global: true });
    } catch (error) {
      console.error("Error signing out:", error);
    }
  };

  return {
    user,
    signIn,
    signOut: handleSignOut,
    getAuthToken,
    handlePostSignIn,
    isAuth,
    isLoading,
    error,
  };
};
