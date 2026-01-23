import type { HeaderConfigI } from "@/types/route-types";
import { Link } from "@tanstack/react-router";
import { ShadowButton } from "./Basic/ShadowButton";

import { LogIn, Menu, X } from 'lucide-react';
import { cn } from "@/lib/utils";
import { useState } from "react";

export default function Header({  }: HeaderConfigI) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  return (
    <header className="relative">
      <div 
        className={cn(
          "flex justify-between items-center border-b-2 border-b-border px-6 py-4",
          "bg-background/80 backdrop-blur-sm"
        )}
      >
        <h2 className="text-2xl font-semibold text-foreground select-none md:block hidden">TrieAuth</h2>
        <div
          className={cn(
            "md:hidden block active:scale-95 active:translate-y-px",
            "select-none cursor-pointer transition-transform duration-100 ease-out"
          )}
          onClick={() => setIsMenuOpen(v => !v)}
        >
          { isMenuOpen ? <X size={24} /> : <Menu size={24} /> }
        </div>
        <div className="md:flex justify-center gap-6 text-lg hidden">
          <Link
            to="/"
            className="custom-underline select-none"
          >
            Features
          </Link>
          <Link
            to="/"
            className="custom-underline select-none"
          >
            Pricing
          </Link>
          <Link
            to="/"
            className="custom-underline select-none"
          >
            Docs
          </Link>
        </div>
        <div>
          <ShadowButton 
            value="Authenticate" 
            leftIcon={ <LogIn size={24}/> }
            className="xs:flex hidden"
          />
          <ShadowButton
            leftIcon={ <LogIn size={16}/> }
            className="xs:hidden flex"
          />
        </div>
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
          className="custom-underline select-none"
        >
          Features
        </Link>
        <Link
          to="/"
          className="custom-underline select-none"
        >
          Pricing
        </Link>
        <Link
          to="/"
          className="custom-underline select-none"
        >
          Docs
        </Link>
      </div>
    </header>
  )
}