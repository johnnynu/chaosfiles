// src/components/LandingPage.tsx
import React from "react";
import { useAuth } from "../hooks/useAuth";

const Dashboard: React.FC = () => {
  const { user, signOut } = useAuth();

  return (
    <div className="sign-in-container">
      <h2>{user && `Welcome to your dashboard, ${user.username}!`}</h2>
      <button onClick={signOut}>
        {user ? "Sign Out" : "Sign In with Google"}
      </button>
    </div>
  );
};

export default Dashboard;
