import type { ObjSuffix, PermissionDomain, PermissionObject } from "../../types/permission-types";
import type { PermissionActionChain, PermissionFinal, PermissionObjectBuilder, PermissionObjectFinal, PermissionRoot } from "./permission.interfaces";
import { PermissionResult } from "./permission.result";
import { assertActionPart, assertNamespace, assertSpecifier } from "./permission.validators";

class Impl {
  segments: { namespace: string; specifier: string; actionParts: string[] }[] = [];
  suffix: ObjSuffix = null;
  wildcardOnly = false;

  commitSegment(ns: string, spec: string, actionParts: string[]) {
    this.segments.push({ namespace: ns, specifier: spec, actionParts });
  }

  buildStructured(): PermissionDomain {
    if (this.wildcardOnly) return { object: { segments: [], suffix: null }, action: "*" };

    if (this.segments.length === 0) throw new Error("No segments added. Use any() or add segments via on()");

    const object: PermissionObject = {
      segments: this.segments.map(s => ({ namespace: s.namespace, specifier: s.specifier })),
      suffix: this.suffix,
    };

    const last = this.segments[this.segments.length - 1];
    const action = last.actionParts.length === 0 ? "*" : last.actionParts.join(':');

    return { object, action };
  }
}

export function permission(): PermissionRoot {
  const s = new Impl();

  const root: PermissionRoot = {
    on(namespace, specifier) {
      const nsStr = namespace as unknown as string;
      const spStr = specifier as unknown as string;
      assertNamespace(nsStr);
      assertSpecifier(spStr);
      // start a fresh segment but do not commit actionParts yet
      s.commitSegment(nsStr, spStr, []); // actionParts empty until can(...) is called
      return objectBuilder();
    },
    onAll(namespace) {
      const nsStr = namespace as unknown as string;
      assertNamespace(nsStr);
      s.commitSegment(nsStr, "*", []); // specifier '*' but only via explicit API
      return objectBuilder();
    },
    any() {
      s.wildcardOnly = true;
      return final();
    }
  };

  function objectBuilder(): PermissionObjectBuilder {
    return {
      on(namespace, specifier) {
        const nsStr = namespace as unknown as string;
        const spStr = specifier as unknown as string;
        assertNamespace(nsStr);
        assertSpecifier(spStr);
        s.commitSegment(nsStr, spStr, []);
        return objectBuilder();
      },
      onAll(namespace) {
        const nsStr = namespace as unknown as string;
        assertNamespace(nsStr);
        s.commitSegment(nsStr, "*", []);
        return objectBuilder();
      },
      forAnyChild() {
        if (s.suffix) throw new Error("Suffix already set");
        s.suffix = '*';
        return objectFinal();
      },
      forAnyDescendant() {
        if (s.suffix) throw new Error("Suffix already set");
        s.suffix = '**';
        return objectFinal();
      },
      done() {
        return objectFinal();
      }
    };
  }

  function objectFinal(): PermissionObjectFinal {
    return {
      can(action) {
        const a = action as unknown as string;
        // validate first token
        const parts = a.split(':').filter(Boolean);
        if (parts.length === 0) throw new Error("Action must be non-empty");
        assertActionPart(parts[0]);
        // store into the last committed segment's actionParts
        const lastIdx = s.segments.length - 1;
        if (lastIdx < 0) throw new Error("No object segment to attach action to");
        // push initial parts
        s.segments[lastIdx].actionParts = parts.slice();
        return actionChainFor();
      },
      canAnyAction() {
        // set actionParts = ['*'] on last segment and finalize
        const lastIdx = s.segments.length - 1;
        if (lastIdx < 0) throw new Error("No object segment to attach action to");
        s.segments[lastIdx].actionParts = ['*'];
        return final();
      }
    };
  }

  function actionChainFor(): PermissionActionChain {
    return {
      and(part) {
        const p = part as unknown as string;
        assertActionPart(p);
        const lastIdx = s.segments.length - 1;
        s.segments[lastIdx].actionParts.push(p);
        return actionChainFor();
      },
      build() {
        const domain = s.buildStructured();
        return new PermissionResult(domain);
      }
    };
  }

  function final(): PermissionFinal {
    return {
      build() {
        const domain = s.buildStructured();
        return new PermissionResult(domain);
      }
    };
  }

  return root;
}