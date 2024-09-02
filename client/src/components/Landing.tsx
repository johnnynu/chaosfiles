import React, { useEffect } from "react";
import { useAuth } from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Loader2 } from "lucide-react";

const LandingPage: React.FC = () => {
  const { signIn, isAuth, isLoading, checkAuthState } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuth = async () => {
      await checkAuthState();
      if (isAuth) {
        navigate("/dashboard");
      }
    };
    checkAuth();
  }, [checkAuthState, isAuth, navigate]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Loader2 className="w-8 h-8 animate-spin" />
      </div>
    );
  }

  if (isAuth) {
    return null;
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>Welcome to ChaosFiles</CardTitle>
          <CardDescription>
            Manage your files with ease and security.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="mb-4 text-sm text-gray-600">
            Sign in to access your files, upload new ones, and collaborate with
            others.
          </p>
        </CardContent>
        <CardFooter>
          <Button className="w-full" onClick={signIn} disabled={isLoading}>
            {isLoading ? (
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
            ) : (
              <svg
                className="w-4 h-4 mr-2"
                aria-hidden="true"
                focusable="false"
                data-prefix="fab"
                data-icon="google"
                role="img"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 488 512"
              >
                <path
                  fill="currentColor"
                  d="M488 261.8C488 403.3 391.1 504 248 504 110.8 504 0 393.2 0 256S110.8 8 248 8c66.8 0 123 24.5 166.3 64.9l-67.5 64.9C258.5 52.6 94.3 116.6 94.3 256c0 86.5 69.1 156.6 153.7 156.6 98.2 0 135-70.4 140.8-106.9H248v-85.3h236.1c2.3 12.7 3.9 24.9 3.9 41.4z"
                ></path>
              </svg>
            )}
            Sign In With Google
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};

export default LandingPage;
