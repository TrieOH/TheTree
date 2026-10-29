import { cleanup, render, screen } from "@solidjs/testing-library";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { Reveal } from "@trieoh/ui-solid";

// jsdom has neither IntersectionObserver nor the Web Animations API, so both
// are replaced with recorders the tests can drive by hand.
class FakeObserver {
  static instances: FakeObserver[] = [];
  observed = new Set<Element>();

  constructor(
    private readonly callback: IntersectionObserverCallback,
    readonly options: IntersectionObserverInit,
  ) {
    FakeObserver.instances.push(this);
  }

  observe(element: Element) {
    this.observed.add(element);
  }

  unobserve(element: Element) {
    this.observed.delete(element);
  }

  disconnect() {
    this.observed.clear();
  }

  fire(element: Element, isIntersecting: boolean) {
    this.callback(
      [{ target: element, isIntersecting } as IntersectionObserverEntry],
      this as unknown as IntersectionObserver,
    );
  }
}

const animate = vi.fn();
const cancel = vi.fn();

// Reveal keeps its observers at module level on purpose (one per trigger
// configuration, shared by every instance), so observers outlive a single test
// and are found by what they currently watch.
const observerFor = (element: Element) =>
  FakeObserver.instances.find((instance) => instance.observed.has(element))!;

const watchedCount = () =>
  FakeObserver.instances.reduce((sum, instance) => sum + instance.observed.size, 0);

beforeEach(() => {
  animate.mockReset();
  cancel.mockReset();
  animate.mockImplementation(() => ({ cancel }));
  vi.stubGlobal("IntersectionObserver", FakeObserver);
  Object.defineProperty(Element.prototype, "animate", {
    configurable: true,
    writable: true,
    value: animate,
  });
  window.matchMedia = vi.fn().mockReturnValue({ matches: false });
});

afterEach(() => {
  cleanup();
  vi.unstubAllGlobals();
});

