import type { ObjSuffix, PermissionDomain, PermissionObject } from "../../types/permission-types";
import type { PermissionActionBuilder, PermissionActionChain, PermissionFinal, PermissionObjectBuilder, PermissionRoot } from "./permission.interfaces";
import { PermissionResult } from "./permission.result";
import { assertActionPart, assertNamespace, assertSpecifier } from "./permission.validators";

class Impl {
  segments: { namespace: string; specifier: string; actionParts: string[] }[] = [];
  suffix: ObjSuffix = null;
  wildcardOnly = false;

  clone(): Impl {
    const newImpl = new Impl();
    newImpl.segments = structuredClone(this.segments);
    newImpl.suffix = this.suffix;
    newImpl.wildcardOnly = this.wildcardOnly;
    return newImpl;
  }

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
  function createWorkflow(s: Impl): PermissionRoot {
    const root: PermissionRoot = {
      on(namespace, specifier) {
        const nsStr = namespace as unknown as string;
        const spStr = specifier as unknown as string;
        assertNamespace(nsStr);
        assertSpecifier(spStr);

        const next = s.clone();
        // start a fresh segment but do not commit actionParts yet
        next.commitSegment(nsStr, spStr, []); // actionParts empty until can(...) is called
        return objectBuilder(next);
      },
      onAll(namespace) {
        const nsStr = namespace as unknown as string;
        assertNamespace(nsStr);

        const next = s.clone();
        next.commitSegment(nsStr, "*", []); // specifier '*' but only via explicit API
        return objectBuilder(next);
      },
      any() {
        const next = s.clone();
        next.wildcardOnly = true;
        return final(next);
      }
    };

    function objectBuilder(currentS: Impl): PermissionObjectBuilder {
      return {
        on(namespace, specifier) {
          const nsStr = namespace as unknown as string;
          const spStr = specifier as unknown as string;
          assertNamespace(nsStr);
          assertSpecifier(spStr);

          const next = currentS.clone();
          next.commitSegment(nsStr, spStr, []);
          return objectBuilder(next);
        },
        onAll(namespace) {
          const nsStr = namespace as unknown as string;
          assertNamespace(nsStr);

          const next = currentS.clone();
          next.commitSegment(nsStr, "*", []);
          return objectBuilder(next);
        },
        forAnyChild() {
          const next = currentS.clone();
          if (next.suffix) throw new Error("Suffix already set");
          next.suffix = '*';
          return actionBuilder(next);
        },
        forAnyDescendant() {
          const next = currentS.clone();
          if (next.suffix) throw new Error("Suffix already set");
          next.suffix = '**';
          return actionBuilder(next);
        },
        
        in(action) {
          const a = action as unknown as string;
          // validate first token
          const parts = a.split(':').filter(Boolean);
          if (parts.length === 0) throw new Error("Action must be non-empty");
          assertActionPart(parts[0]);

          // store into the last committed segment's actionParts
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          if (lastIdx < 0) throw new Error("No object segment to attach action to");
          // push initial parts
          next.segments[lastIdx].actionParts = parts.slice();
          return actionChainFor(next);
        },
        inAny() {
          // set actionParts = ['*'] on last segment and finalize
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          if (lastIdx < 0) throw new Error("No object segment to attach action to");
          next.segments[lastIdx].actionParts = ['*'];
          return final(next);
        },
        can(action) {
          const a = action as unknown as string;
          // validate first token
          const parts = a.split(':').filter(Boolean);
          if (parts.length === 0) throw new Error("Action must be non-empty");
          assertActionPart(parts[0]);

          // store into the last committed segment's actionParts
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          if (lastIdx < 0) throw new Error("No object segment to attach action to");
          // push initial parts
          next.segments[lastIdx].actionParts = parts.slice();
          return final(next);
        },
      };
    }

    function actionBuilder(currentS: Impl): PermissionActionBuilder {
      return {
        in(action) {
          const a = action as unknown as string;
          // validate first token
          const parts = a.split(':').filter(Boolean);
          if (parts.length === 0) throw new Error("Action must be non-empty");
          assertActionPart(parts[0]);

          // store into the last committed segment's actionParts
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          if (lastIdx < 0) throw new Error("No object segment to attach action to");
          // push initial parts
          next.segments[lastIdx].actionParts = parts.slice();
          return actionChainFor(next);
        },
        inAny() {
          // set actionParts = ['*'] on last segment and finalize
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          if (lastIdx < 0) throw new Error("No object segment to attach action to");
          next.segments[lastIdx].actionParts = ['*'];
          return final(next);
        },
        can(action) {
          const a = action as unknown as string;
          // validate first token
          const parts = a.split(':').filter(Boolean);
          if (parts.length === 0) throw new Error("Action must be non-empty");
          assertActionPart(parts[0]);

          // store into the last committed segment's actionParts
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          if (lastIdx < 0) throw new Error("No object segment to attach action to");
          // push initial parts
          next.segments[lastIdx].actionParts = parts.slice();
          return final(next);
        },
      };
    }

    function actionChainFor(currentS: Impl): PermissionActionChain {
      return {
        in(part) {
          const p = part as unknown as string;
          assertActionPart(p);
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          next.segments[lastIdx].actionParts.push(p);
          return actionChainFor(next);
        },
        inAny() {
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          next.segments[lastIdx].actionParts.push('*');
          return actionChainFor(next);
        },
        can(action) {
          const a = action as unknown as string;
          // validate first token
          const parts = a.split(':').filter(Boolean);
          if (parts.length === 0) throw new Error("Action must be non-empty");
          assertActionPart(parts[0]);

          // store into the last committed segment's actionParts
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          if (lastIdx < 0) throw new Error("No object segment to attach action to");
          // push initial parts
          next.segments[lastIdx].actionParts = parts.slice();
          return final(next);
        },
        forAnyChild() {
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          next.segments[lastIdx].actionParts.push('*');
          return final(next);
        },
        forAnyDescendant() {
          const next = currentS.clone();
          const lastIdx = next.segments.length - 1;
          next.segments[lastIdx].actionParts.push('**');
          return final(next);
        },
      };
    }

    function final(currentS: Impl): PermissionFinal {
      return {
        finish() {
          return new PermissionResult(currentS.buildStructured());
        }
      };
    }
    return root;
  }
  return createWorkflow(new Impl());
}