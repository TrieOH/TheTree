const NAMESPACE_RE = /^[a-zA-Z][a-zA-Z0-9_]*$/;
const SPECIFIER_RE = /^([a-zA-Z0-9_.-]+|[0-9]+)$/;
const ACTION_PART_RE = /^[a-zA-Z0-9_]+$/;

export function assertNamespace(ns: string) {
  if (!NAMESPACE_RE.test(ns)) throw new Error(`Invalid namespace "${ns}"`);
  if (ns === "*" || ns === "**") throw new Error(`Namespace cannot be *`);
}

export function assertSpecifier(spec: string) {
  if (spec === "*" || spec === "**") throw new Error(`Specifier cannot be *`);
  if (!SPECIFIER_RE.test(spec)) throw new Error(`Invalid specifier "${spec}"`);
}

export function assertActionPart(part: string) {
  if (part === "*" || part === "**") throw new Error(`Action token cannot be *`);
  if (!ACTION_PART_RE.test(part)) throw new Error(`Invalid action part "${part}"`);
}


type Digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9";

type Letter = 
  | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' 
  | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z'
  | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' 
  | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z';

type ValidChar = Letter | Digit | '_';

type IsValidSequence<T extends string> = 
  T extends `${infer Head}${infer Tail}`
    ? Head extends ValidChar
      ? Tail extends "" ? true : IsValidSequence<Tail>
      : false
    : false;

export type ValidateNamespace<T extends string> = 
  T extends "" ? "Namespace cannot be empty" :
  T extends `${infer Head}${infer Tail}`
    ? Head extends Letter
      ? (Tail extends "" ? T : IsValidSequence<Tail> extends true ? T : "Namespace contains invalid characters")
      : "Namespace must start with a letter (a-zA-Z)"
    : T;

export type ValidateSpecifier<T extends string> = 
  T extends "" ? "Specifier cannot be empty" :
  IsValidSequence<T> extends true ? T : "Specifier contains invalid characters";


export type ValidateAction<T extends string> = 
  T extends "" 
    ? "Action cannot be empty" 
    : IsValidSequence<T> extends true 
      ? T 
      : "Action contains invalid characters (only a-z, A-Z, 0-9 and _)";
