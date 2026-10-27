import type { Values } from "../catalog";

export function usage(name: string, v: Values): string {
  let imports = name;
  let setup = "";
  let jsx = "";
  const prop = (key: string, customName?: string) => {
    const val = v[key];
    const attr = customName ?? key;
    if (val === undefined || val === "") return "";
    if (typeof val === "boolean") return val ? attr : "";
    return `${attr}="${val}"`;
  };
  const q = (key: string) => JSON.stringify(v[key]);
  const state = (initial: string, varName = "value") =>
    `const [${varName}, set${varName[0].toUpperCase()}${varName.slice(1)}] = createSignal(${initial});\\n\\n`;
  switch (name) {
    case "Sidebar":
      jsx = `<Sidebar currentId="overview" sections={[{ title: "Menu", items: [{ id: "overview", label: "Visão Geral" }] }]} />`;
      break;
    case "CatalogNav":
      jsx = `<CatalogNav\\n  activeId="button"\\n  top={{ id: "all", label: "Todos os componentes", href: "#/" }}\\n  sections={[{ id: "actions", title: "Ações", href: "#/group/acoes", items: [{ id: "button", label: "Button", href: "#/components/Button" }] }]}\\n/>`;
      break;
    case "Breadcrumb":
      jsx = `<Breadcrumb items={[{ label: "Início", href: "#/" }, { label: "Componentes" }, { label: "Breadcrumb", current: true }]} />`;
      break;
    case "Loading":
      jsx = `<Loading ${prop("size")} ${prop("tone")} ${prop("layout")} ${prop("showLabel")} label={${q("label")}} />`;
      break;
    case "Avatar":
      jsx = `<Avatar name={${q("name")}} initials={${q("initials")}} src={${q("src")}} ${prop("size")} />`;
      break;
    case "Badge":
      jsx = `<Badge ${prop("variant")}>${v.content}</Badge>`;
      break;
    case "Checkbox":
      setup = state(String(v.checked), "checked");
      jsx = `<Checkbox checked={checked()} onChange={setChecked} label={${q("label")}} description={${q("description")}} ${prop("indeterminate")} />`;
      break;
    case "SearchInput":
      setup = state(q("value") ?? '""');
      jsx = `<SearchInput value={value()} onInput={e => setValue(e.currentTarget.value)} ${prop("placeholder")} ${prop("loading")} />`;
      break;
    case "Card":
      jsx = `<Card ${prop("variant")}>\\n  <h3 class="text-base font-semibold">${v.title}</h3>\\n  <p class="text-sm text-muted-foreground">${v.content}</p>\\n</Card>`;
      break;
    case "Link":
      jsx = `<Link href="#/">{${q("children")}}</Link>`;
      break;
    case "CodeBlock":
      imports += ", Toaster";
      jsx = `<Toaster position="bottom-right" />\\n<CodeBlock code={${q("code")}} />`;
      break;
    case "Button":
      jsx = `<Button ${prop("variant")} ${prop("size")} ${prop("disabled")}\\n  aria-label={${q("children")}}\\n  onClick={() => console.log("Ação executada")}\\n>\\n  ${v.size === "icon" ? '<IconArrowUpRight class="size-4" />' : `{${q("children")}}`}\\n</Button>`;
      break;
    case "Input":
    case "Textarea":
      setup = state(q("value") ?? '""');
      jsx = `<${name} aria-label="${name === "Input" ? "E-mail" : "Descrição"}"\\n  ${prop("placeholder")} ${prop("disabled")}\\n  value={value()} onInput={event => setValue(event.currentTarget.value)}\\n/>`;
      break;
    case "Label":
      imports += ", Input";
      jsx = `<Label for={${q("for")}}>{${q("children")}}</Label>\\n<Input id={${q("for")}} />`;
      break;
    case "Field":
      imports += ", Input";
      jsx = `<Field label={${q("label")}} description={${q("description")}} error={${q("error")}} for="field-demo">\\n  <Input id="field-demo" placeholder="Digite aqui..." />\\n</Field>`;
      break;
    case "Combobox":
      setup = state(q("value"));
      jsx = `<Combobox value={value()} onChange={setValue} options={[{ value: "1", label: "Opção 1" }, { value: "2", label: "Opção 2" }]} ${prop("placeholder")} />`;
      break;
    case "SectionTabs":
      setup = state(q("value"));
      jsx = `<SectionTabs value={value()} onChange={setValue} items={[{ id: "tab1", label: "Geral" }, { id: "tab2", label: "Segurança" }]} />`;
      break;
    case "Accordion":
      jsx = `<Accordion items={[{ id: "1", title: "Primeiro item", content: "Conteúdo do item." }, { id: "2", title: "Segundo item", content: "Outro conteúdo." }]} />`;
      break;
    case "Dialog":
    case "AlertModal":
    case "Drawer":
    case "MultiStepDialog":
      setup = state("false", "open");
      imports += ", Button";
      jsx = `<Button onClick={() => setOpen(true)}>Abrir ${name}</Button>\\n\\n<${name} open={open()} onOpenChange={setOpen} title={${q("title")}}>\\n  <p>Conteúdo do painel.</p>\\n</${name}>`;
      break;
    case "Tooltip":
      imports += ", Button";
      jsx = `<Tooltip content={${q("content")}}>\\n  <Button variant="outline">Passe o mouse</Button>\\n</Tooltip>`;
      break;
    case "Toaster":
      imports = "toast, Button";
      jsx = `<Button onClick={() => toast.success("Ação realizada com sucesso!")}>Disparar Toast</Button>`;
      break;
    case "EmptyState":
      jsx = `<EmptyState title={${q("title")}} description={${q("description")}} />`;
      break;
    case "Skeleton":
      jsx = `<Skeleton ${prop("width")} ${prop("height")} />`;
      break;
    case "CardSkeleton":
      jsx = `<CardSkeleton lines={${Number(v.lines)}} hasImage={${Boolean(v.hasImage)}} />`;
      break;
    case "MultiStepStepper":
      setup = state(q("currentStep"), "currentStep");
      jsx = `<MultiStepStepper steps={[{ id: "1", title: "Conta", description: "Dados de acesso" }, { id: "2", title: "Perfil", description: "Informações públicas" }]} currentStep={currentStep()} onStepClick={setCurrentStep} />`;
      break;
    case "MultiStepFields":
      imports += ", createSignal";
      setup = `const [data, setData] = createSignal({ name: "", bio: "" });\\n\\n`;
      jsx = `<MultiStepFields values={data()} onChange={(key, val) => setData(prev => ({ ...prev, [key]: val }))} fields={[{ kind: "text", name: "name", label: "Nome" }, { kind: "textarea", name: "bio", label: "Biografia" }]} />`;
      break;
    case "MultiStepSummaryCard":
      jsx = `<MultiStepSummaryCard summary={{\\n  title: ${q("title")}, badge: ${q("badge")}, items: [{ label: "Projeto", value: ${q("value")}, mono: ${q("mono")}, fullWidth: ${q("fullWidth")} }],\\n}} />`;
      break;
    case "PaginatedContainer":
      setup = state(q("filterValue"), "search");
      imports += ", EmptyState";
      jsx = `<PaginatedContainer\\n  items={[{ name: "Website" }, { name: "Aplicativo" }, { name: "Design system" }]}\\n  pageSize={${Number(v.pageSize)}} layout={${q("layout")}} gap={${q("gap")}} filterPlaceholder={${q("filterPlaceholder")}} itemLabel={${q("itemLabel")}} minItemWidth={${q("minItemWidth")}}\\n  filterValue={search()} onFilterChange={setSearch}\\n  emptyState={<EmptyState title="Nenhum projeto encontrado" description="Tente outro termo de busca." />}\\n  sortFields={[{ key: "name", label: "Nome" }]}\\n  renderItems={items => <For each={items}>{item => <div>{item.name}</div>}</For>}\\n/>`;
      break;
    case "OverlayScrollbar":
      jsx = v.enabled
        ? "<OverlayScrollbar />"
        : "<Show when={enabled()}>\\n  <OverlayScrollbar />\\n</Show>";
      jsx +=
        '\\n<main>\\n  <For each={Array.from({ length: 30 }, (_, i) => i + 1)}>\\n    {n => <p class="p-8">Conteúdo {n}</p>}\\n  </For>\\n</main>';
      if (!v.enabled) setup = state("false", "enabled");
      break;
    case "Spinner":
      jsx = `<Spinner size={${q("size")}} activeColor="var(--primary)" trackColor="var(--accent)" />`;
      break;
    case "DesktopOnly":
      jsx = `<DesktopOnly\\n  threshold={${Number(v.threshold) || 1200}}\\n  title={${q("title")}}\\n  description={${q("description")}}\\n  actionLabel={${q("actionLabel")}}\\n>\\n  <div class="p-6 border border-border rounded-xl bg-card text-center">\\n    <p class="text-sm font-medium">Conteúdo visível apenas em telas desktop.</p>\\n  </div>\\n</DesktopOnly>`;
      break;
    case "CanvasControls":
      setup = state("1", "scale");
      jsx = `<div class="relative h-44 w-full rounded-xl border border-border bg-muted/20 flex items-center justify-center overflow-hidden">\\n  <div style={{ transform: \`scale(\${scale()})\` }} class="px-4 py-2 bg-card border border-border rounded-lg shadow-sm text-xs font-medium">\\n    Elemento no Canvas\\n  </div>\\n  <CanvasControls\\n    scale={scale()}\\n    onZoomIn={() => setScale(s => Math.min(s + 0.1, 2))}\\n    onZoomOut={() => setScale(s => Math.max(s - 0.1, 0.2))}\\n    onReset={() => setScale(1)}\\n  />\\n</div>`;
      break;
    case "WaveText":
      jsx = `<WaveText\\n  text={${q("text")}}\\n  duration={${Number(v.duration) || 1400}}\\n  delay={${Number(v.delay) || 500}}\\n  lift={${Number(v.lift) || 12}}\\n  waveWidth={${Number(v.waveWidth) || 20}}\\n/>`;
      break;
    case "WaveSpinnerLoading":
      jsx = `<WaveSpinnerLoading\\n  text={${q("text")}}\\n  spinnerSize={${q("spinnerSize")}}\\n/>`;
      break;
    default:
      throw new Error(`Missing example for ${name}`);
  }
  if (name === "AlertModal" && v.loading)
    jsx = jsx.replace(
      "<AlertModal",
      '<AlertModal children={<Button variant="outline" onClick={() => setOpen(false)}>Encerrar simulação</Button>}',
    );
  if (name === "Drawer")
    jsx = jsx.replace(
      "<p>Conteúdo do painel.</p>",
      `<p>{${q("description")}}</p>`,
    );
  const extraProps: Record<string, string[]> = {
    Checkbox: [
      "position",
      "alignment",
      "size",
      "variant",
      "controlShape",
      "indeterminate",
      "required",
      "name",
      "value",
    ],
    SearchInput: ["disabled", "readonly", "required", "name"],
    Button: ["type", "title", "name", "value"],
    Input: [
      "type",
      "required",
      "readonly",
      "name",
      "min",
      "max",
      "step",
      "minlength",
      "maxlength",
      "autocomplete",
      "inputmode",
    ],
    Textarea: ["required", "readonly", "rows", "name"],
    Field: ["required", "name"],
    Combobox: ["disabled", "searchPlaceholder", "emptyMessage"],
    SectionTabs: ["ariaLabel"],
    CodeBlock: ["filename", "showLineNumbers", "format"],
    Dialog: ["description", "closable", "preventClose"],
    AlertModal: ["loading", "confirmLabel", "cancelLabel"],
    MultiStepDialog: [
      "size",
      "description",
      "loading",
      "fixedHeight",
      "nextLabel",
      "submitLabel",
    ],
    Tooltip: ["placement", "align"],
    MultiStepStepper: ["orientation", "showDescriptions", "ariaLabel"],
    CatalogNav: ["mode"],
    DesktopOnly: ["threshold", "title", "description", "actionLabel"],
    CanvasControls: ["scale"],
    WaveText: ["text", "duration", "delay", "lift", "waveWidth"],
    WaveSpinnerLoading: ["text", "spinnerSize"],
    Spinner: ["size"],
  };
  if (
    [
      "Button",
      "Input",
      "Textarea",
      "Checkbox",
      "Combobox",
      "Field",
      "SearchInput",
    ].includes(name)
  )
    extraProps[name] = [...(extraProps[name] ?? []), "id"];
  const keys = [...(extraProps[name] ?? []), ...(v.shape ? ["shape"] : [])];
  if (keys.length) {
    const attributes = keys
      .filter(
        (key) =>
          v[key] !== undefined &&
          !(
            ["id", "min", "max", "step", "minlength", "maxlength"].includes(
              key,
            ) && v[key] === ""
          ),
      )
      .map((key) => prop(key))
      .filter(Boolean)
      .join(" ");
    if (attributes) {
      const isSelfClosing = jsx.includes("/>");
      const insertAt = isSelfClosing ? jsx.indexOf("/>") : jsx.indexOf(">");
      jsx = `${jsx.slice(0, insertAt)} ${attributes}${jsx.slice(insertAt)}`;
    }
  }
  const cleanSetup = setup ? `${setup.trim()}\\n\\n` : "";
  const code = `import { ${imports} } from "@trieoh/ui-solid";\\n\\n${cleanSetup}${jsx}`;
  return code.replace(/\\n/g, "\n");
}

export const getComponentUsage = usage;
