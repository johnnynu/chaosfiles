// src/components/LandingPage.tsx
import React, { useEffect } from "react";
import { useAuth } from "../hooks/useAuth";
import { useNavigate } from "react-router-dom";

const LandingPage: React.FC = () => {
  const { signIn, isAuth, isLoading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuth && !isLoading) {
      navigate("/dashboard");
    }
  }, [isAuth, isLoading, navigate]);

  return (
    <div className="landing-page">
      <h1>Welcome to ChaosFiles</h1>
      <p>Manage your files with ease.</p>
      <button onClick={signIn}>Sign In With Google</button>
    </div>
  );
};

export default LandingPage;
