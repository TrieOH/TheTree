import "./tailwind.css";

import { SignIn } from "./SignIn/SignIn";
import { SignUp } from "./SignUp/SignUp";
import { Copyright } from "./Extra/Copyright";
import { BasicLogoutButton } from "./Logout/BasicLogoutButton";
import { Sessions } from "./Session/Sessions";
import ForgotPassword from "./ForgotPassword/ForgotPassword";
import { ResetPassword } from "./ResetPassword/ResetPassword";
import { VerifyEmail } from "./VerifyEmail/VerifyEmail";
import { ResendVerifyEmail } from "./VerifyEmail/ResendVerifyEmail";
import BasicInputField from "./Form/BasicInputField"

import { ModernSignIn } from "./Modern/ModernSignIn";
import { ModernSignUp } from "./Modern/ModernSignUp";
import { ModernForgotPassword } from "./Modern/ModernForgotPassword";
import { ModernResetPassword } from "./Modern/ModernResetPassword";
import { ModernVerifyEmail } from "./Modern/ModernVerifyEmail";
import { ModernAuth } from "./Modern/ModernAuth";

export {
  SignIn,
  SignUp,
  ModernSignIn,
  ModernSignUp,
  ModernForgotPassword,
  ModernResetPassword,
  ModernVerifyEmail,
  ModernAuth,
  BasicLogoutButton,
  Copyright,
  Sessions,
  BasicInputField,
  ForgotPassword,
  ResetPassword,
  VerifyEmail,
  ResendVerifyEmail
};
