import { cleanup, fireEvent, render, screen, waitFor } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { EventMemberWithEmailI } from "@/features/events/model/member";
import { AdminAddMemberCard } from "@/features/events/ui/AdminAddMemberCard";
import { AdminEventMemberCard } from "@/features/events/ui/AdminEventMemberCard";
import { ManageEventMemberDialog } from "@/features/events/ui/ManageEventMemberDialog";
import { RemoveEventMemberDialog } from "@/features/events/ui/RemoveEventMemberDialog";

afterEach(cleanup);

const mockMember: EventMemberWithEmailI = {
  id: "member-1",
  event_id: "event-1",
  user_id: "user-abc-12345678",
  role: "admin",
  email: "maria@example.com",
  created_at: "2026-03-10T12:00:00.000Z",
};

describe("AdminAddMemberCard", () => {
  it("renders add member card and triggers onAdd callback", () => {
    const onAdd = vi.fn();
    render(() => <AdminAddMemberCard onAdd={onAdd} animate={false} />);

    const btn = screen.getByRole("button", { name: /Adicionar membro/i });
    expect(btn).toBeInTheDocument();
    fireEvent.click(btn);
    expect(onAdd).toHaveBeenCalledTimes(1);
  });
});

describe("AdminEventMemberCard", () => {
  it("renders member email, role subtitle, and initials, and allows removal", () => {
    const onRemove = vi.fn();
    render(() => (
      <AdminEventMemberCard member={mockMember} onRemove={onRemove} animate={false} />
    ));

    expect(screen.getByText("maria@example.com")).toBeInTheDocument();
    expect(screen.getByText("Administrador")).toBeInTheDocument();
    expect(screen.getByText("MA")).toBeInTheDocument();

    const removeBtn = screen.getByLabelText("Remover maria@example.com");
    fireEvent.click(removeBtn);
    expect(onRemove).toHaveBeenCalledWith(mockMember);
  });

  it("renders member avatar image when pfp_url is provided", () => {
    render(() => (
      <AdminEventMemberCard
        member={{ ...mockMember, pfp_url: "https://example.com/avatar.png" }}
        onRemove={() => undefined}
        animate={false}
      />
    ));

    const img = screen.getByRole("img", { name: "Avatar de maria@example.com" });
    expect(img).toHaveAttribute("src", "https://example.com/avatar.png");
  });

  it("copies user id to clipboard when clicked", async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    Object.assign(navigator, {
      clipboard: { writeText },
    });

    render(() => (
      <AdminEventMemberCard member={mockMember} onRemove={() => undefined} animate={false} />
    ));

    const copyBtn = screen.getByTitle("Clique para copiar o ID completo");
    fireEvent.click(copyBtn);
    expect(writeText).toHaveBeenCalledWith("user-abc-12345678");
  });
});

describe("ManageEventMemberDialog", () => {
  it("submits the entered email and selected role", async () => {
    const onSubmit = vi.fn().mockResolvedValue(true);
    const onOpenChange = vi.fn();

    render(() => (
      <ManageEventMemberDialog
        open={true}
        onOpenChange={onOpenChange}
        onSubmit={onSubmit}
      />
    ));

    expect(screen.getByRole("heading", { name: "Adicionar membro" })).toBeInTheDocument();

    const emailInput = screen.getByPlaceholderText("membro@exemplo.com");
    fireEvent.input(emailInput, { target: { value: "novo@membro.com" } });

    // Select role "Administrador"
    const adminRoleBtn = screen.getByRole("button", { name: /Administrador/i });
    fireEvent.click(adminRoleBtn);

    const submitBtn = screen.getByRole("button", { name: "Adicionar membro" });
    await waitFor(() => {
      expect(submitBtn).not.toBeDisabled();
    });

    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        email: "novo@membro.com",
        role: "admin",
      });
    });
  });
});

describe("RemoveEventMemberDialog", () => {
  it("requires email confirmation before allowing removal", async () => {
    const onRemove = vi.fn().mockResolvedValue(true);
    const onOpenChange = vi.fn();

    render(() => (
      <RemoveEventMemberDialog
        open={true}
        onOpenChange={onOpenChange}
        member={mockMember}
        onRemove={onRemove}
      />
    ));

    expect(screen.getByText("Remover membro?")).toBeInTheDocument();

    const removeBtn = screen.getByRole("button", { name: "Remover membro" });
    expect(removeBtn).toBeDisabled();

    const confirmInput = screen.getByPlaceholderText("maria@example.com");
    fireEvent.input(confirmInput, { target: { value: "maria@example.com" } });

    await waitFor(() => {
      expect(removeBtn).not.toBeDisabled();
    });

    fireEvent.click(removeBtn);

    await waitFor(() => {
      expect(onRemove).toHaveBeenCalledWith("user-abc-12345678", "maria@example.com");
    });
  });
});