describe("Reveal", () => {
  it("stays hidden until it scrolls into view, then slides up once", () => {
    render(() => <Reveal>conteúdo</Reveal>);
    const element = screen.getByText("conteúdo");

    expect(element).toHaveAttribute("data-reveal", "hidden");
    expect(animate).not.toHaveBeenCalled();

    const observer = observerFor(element);
    observer.fire(element, true);

    expect(element).not.toHaveAttribute("data-reveal");
    expect(animate).toHaveBeenCalledTimes(1);
    const [frames, timing] = animate.mock.calls[0];
    expect(frames).toEqual([
      { opacity: 0, transform: "translate(0, 15px)" },
      { opacity: 1, transform: "translate(0, 0)" },
    ]);
    expect(timing).toMatchObject({ delay: 0, duration: 400, fill: "backwards" });
    // Played once: the element is no longer watched.
    expect(observer.observed.has(element)).toBe(false);
  });

  it("mounts visible with animate={false} and never observes", () => {
    render(() => <Reveal animate={false}>conteúdo</Reveal>);

    expect(screen.getByText("conteúdo")).not.toHaveAttribute("data-reveal");
    expect(watchedCount()).toBe(0);
  });

  it("mounts visible when the user prefers reduced motion", () => {
    window.matchMedia = vi.fn().mockReturnValue({ matches: true });
    render(() => <Reveal>conteúdo</Reveal>);

    expect(screen.getByText("conteúdo")).not.toHaveAttribute("data-reveal");
    expect(watchedCount()).toBe(0);
  });

  it.each([
    ["left", "translate(-15px, 0)"],
    ["right", "translate(15px, 0)"],
    ["down", "translate(0, -15px)"],
    ["up", "translate(0, 15px)"],
  ] as const)("slides in from the %s", (direction, from) => {
    render(() => <Reveal direction={direction}>conteúdo</Reveal>);
    const element = screen.getByText("conteúdo");
    observerFor(element).fire(element, true);

    expect(animate.mock.calls[0][0][0].transform).toBe(from);
  });

  it.each([
    ["fade", [{ opacity: 0 }, { opacity: 1 }]],
    [
      "scale",
      [
        { opacity: 0, transform: "scale(0.94)" },
        { opacity: 1, transform: "scale(1)" },
      ],
    ],
    [
      "blur",
      [
        { opacity: 0, filter: "blur(8px)" },
        { opacity: 1, filter: "blur(0)" },
      ],
    ],
  ] as const)("supports the %s variant", (variant, frames) => {
    render(() => <Reveal variant={variant}>conteúdo</Reveal>);
    const element = screen.getByText("conteúdo");
    observerFor(element).fire(element, true);

    expect(animate.mock.calls[0][0]).toEqual(frames);
  });

  it("applies distance, delay, duration and easing", () => {
    render(() => (
      <Reveal distance={40} delay={0.2} duration={1} easing="ease-in">
        conteúdo
      </Reveal>
    ));
    const element = screen.getByText("conteúdo");
    observerFor(element).fire(element, true);

    const [frames, timing] = animate.mock.calls[0];
    expect(frames[0].transform).toBe("translate(0, 40px)");
    expect(timing).toMatchObject({
      delay: 200,
      duration: 1000,
      easing: "ease-in",
    });
  });

  it("staggers direct children instead of the wrapper", () => {
    render(() => (
      <Reveal stagger={0.1} delay={0.05} data-testid="group">
        <p>a</p>
        <p>b</p>
        <p>c</p>
      </Reveal>
    ));
    const group = screen.getByText("a").parentElement!;
    observerFor(group).fire(group, true);

    expect(animate).toHaveBeenCalledTimes(3);
    expect(animate.mock.calls.map(([, timing]) => timing.delay)).toEqual([
      50, 150, 250,
    ]);
    expect(group.className).toContain("data-[reveal=hidden]:*:opacity-0");
  });

  it("replays on every return when once is false", () => {
    const onReveal = vi.fn();
    render(() => (
      <Reveal once={false} onReveal={onReveal}>
        conteúdo
      </Reveal>
    ));
    const element = screen.getByText("conteúdo");
    const observer = observerFor(element);

    observer.fire(element, true);
    expect(animate).toHaveBeenCalledTimes(1);
    expect(element).not.toHaveAttribute("data-reveal");

    observer.fire(element, false);
    expect(element).toHaveAttribute("data-reveal", "hidden");
    expect(cancel).toHaveBeenCalled();

    observer.fire(element, true);
    expect(animate).toHaveBeenCalledTimes(2);
    expect(onReveal).toHaveBeenCalledTimes(2);
  });

  it("shares one observer between instances with the same trigger", () => {
    render(() => (
      <>
        <Reveal>um</Reveal>
        <Reveal>dois</Reveal>
        <Reveal threshold={0.5}>três</Reveal>
      </>
    ));

    // Two Reveals share the default configuration, the third has its own.
    const shared = observerFor(screen.getByText("um"));
    expect(observerFor(screen.getByText("dois"))).toBe(shared);
    expect(observerFor(screen.getByText("três"))).not.toBe(shared);
    expect(shared.options.threshold).toBe(0.08);
  });

  it("stops watching when it unmounts before ever being seen", () => {
    const view = render(() => <Reveal>conteúdo</Reveal>);
    const element = screen.getByText("conteúdo");
    const observer = observerFor(element);
    expect(observer.observed.has(element)).toBe(true);

    view.unmount();

    expect(observer.observed.has(element)).toBe(false);
  });

  it("renders the requested element and merges the class", () => {
    render(() => (
      <Reveal as="section" class="h-full">
        conteúdo
      </Reveal>
    ));
    const element = screen.getByText("conteúdo");

    expect(element.tagName).toBe("SECTION");
    expect(element.className).toContain("h-full");
    expect(element.className).toContain("data-[reveal=hidden]:opacity-0");
  });
});
