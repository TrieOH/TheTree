import { cleanup, render, screen, within } from "@solidjs/testing-library";
import { MultiStepStepper } from "@trieoh/ui-solid";
import { afterEach, describe, expect, it } from "vitest";

afterEach(cleanup);

const steps = [
  { id: "edition", title: "Edição" },
  { id: "schedule", title: "Cronograma" },
  { id: "summary", title: "Resumo" },
];

describe("MultiStepStepper layout", () => {
  it("uses a horizontal list and a compact variant for narrow containers", () => {
    render(() => <MultiStepStepper steps={steps} currentStep={0} />);
    const nav = screen.getByRole("navigation", { name: "Etapas" });
    const list = within(nav).getByRole("list", { name: "" });
    const compact = within(nav).getByRole("list", { name: "Progresso das etapas" }).parentElement;

    expect(nav).toHaveClass("@container/stepper");
    expect(list).toHaveClass("flex", "@max-[28rem]/stepper:hidden");
    expect(compact).toHaveClass("hidden", "@max-[28rem]/stepper:block");
  });

  it("keeps explicit horizontal orientation at narrow widths", () => {
    render(() => <MultiStepStepper steps={steps} currentStep={0} orientation="horizontal" />);
    const list = screen.getByRole("list", { name: "" });
    expect(list).toHaveClass("flex");
    expect(list).not.toHaveClass("@max-[28rem]/stepper:hidden");
  });

  it("lays out explicit vertical steps as a grid", () => {
    render(() => <MultiStepStepper steps={steps} currentStep={0} orientation="vertical" />);
    const list = screen.getByRole("list", { name: "" });
    expect(list).toHaveClass("grid");
    expect(list).not.toHaveClass("flex", "@max-[28rem]/stepper:hidden");
  });
});
