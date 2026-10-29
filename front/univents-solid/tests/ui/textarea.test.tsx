import { cleanup, fireEvent, render, screen } from "@solidjs/testing-library";
import { createSignal, flush } from "solid-js";
import { afterEach, describe, expect, it, vi } from "vitest";

import { Textarea } from "@trieoh/ui-solid";

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
});

const type = (element: HTMLElement, value: string) => {
  (element as HTMLTextAreaElement).value = value;
  fireEvent.input(element, { target: { value } });
  flush();
};

const set = <T,>(setter: (value: T) => void, value: T) => {
  setter(value);
  flush();
};

// The field is border-box: the height covers content plus borders (and, for the
// cap, padding). jsdom gives textareas non-zero defaults, so they are read
// instead of assumed.
const sizes = (element: HTMLElement) => {
  const style = getComputedStyle(element);
  const px = (value: string) => parseFloat(value) || 0;
  return {
    borders: px(style.borderTopWidth) + px(style.borderBottomWidth),
    padding: px(style.paddingTop) + px(style.paddingBottom),
  };
};

// jsdom has no layout, so the content height autosize reads is supplied.
const withScrollHeight = (height: () => number) =>
  vi
    .spyOn(HTMLElement.prototype, "scrollHeight", "get")
    .mockImplementation(height);

describe("Textarea", () => {
  it("keeps its original defaults", () => {
    render(() => <Textarea aria-label="Descrição" />);
    const field = screen.getByLabelText("Descrição");

    expect(field).toHaveAttribute("rows", "3");
    expect(field).not.toHaveAttribute("maxlength");
    expect(field.className).toContain("resize-y");
    expect(field.style.height).toBe("");
    // No counter unless asked for.
    expect(screen.queryByText(/caracteres|\/\d+/)).toBeNull();
  });

  it("limits the length natively", () => {
    render(() => <Textarea aria-label="Bio" maxlength={120} />);

    expect(screen.getByLabelText("Bio")).toHaveAttribute("maxlength", "120");
  });

  it("shows the counter on its own when there is a limit, unless turned off", () => {
    render(() => <Textarea aria-label="Bio" maxlength={50} />);
    expect(screen.getByText("0/50")).toBeInTheDocument();
    cleanup();

    render(() => <Textarea aria-label="Bio" maxlength={50} showCount={false} />);
    expect(screen.queryByText("0/50")).toBeNull();
    expect(screen.getByLabelText("Bio")).not.toHaveAttribute("aria-describedby");
  });

  it("flags text that is already over a newly set limit", () => {
    render(() => <Textarea aria-label="Bio" maxlength={5} value="comprimento maior" />);

    // The browser never truncates existing text, so the counter is what shows it.
    expect(screen.getByText("17/5").className).toContain("text-destructive");
  });

  it.each([
    ["none", "resize-none"],
    ["horizontal", "resize-x"],
    ["both", "resize"],
  ] as const)("applies resize=%s", (resize, expected) => {
    render(() => <Textarea aria-label="Texto" resize={resize} />);

    expect(screen.getByLabelText("Texto").className).toContain(expected);
  });

  describe("character count", () => {
    it("counts typed characters against the limit", () => {
      render(() => <Textarea aria-label="Bio" showCount maxlength={10} />);
      expect(screen.getByText("0/10")).toBeInTheDocument();

      type(screen.getByLabelText("Bio"), "olá");
      expect(screen.getByText("3/10")).toBeInTheDocument();
    });

    it("counts without a limit as plain characters", () => {
      render(() => <Textarea aria-label="Bio" showCount />);

      type(screen.getByLabelText("Bio"), "abcd");
      expect(screen.getByText("4 caracteres")).toBeInTheDocument();
    });

    it("follows a controlled value", () => {
      const [value, setValue] = createSignal("abc");
      render(() => (
        <Textarea aria-label="Bio" showCount maxlength={20} value={value()} />
      ));
      expect(screen.getByText("3/20")).toBeInTheDocument();

      set(setValue, "abcdef");
      expect(screen.getByText("6/20")).toBeInTheDocument();
    });

    it("emphasises the counter near and at the limit", () => {
      const [value, setValue] = createSignal("1234");
      render(() => (
        <Textarea aria-label="Bio" showCount maxlength={10} value={value()} />
      ));
      const counter = () => screen.getByText(/\/10$/);

      expect(counter().className).not.toContain("text-foreground");

      set(setValue, "123456789");
      expect(counter().className).toContain("text-foreground");

      set(setValue, "1234567890");
      expect(counter().className).toContain("text-destructive");
    });

    it("describes the field by its counter and keeps the caller's description", () => {
      render(() => (
        <Textarea aria-label="Bio" showCount aria-describedby="ajuda" />
      ));
      const field = screen.getByLabelText("Bio");
      const counter = screen.getByText("0 caracteres");

      expect(field.getAttribute("aria-describedby")).toBe(`ajuda ${counter.id}`);
    });
  });

  describe("autosize", () => {
    it("grows to its content and stops resizing by hand", () => {
      let content = 60;
      withScrollHeight(() => content);
      render(() => <Textarea aria-label="Nota" autosize />);
      const field = screen.getByLabelText("Nota") as HTMLTextAreaElement;

      const { borders } = sizes(field);
      expect(field.className).toContain("resize-none");
      expect(field.style.height).toBe(`${60 + borders}px`);

      content = 140;
      type(field, "uma\nbem\nlonga\nnota");
      expect(field.style.height).toBe(`${140 + borders}px`);
      expect(field.style.overflowY).toBe("hidden");
    });

    it("caps the height at maxRows and then scrolls", () => {
      withScrollHeight(() => 500);
      render(() => <Textarea aria-label="Nota" autosize maxRows={2} />);
      const field = screen.getByLabelText("Nota") as HTMLTextAreaElement;

      // jsdom reports no line height, so a line is the 14px*1.5 fallback.
      const { borders, padding } = sizes(field);
      expect(field.style.height).toBe(`${42 + padding + borders}px`);
      expect(field.style.overflowY).toBe("auto");
    });

    it("re-fits when the value changes from outside", () => {
      let content = 50;
      withScrollHeight(() => content);
      const [value, setValue] = createSignal("");
      render(() => <Textarea aria-label="Nota" autosize value={value()} />);
      const field = screen.getByLabelText("Nota") as HTMLTextAreaElement;
      const { borders } = sizes(field);
      expect(field.style.height).toBe(`${50 + borders}px`);

      content = 90;
      set(setValue, "texto carregado de outro lugar");
      expect(field.style.height).toBe(`${90 + borders}px`);
    });

    it("does nothing to the height when off", () => {
      withScrollHeight(() => 300);
      render(() => <Textarea aria-label="Nota" />);
      const field = screen.getByLabelText("Nota") as HTMLTextAreaElement;

      type(field, "x");
      expect(field.style.height).toBe("");
    });
  });

  it("still forwards input events and the element ref", () => {
    const onInput = vi.fn();
    let element: HTMLTextAreaElement | undefined;
    render(() => (
      <Textarea
        aria-label="Nota"
        autosize
        showCount
        onInput={onInput}
        ref={(el) => {
          element = el;
        }}
      />
    ));

    type(screen.getByLabelText("Nota"), "oi");

    expect(onInput).toHaveBeenCalledTimes(1);
    expect(element).toBe(screen.getByLabelText("Nota"));
  });
});
