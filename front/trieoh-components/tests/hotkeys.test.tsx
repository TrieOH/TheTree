import { KbdExample } from "../src/examples/InteractionExamples";
import { afterEach, beforeEach, expect, it, vi } from "vitest";
import { createSignal } from "solid-js";
import { render } from "@solidjs/web";
import { captureArtifact, expectNoDiagnostics } from "@solidjs/diagnostics";
import {
  DropdownMenu,
  Button,
  createHotkeys,
  CommandPalette,
  Timeline,
  Kbd,
  parseHotkey,
  type Hotkey,
  type HotkeyInput,
} from "@trieoh/ui-solid";
let host: HTMLDivElement;
const disposers: (() => void)[] = [];
const settle = () => new Promise((resolve) => setTimeout(resolve, 25));
const press = (
  target: EventTarget,
  key: string,
  extra: KeyboardEventInit = {},
) => {
  const event = new KeyboardEvent("keydown", {
    key,
    bubbles: true,
    cancelable: true,
    ...extra,
  });
  target.dispatchEvent(event);
  return event;
};
beforeEach(() => {
  host = document.createElement("div");
  document.body.append(host);
});
afterEach(() => {
  for (const dispose of disposers.splice(0).reverse()) dispose();
  host.remove();
});
it("handles Ctrl and Command through Mod, exact modifiers, composition and repeat", async () => {
  const run = vi.fn();
  disposers.push(
    render(() => {
      createHotkeys(() => [{ hotkey: "Mod+K", callback: run }]);
      return <div />;
    }, host),
  );
  await settle();
  expect(press(document.body, "k", { ctrlKey: true }).defaultPrevented).toBe(
    true,
  );
  press(document.body, "K", { metaKey: true });
  press(document.body, "k", { ctrlKey: true, shiftKey: true });
  press(document.body, "k", { ctrlKey: true, repeat: true });
  press(document.body, "k", { ctrlKey: true, isComposing: true });
  expect(run).toHaveBeenCalledTimes(2);
});
it("ignores inputs and editable descendants, with a reactive opt-out", async () => {
  const run = vi.fn(),
    [ignore, setIgnore] = createSignal(true);
  const { artifact } = await captureArtifact(async () => {
    disposers.push(
      render(() => {
        createHotkeys(() => [{ hotkey: "Mod+K", callback: run }], {
          get ignoreInputs() {
            return ignore();
          },
        });
        return (
          <>
            <input />
            <div contenteditable="true">
              <span>Editable</span>
            </div>
          </>
        );
      }, host),
    );
    await settle();
    press(host.querySelector("input")!, "k", { ctrlKey: true });
    press(host.querySelector("span")!, "k", { metaKey: true });
    expect(run).not.toHaveBeenCalled();
    setIgnore(false);
    await settle();
    press(host.querySelector("input")!, "k", { ctrlKey: true });
    expect(run).toHaveBeenCalledTimes(1);
  });
  expectNoDiagnostics(artifact);
});
it("uses live enabled conditions and removes shortcuts when the owner unmounts", async () => {
  const run = vi.fn(),
    [enabled, setEnabled] = createSignal(false);
  const dispose = render(() => {
    createHotkeys(() => [
      { hotkey: "Mod+Shift+P", callback: run, options: { enabled: enabled() } },
    ]);
    return <div />;
  }, host);
  disposers.push(dispose);
  await settle();
  press(document.body, "P", { ctrlKey: true, shiftKey: true });
  expect(run).not.toHaveBeenCalled();
  setEnabled(true);
  await settle();
  press(document.body, "P", { ctrlKey: true, shiftKey: true });
  expect(run).toHaveBeenCalledTimes(1);
  dispose();
  press(document.body, "P", { ctrlKey: true, shiftKey: true });
  expect(run).toHaveBeenCalledTimes(1);
});
it("resolves shortcut collisions once and respects scoped registrations and handled events", async () => {
  const first = vi.fn(),
    second = vi.fn();
  let scope!: HTMLDivElement;
  disposers.push(
    render(() => {
      createHotkeys(() => [{ hotkey: "Mod+K", callback: first }]);
      return (
        <div
          ref={(element) => {
            scope = element;
          }}
          tabindex={0}
        >
          Scope
        </div>
      );
    }, host),
  );
  await settle();
  const secondHost = document.createElement("div");
  host.append(secondHost);
  const disposeSecond = render(() => {
    createHotkeys(() => [{ hotkey: "Mod+K", callback: second }], {
      scope: () => scope,
    });
    return <div />;
  }, secondHost);
  disposers.push(disposeSecond);
  await settle();
  press(document.body, "k", { ctrlKey: true });
  expect(first).toHaveBeenCalledTimes(1);
  expect(second).not.toHaveBeenCalled();
  press(scope, "k", { ctrlKey: true });
  expect(first).toHaveBeenCalledTimes(1);
  expect(second).toHaveBeenCalledTimes(1);
  const handled = new KeyboardEvent("keydown", {
    key: "k",
    ctrlKey: true,
    bubbles: true,
    cancelable: true,
  });
  handled.preventDefault();
  scope.dispatchEvent(handled);
  expect(second).toHaveBeenCalledTimes(1);
  disposeSecond();
  press(scope, "k", { ctrlKey: true });
  expect(first).toHaveBeenCalledTimes(2);
});
it("supports configurable command shortcuts and a fixed results area through filtering", async () => {
  disposers.push(
    render(
      () => (
        <CommandPalette
          shortcut="Mod+J"
          resultsHeight={360}
          items={[
            { id: "a", label: "Eventos", onSelect: () => {} },
            { id: "b", label: "Pagamentos", onSelect: () => {} },
          ]}
        />
      ),
      host,
    ),
  );
  await settle();
  press(document.body, "k", { ctrlKey: true });
  await settle();
  expect(document.querySelector('[role="combobox"]')).toBeNull();
  press(document.body, "j", { metaKey: true });
  await settle();
  const results = document.querySelector<HTMLElement>('[role="listbox"]')!;
  const initial = results.style.height;
  expect(initial).toContain("360px");
  const input = document.querySelector<HTMLInputElement>('[role="combobox"]')!;
  input.value = "inexistente";
  input.dispatchEvent(new Event("input", { bubbles: true }));
  await settle();
  expect(results.textContent).toContain("Nenhum resultado");
  expect(results.style.height).toBe(initial);
});
it("switches timeline orientation without changing item order or contents", async () => {
  const [orientation, setOrientation] = createSignal<"vertical" | "horizontal">(
    "vertical",
  );
  disposers.push(
    render(
      () => (
        <Timeline
          orientation={orientation()}
          items={[
            { id: "a", title: "Criada", date: "2026-11-19" },
            { id: "b", title: "Paga", date: "2026-11-20" },
          ]}
        />
      ),
      host,
    ),
  );
  await settle();
  const list = host.querySelector("ol")!;
  const contents = list.textContent;
  expect(list.getAttribute("tabindex")).toBeNull();
  setOrientation("horizontal");
  await settle();
  expect(list.textContent).toBe(contents);
  expect(list.getAttribute("tabindex")).toBe("0");
  expect(list.className).toContain("overflow-x-auto");
});

