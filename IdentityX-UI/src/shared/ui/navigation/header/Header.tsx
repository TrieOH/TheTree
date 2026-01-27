import { Link, useNavigate, useRouteContext } from "@tanstack/react-router";
import { LogIn, Menu, User, X } from 'lucide-react';
import { useState } from "react";
import { cn } from "@/shared/lib/utils";
import type { HeaderConfigI } from "@/shared/types/route-types";
import { ShadowButton } from "../../buttons/ShadowButton";

export default function Header({  }: HeaderConfigI) {
  const navigate = useNavigate();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const { auth } = useRouteContext({ from: '__root__' })
  return (
    <header className="relative">
      <div 
        className={cn(
          "flex justify-between items-center border-b-2 border-b-border px-6 py-4",
          "bg-background/80 backdrop-blur-sm select-none"
        )}
      >
        <h2 className="text-2xl font-semibold text-foreground md:block hidden">TrieAuth</h2>
        <button
          type="button"
          className={cn(
            "md:hidden block active:scale-95 active:translate-y-px",
            "cursor-pointer transition-transform duration-100 ease-out"
          )}
          onClick={() => setIsMenuOpen(v => !v)}
        >
          { isMenuOpen ? <X size={24} /> : <Menu size={24} /> }
        </button>
        <div className="md:flex justify-center gap-6 text-lg hidden">
          <Link
            to="/"
            className="custom-underline"
          >
            Features
          </Link>
          <Link
            to="/"
            className="custom-underline"
          >
            Pricing
          </Link>
          <Link
            to="/"
            className="custom-underline"
          >
            Docs
          </Link>
        </div>
        {!auth?.isAuthenticated ?
          <div>
            <ShadowButton 
              value="Authenticate" 
              leftIcon={ <LogIn size={24}/> }
              className="xs:flex hidden"
              onClick={() => navigate({to: "/auth"})}
            />
            <ShadowButton
              leftIcon={ <LogIn size={16}/> }
              className="xs:hidden flex"
              onClick={() => navigate({to: "/auth"})}
            />
          </div>
          : <div>
            <ShadowButton 
              value="Dashboard" 
              leftIcon={ <User size={24}/> }
              className="xs:flex hidden"
              onClick={() => navigate({to: "/projects"})}
            />
            <ShadowButton
              leftIcon={ <User size={16}/> }
              className="xs:hidden flex"
              onClick={() => navigate({to: "/projects"})}
            />
          </div>
        }
      </div>
      <div 
        className={cn(
          "absolute w-full flex flex-col md:hidden justify-center items-center gap-4",
          "text-lg border-b-2 border-b-border py-4 bg-background/80 backdrop-blur-sm",
          !isMenuOpen && "hidden"
        )}
      >
        <Link
          to="/"
          className="custom-underline"
        >
          Features
        </Link>
        <Link
          to="/"
          className="custom-underline"
        >
          Pricing
        </Link>
        <Link
          to="/"
          className="custom-underline"
        >
          Docs
        </Link>
      </div>
    </header>
  )
}