import path from "node:path";
import { API, SignatureKind, SymbolFlags } from "typescript/unstable/sync";
import {
  isBinaryExpression,
  isPropertyAccessExpression,
  isStringLiteral,
  isNumericLiteral,
  SyntaxKind,
  type Node,
} from "typescript/unstable/ast";
import type { Plugin } from "vite";

export function apiReference(): Plugin {
  const moduleId = "virtual:api-reference";
  const resolvedId = `\0${moduleId}`;
  return {
    name: "showcase-api-reference",
    resolveId(id) {
      if (id === moduleId) return resolvedId;
    },
    load(id) {
      if (id !== resolvedId) return;
      const entry = path.resolve(
        import.meta.dirname,
        "../../../lib/ts/ui-solid/src/index.ts",
      );
      const api = new API();
      try {
        const snapshot = api.updateSnapshot({
          openProjects: [path.resolve(import.meta.dirname, "../tsconfig.json")],
        });
        const project = snapshot.getProjects()[0];
        const { program, checker } = project;
        const source = program.getSourceFile(entry)!;
        const reference: Record<
          string,
          {
            name: string;
            type: string;
            required: boolean;
            defaultValue?: string;
            description: string;
          }[]
        > = {};
        for (const file of program.getSourceFileNames()) {
          if (file.includes("/ui-solid/src/")) this.addWatchFile(file);
        }
        for (const symbol of checker.getExportsOfModule(
          checker.getSymbolAtLocation(source)!,
        )) {
          const component = checker.getTypeOfSymbolAtLocation(symbol, source);
          const signature = checker.getSignaturesOfType(
            component,
            SignatureKind.Call,
          )[0];
          const parameter = signature?.getParameters()[0];
          if (!parameter || !/^[A-Z]/.test(symbol.name)) continue;
          const defaults = new Map<string, string>();
          const collectDefaults = (node: Node) => {
            if (
              isBinaryExpression(node) &&
              node.operatorToken.kind === SyntaxKind.QuestionQuestionToken &&
              isPropertyAccessExpression(node.left) &&
              node.left.expression.getText() === parameter.name &&
              (isStringLiteral(node.right) ||
                isNumericLiteral(node.right) ||
                node.right.kind === SyntaxKind.TrueKeyword ||
                node.right.kind === SyntaxKind.FalseKeyword)
            ) {
              defaults.set(node.left.name.getText(), node.right.getText());
            }
            node.forEachChild(collectDefaults);
          };
          const implementation = signature.declaration?.resolve();
          if (implementation) collectDefaults(implementation);
          const props = checker.getTypeOfSymbolAtLocation(
            parameter,
            parameter.valueDeclaration!.resolve()!,
          );
          reference[symbol.name] = checker
            .getPropertiesOfType(props)
            .flatMap((prop) => {
              const declaration =
                prop.valueDeclaration ?? prop.declarations?.[0];
              if (!declaration || !declaration.path.includes("/ui-solid/src/"))
                return [];
              const type = checker.getTypeOfSymbolAtLocation(
                prop,
                declaration.resolve()!,
              );
              return [
                {
                  name: prop.name,
                  type: checker.typeToString(type, declaration.resolve()!),
                  required: !(prop.flags & SymbolFlags.Optional),
                  defaultValue: defaults.get(prop.name),
                  description: prop.getDocumentationComment(checker),
                },
              ];
            });
        }
        return `export default ${JSON.stringify(reference)}`;
      } finally {
        api.close();
      }
    },
    handleHotUpdate(context) {
      if (!context.file.includes("/ui-solid/src/")) return;
      const module = context.server.moduleGraph.getModuleById(resolvedId);
      if (module) {
        context.server.moduleGraph.invalidateModule(module);
        return [...context.modules, module];
      }
    },
  };
}
