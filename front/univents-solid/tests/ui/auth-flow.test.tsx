import { render, screen, fireEvent } from "@solidjs/testing-library";
import { describe, expect, it } from "vitest";
import { AuthProvider, ModernAuth } from "@trieoh/identityx-sdk-ts-solid";

describe("ModernAuth flow navigation", () => {
  it("transitions between signin, signup and forgot-password when clicking toggle buttons", async () => {
    render(() => (
      <AuthProvider isProjectMode={false} waitSession={false}>
        <ModernAuth initialView="signin" />
      </AuthProvider>
    ));

    // Initially in signin
    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("Bem-vindo de volta");
    expect(screen.getByRole("button", { name: "Cadastre-se" })).toBeInTheDocument();

    // Click "Cadastre-se"
    const signUpToggle = screen.getByRole("button", { name: "Cadastre-se" });
    fireEvent.click(signUpToggle);
    await Promise.resolve();

    // Expect to transition to signup
    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("Criar sua conta");
    expect(screen.getByRole("button", { name: "Faça login" })).toBeInTheDocument();

    // Click "Faça login" to return to signin
    const signInToggle = screen.getByRole("button", { name: "Faça login" });
    fireEvent.click(signInToggle);
    await Promise.resolve();

    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("Bem-vindo de volta");

    // Click "Esqueceu a senha?"
    const forgotPasswordLink = screen.getByRole("button", { name: "Esqueceu a senha?" });
    fireEvent.click(forgotPasswordLink);
    await Promise.resolve();

    // Expect to transition to forgot-password
    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("Recuperar senha");
    expect(screen.getByRole("button", { name: "Voltar ao login" })).toBeInTheDocument();

    // Click "Voltar ao login"
    const backToLoginToggle = screen.getByRole("button", { name: "Voltar ao login" });
    fireEvent.click(backToLoginToggle);
    await Promise.resolve();

    // Expect to be back at signin
    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("Bem-vindo de volta");
  });
});
