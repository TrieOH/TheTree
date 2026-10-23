type Values = Record<string, string | boolean>;

/** Runnable Solid examples; values use JS literals so edited text stays valid JSX. */
export function usage(name: string, v: Values): string {
  const q = (key: string) => JSON.stringify(v[key]);
  const prop = (key: string) => `${key}={${q(key)}}`;
  const state = (initial: string, variable = "value") =>
    `const [${variable}, set${variable[0].toUpperCase() + variable.slice(1)}] = createSignal(${initial});`;
  let imports = name;
  let setup = "";
  let jsx = "";
  switch (name) {
    case "Loading":
      jsx = `<Loading label={${q("label")}} showLabel={${q("showLabel")}} size={${q("size")}} tone={${q("tone")}} layout={${q("layout")}} />`;
      break;
    case "Avatar":
      jsx = `<Avatar initials={${q("initials")}} name={${q("name")}} src={${q("src")}} size={${q("size")}} />`;
      break;
    case "Badge":
      jsx = `<Badge variant={${q("variant")}}>{${q("children")}}</Badge>`;
      break;
    case "Checkbox":
      setup = state(q("checked"), "checked");
      jsx = `<Checkbox label={${q("label")}} description={${q("description")}} disabled={${q("disabled")}} checked={checked()} onChange={setChecked} />`;
      break;
    case "SearchInput":
      setup = state(q("value"));
      jsx = `<SearchInput aria-label="Buscar projetos" placeholder={${q("placeholder")}} value={value()} onInput={event => setValue(event.currentTarget.value)} />`;
      break;
    case "Card":
      jsx = `<Card class="p-5">{${q("children")}}</Card>`;
      break;
    case "Link":
      jsx = `<Link href="#/">{${q("children")}}</Link>`;
      break;
    case "CodeBlock":
      imports += ", Toaster";
      jsx = `<Toaster position="bottom-right" />\n<CodeBlock code={${q("code")}} />`;
      break;
    case "Button":
      jsx = `<Button ${prop("variant")} ${prop("size")} ${prop("disabled")}\n  aria-label={${q("children")}}\n  onClick={() => console.log("Ação executada")}\n>\n  ${v.size === "icon" ? '<span aria-hidden="true">↗</span>' : `{${q("children")}}`}\n</Button>`;
      break;
    case "Input":
    case "Textarea":
      setup = state(q("value") ?? '""');
      jsx = `<${name} aria-label="${name === "Input" ? "E-mail" : "Descrição"}"\n  ${prop("placeholder")} ${prop("disabled")}\n  value={value()} onInput={event => setValue(event.currentTarget.value)}\n/>`;
      break;
    case "Label":
      imports += ", Input";
      jsx = `<Label for={${q("for")}}>{${q("children")}}</Label>\n<Input id={${q("for")}} />`;
      break;
    case "Field":
      imports += ", Input";
      jsx = `<Field ${prop("label")} ${prop("hint")} ${prop("error")}>\n  {ids => <Input id={ids.id} aria-describedby={ids.describedBy} />}\n</Field>`;
      break;
    case "Combobox":
      setup = state(q("value"));
      jsx = `<Combobox\n  options={[\n    { value: "sp", label: "São Paulo" },\n    { value: "rj", label: "Rio de Janeiro" },\n    { value: "bh", label: "Belo Horizonte" },\n  ]}\n  ${prop("placeholder")} value={value()} onChange={setValue}\n/>`;
      break;
    case "SectionTabs":
      setup = state(q("active"), "active");
      jsx = `<SectionTabs\n  items={[{ id: "overview", label: "Visão geral" }, { id: "activity", label: "Atividade", count: ${Number(v.count)} }, { id: "settings", label: "Ajustes", badge: ${q("badge")} }]}\n  active={active()} onChange={setActive}\n/>`;
      break;
    case "Accordion":
      jsx = `<Accordion ${prop("multiple")} items={[\n  { value: "start", title: ${q("title")}, content: "Escolha um componente e explore suas propriedades." },\n  { value: "more", title: "Posso personalizar?", content: "Sim. Ajuste as propriedades." },\n]} />`;
      break;
    case "Dialog":
    case "AlertModal":
    case "Drawer":
      imports += ", Button";
      setup = state("false", "open");
      jsx = `<Button onClick={() => setOpen(true)}>Abrir ${name}</Button>\n<${name} open={open()} ${prop("title")}\n  ${name === "Drawer" ? "onClose={() => setOpen(false)}" : "onOpenChange={setOpen}"}${name === "Dialog" ? ` ${prop("size")}` : ""}${name === "AlertModal" ? `\n  ${prop("variant")} description={${q("description")}}\n  onConfirm={() => { setOpen(false); }}` : ""}\n${name === "AlertModal" ? "/>" : ">\n  <p>Conteúdo do painel.</p>\n</" + name + ">"}`;
      break;
    case "MultiStepDialog":
      imports += ", Button, type MultiStepItem";
      setup = `const [open, setOpen] = createSignal(false);
const [values, setValues] = createSignal({ name: "", description: "" });
const steps: MultiStepItem<{ name: string; description: string }>[] = [
  { id: "details", title: "Detalhes", fields: [
    { name: "name", label: "Nome" },
    { name: "description", label: "Descrição" },
  ] },
  { id: "review", title: "Revisão", summary: {
    title: "Revise os dados",
    items: [
      { label: "Nome", value: () => values().name },
      { label: "Descrição", value: () => values().description },
    ],
  } },
];`;
      jsx = `<Button onClick={() => setOpen(true)}>Abrir formulário</Button>
<MultiStepDialog title={${q("title")}} open={open()} onOpenChange={setOpen}
  steps={steps} values={values()}
  onChange={(key, value) => setValues(previous => ({ ...previous, [key]: String(value) }))}
  onSubmit={() => { console.log(values()); return true; }}
/>`;
      break;
    case "Tooltip":
      imports += ", Button";
      jsx = `<Tooltip ${prop("label")}><Button variant="outline">Ajuda</Button></Tooltip>`;
      break;
    case "Toaster":
      imports += ", toast, Button";
      jsx = `<Toaster position={${q("position")}} />\n<Button onClick={() => toast.${v.type}(${q("message")}, { description: ${q("description")}, duration: ${Number(v.duration)}, position: ${q("position")}, dismissible: ${q("dismissible")} })}>Mostrar notificação</Button>`;
      break;
    case "EmptyState":
      if (v.showAction) imports += ", Button";
      jsx = `<EmptyState eyebrow={${q("eyebrow")}} ${v.showAction ? `action={<Button onClick={() => console.log("Ação executada")}>{${q("actionLabel")}}</Button>}` : ""} ${prop("title")} ${prop("description")} />`;
      break;
    case "Skeleton":
      jsx = `<div class="flex items-center gap-3">\n  <Skeleton class="size-11${v.rounded ? " rounded-full" : ""}" />\n  <div class="flex-1 space-y-2">\n    <Skeleton class="h-3 w-3/4" />\n    <Skeleton class="h-3 w-1/2" />\n  </div>\n</div>`;
      break;
    case "CardSkeleton":
      jsx = `<CardSkeleton class="w-full${v.compact ? " [&>div:first-child]:h-12" : ""}" />`;
      break;
    case "MultiStepStepper":
      setup = state(String(Number(v.currentStep)), "currentStep");
      jsx = `<MultiStepStepper\n  steps={[\n    { id: "details", title: "Detalhes", description: "Nome e informações do projeto." },\n    { id: "access", title: "Acesso", description: "Defina quem pode participar." },\n    { id: "review", title: "Revisão", description: "Confira os dados antes de criar." },\n  ]}\n  currentStep={currentStep()} ${v.interactive ? "onStepClick={setCurrentStep}" : ""}\n/>`;
      break;
    case "MultiStepFields":
      setup = 'const [values, setValues] = createSignal({ name: "" });';
      jsx = `<MultiStepFields\n  fields={[{ name: "name", label: ${q("label")}, disabled: ${q("disabled")}, required: ${q("required")}, kind: ${q("kind")}, placeholder: ${q("placeholder")}, hint: ${q("hint")}, layout: ${q("layout")} }]}\n  values={values()}\n  onChange={(key, value) => setValues(previous => ({ ...previous, [key]: value }))}\n/>`;
      break;
    case "MultiStepSummaryCard":
      jsx = `<MultiStepSummaryCard summary={{\n  title: ${q("title")}, badge: ${q("badge")}, items: [{ label: "Projeto", value: ${q("value")}, mono: ${q("mono")}, fullWidth: ${q("fullWidth")} }],\n}} />`;
      break;
    case "PaginatedContainer":
      setup = state(q("filterValue"), "search");
      imports += ", EmptyState";
      jsx = `<PaginatedContainer\n  items={[{ name: "Website" }, { name: "Aplicativo" }, { name: "Design system" }]}\n  pageSize={${Number(v.pageSize)}} layout={${q("layout")}} gap={${q("gap")}} filterPlaceholder={${q("filterPlaceholder")}} itemLabel={${q("itemLabel")}} minItemWidth={${q("minItemWidth")}}\n  filterValue={search()} onFilterChange={setSearch}\n  emptyState={<EmptyState title="Nenhum projeto encontrado" description="Tente outro termo de busca." />}\n  sortFields={[{ key: "name", label: "Nome" }]}\n  renderItems={items => <For each={items}>{item => <div>{item.name}</div>}</For>}\n/>`;
      break;
    case "OverlayScrollbar":
      jsx = v.enabled
        ? "<OverlayScrollbar />"
        : "<Show when={enabled()}>\n  <OverlayScrollbar />\n</Show>";
      jsx +=
        '\n<main>\n  <For each={Array.from({ length: 30 }, (_, i) => i + 1)}>\n    {n => <p class="p-8">Conteúdo {n}</p>}\n  </For>\n</main>';
      if (!v.enabled) setup = state("false", "enabled");
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
    CodeBlock: ["filename"],
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
      .map((key) =>
        ["rows", "minlength", "maxlength"].includes(key)
          ? `${key}={${Number(v[key])}}`
          : prop(key),
      )
      .join(" ");
    jsx =
      name === "Skeleton"
        ? jsx.replaceAll("<Skeleton", `<Skeleton ${attributes}`)
        : jsx.replace(`<${name}`, `<${name} ${attributes}`);
  }
  if (v.class) {
    const rootClass = new RegExp(`(<${name}\\b[^>]*?)class="([^"]*)"`);
    jsx = rootClass.test(jsx)
      ? jsx.replace(
          rootClass,
          (_, prefix: string, value: string) =>
            `${prefix}class={${JSON.stringify(`${value} ${v.class}`)}}`,
        )
      : jsx.replace(`<${name}`, `<${name} class={${q("class")}}`);
  }
  const solidImports = [
    setup ? "createSignal" : "",
    name === "PaginatedContainer" || name === "OverlayScrollbar" ? "For" : "",
    name === "OverlayScrollbar" && !v.enabled ? "Show" : "",
  ].filter(Boolean);
  return `${solidImports.length ? `import { ${solidImports.join(", ")} } from "solid-js";\n` : ""}import { ${imports} } from "@trieoh/ui-solid";\n\nexport function Example() {\n${setup ? `  ${setup}\n\n` : ""}  return (\n    <>\n${jsx
    .split("\n")
    .map((line) => `      ${line}`)
    .join("\n")}\n    </>\n  );\n}`;
}
