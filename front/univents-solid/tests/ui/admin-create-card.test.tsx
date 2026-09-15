import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { afterEach, describe, expect, it, vi } from "vitest";

import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

afterEach(cleanup);

describe("AdminCreateCard", () => {
  it("renders title, description, and icon, and handles clicks", () => {
    const onClick = vi.fn();
    render(() => (
      <AdminCreateCard
        title="Novo item"
        description="Descrição do item"
        icon={(props) => <span data-testid="test-icon" class={props.class}>icon</span>}
        onClick={onClick}
        animate={false}
      />
    ));

    const btn = screen.getByRole("button", { name: /Novo item/i });
    expect(btn).toBeInTheDocument();
    expect(screen.getByText("Novo item")).toBeInTheDocument();
    expect(screen.getByText("Descrição do item")).toBeInTheDocument();
    expect(screen.getByTestId("test-icon")).toBeInTheDocument();

    // Check dashed border classes
    expect(btn.className).toContain("border-dashed");

    fireEvent.click(btn);
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it("supports stacked layout for taller cards", () => {
    const onClick = vi.fn();
    render(() => (
      <AdminCreateCard
        title="Nova edição"
        description="Criar uma nova edição"
        icon={<span>+</span>}
        onClick={onClick}
        layout="stacked"
        animate={false}
      />
    ));

    const btn = screen.getByRole("button", { name: /Nova edição/i });
    expect(btn).toBeInTheDocument();
    expect(btn.className).toContain("flex-col");
    expect(btn.className).toContain("justify-between");
  });

  it("supports JSX icon element directly", () => {
    const onClick = vi.fn();
    render(() => (
      <AdminCreateCard
        title="Item com JSX"
        description="Outra descrição"
        icon={<span data-testid="direct-icon">★</span>}
        onClick={onClick}
        animate={false}
      />
    ));

    expect(screen.getByTestId("direct-icon")).toBeInTheDocument();
  });
});