it("executes a dropdown shortcut once and closes the active menu", async () => {
  const run = vi.fn();
  disposers.push(
    render(
      () => (
        <DropdownMenu
          ariaLabel="Ações"
          trigger={(trigger) => <Button {...trigger}>Abrir menu</Button>}
          items={[
            { id: "edit", label: "Editar", shortcut: "Mod+E", onSelect: run },
          ]}
        />
      ),
      host,
    ),
  );
  await settle();
  host.querySelector<HTMLButtonElement>("button")!.click();
  await settle();
  press(document.activeElement!, "e", { ctrlKey: true });
  await settle();
  expect(run).toHaveBeenCalledTimes(1);
  expect(document.querySelector('[role="menu"]')).toBeNull();
  expect(document.activeElement).toBe(host.querySelector("button"));
});

it("supports global dropdown shortcuts and ignores disabled, closed and typing contexts", async () => {
  const run = vi.fn(),
    [disabled, setDisabled] = createSignal(false),
    [scope, setScope] = createSignal<"menu" | "global" | "none">("menu");
  disposers.push(
    render(
      () => (
        <>
          <DropdownMenu
            ariaLabel="Ações"
            disabled={disabled()}
            shortcutScope={scope()}
            trigger={(trigger) => <Button {...trigger}>Abrir</Button>}
            items={[
              { id: "edit", label: "Editar", shortcut: "Mod+E", onSelect: run },
              {
                id: "disabled",
                label: "Bloqueado",
                disabled: true,
                shortcut: "Mod+D",
                onSelect: run,
              },
            ]}
          />
          <input />
        </>
      ),
      host,
    ),
  );
  await settle();
  press(document.body, "e", { ctrlKey: true });
  expect(run).not.toHaveBeenCalled();
  setScope("global");
  await settle();
  press(document.body, "e", { metaKey: true });
  expect(run).toHaveBeenCalledTimes(1);
  press(document.body, "d", { ctrlKey: true });
  press(host.querySelector("input")!, "e", { ctrlKey: true });
  expect(run).toHaveBeenCalledTimes(1);
  setDisabled(true);
  await settle();
  press(document.body, "e", { ctrlKey: true });
  expect(run).toHaveBeenCalledTimes(1);
  setDisabled(false);
  setScope("none");
  await settle();
  press(document.body, "e", { ctrlKey: true });
  expect(run).toHaveBeenCalledTimes(1);
});

