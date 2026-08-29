const nextFrame = (view: Window) =>
  new Promise<void>((resolve) => view.requestAnimationFrame(() => resolve()));

export async function printElement(element: HTMLElement, title: string) {
  const iframe = document.createElement("iframe");
  iframe.setAttribute("aria-hidden", "true");
  Object.assign(iframe.style, {
    position: "fixed",
    width: "0",
    height: "0",
    border: "0",
    right: "0",
    bottom: "0",
  });
  document.body.appendChild(iframe);

  const printDocument = iframe.contentDocument;
  const printWindow = iframe.contentWindow;
  if (!printDocument || !printWindow) {
    iframe.remove();
    return;
  }

  const styles = [...document.querySelectorAll("link[rel='stylesheet'], style")]
    .map((node) => node.outerHTML)
    .join("");
  printDocument.open();
  printDocument.write(`<!doctype html>
    <html class="${document.documentElement.className}">
      <head><base href="${document.baseURI}"><title>${title}</title>${styles}<style>@page { margin: 2mm !important; } body { margin: 0 !important; }</style></head>
      <body>${element.outerHTML}</body>
    </html>`);
  printDocument.close();

  await Promise.all(
    [
      ...printDocument.querySelectorAll<HTMLLinkElement>(
        "link[rel='stylesheet']",
      ),
    ]
      .filter((link) => !link.sheet)
      .map(
        (link) =>
          new Promise<void>((resolve) => {
            link.addEventListener("load", () => resolve(), { once: true });
            link.addEventListener("error", () => resolve(), { once: true });
          }),
      ),
  );
  await printDocument.fonts?.ready;
  await Promise.all(
    [...printDocument.images].map((image) =>
      image.decode().catch(() => undefined),
    ),
  );
  await nextFrame(printWindow);
  await nextFrame(printWindow);

  printWindow.addEventListener("afterprint", () => iframe.remove(), {
    once: true,
  });
  printWindow.focus();
  printWindow.print();
}
