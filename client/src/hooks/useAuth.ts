import { useAtom } from "jotai";
import { userAtom, errorAtom } from "../atoms";
import { signInWithRedirect, signOut } from "aws-amplify/auth";

export const useAuth = () => {
  const [user] = useAtom(userAtom);
  const [error] = useAtom(errorAtom);

  const signIn = async () => {
    try {
      await signInWithRedirect({ provider: "Google" });
    } catch (error) {
      console.error("Error signing in:", error);
    }
  };

  const handleSignOut = async () => {
    try {
      await signOut();
    } catch (error) {
      console.error("Error signing out:", error);
    }
  };

  return {
    user,
    signIn,
    signOut: handleSignOut,
    error,
  };
};