it("uses one typed definition for display and command execution", async () => {
  const shortcut = { key: "k", modifiers: ["Mod"] } satisfies Hotkey;
  const run = vi.fn();
  disposers.push(
    render(() => {
      createHotkeys(() => [{ hotkey: shortcut, callback: run }]);
      return <Kbd keys={shortcut} />;
    }, host),
  );
  await settle();
  expect(host.textContent).toContain("Ctrl/⌘");
  press(document.body, "k", { ctrlKey: true });
  press(document.body, "k", { metaKey: true });
  expect(run).toHaveBeenCalledTimes(2);
  // @ts-expect-error Unknown key names must fail at the public definition.
  const invalidKey: Hotkey = { key: "KeyK" };
  // @ts-expect-error Misspelled modifiers must fail at the public definition.
  const invalidModifier: Hotkey = { key: "k", modifiers: ["Mdo"] };
  // @ts-expect-error Unchecked string typos must fail too.
  const invalidString: HotkeyInput = "Mdo+K";
  void invalidKey;
  void invalidModifier;
  void invalidString;
});
it("validates dynamic shortcut strings and normalizes aliases", () => {
  expect(parseHotkey("control+shift+k")).toEqual({
    key: "k",
    modifiers: ["Ctrl", "Shift"],
  });
  expect(parseHotkey("Mod+Plus")).toEqual({ key: "+", modifiers: ["Mod"] });
  expect(parseHotkey("Mdo+K")).toBeNull();
  expect(parseHotkey("Mod+Mod+K")).toBeNull();
  expect(parseHotkey("Mod+Unknown")).toBeNull();
});
it("opens search from the standalone Kbd preview", async () => {
  disposers.push(
    render(
      () => <KbdExample values={{ keys: "Mod+K" }} update={() => {}} />,
      host,
    ),
  );
  await settle();
  press(document.body, "k", { ctrlKey: true });
  await settle();
  expect(document.querySelector('[role="combobox"]')).not.toBeNull();
  expect(document.querySelector('[role="listbox"]')?.textContent).toContain(
    "Meus eventos",
  );
  expect(document.querySelectorAll('[role="option"] svg')).toHaveLength(4);
  expect(document.querySelector('[role="listbox"]')?.textContent).toContain(
    "Gerenciar eventos e edições",
  );
  expect(document.querySelector('[role="listbox"]')?.textContent).toContain(
    "Ações",
  );
});
