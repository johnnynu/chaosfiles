import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuth } from "../hooks/useAuth";
import Chaos from "@/assets/doodle-stars.png";

const Navbar: React.FC = () => {
  const { user, signOut } = useAuth();
  const navigate = useNavigate();

  const UserMenu = () => (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost">Hi, {user?.username}</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => navigate("/files")}>
          File Manager
        </DropdownMenuItem>
        <DropdownMenuItem onClick={signOut}>Sign out</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );

  return (
    <nav className="flex justify-between items-center mb-8">
      <div className="flex items-center">
        <Link to="/dashboard" className="flex items-center">
          <img src={Chaos} alt="ChaosFiles Logo" className="h-10 mr-4" />
          <h1 className="text-3xl font-bold">ChaosFiles</h1>
        </Link>
      </div>
      <UserMenu />
    </nav>
  );
};

export default Navbar;
