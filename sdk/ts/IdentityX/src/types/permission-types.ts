export type ObjSegment = {
  namespace: string;
  specifier: string;
};

export type ObjSuffix = "*" | "**" | null;

export interface PermissionObject {
  segments: ObjSegment[];
  suffix: ObjSuffix;
}

export interface PermissionDomain {
  object: PermissionObject;
  action: string; // "field:update" | "*"
}

export interface PermissionApi {
  object: string; // "project:1/schema:2/**"
  action: string; // "field:update"
}

type Digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9";

type Letter = 
  | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' 
  | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z'
  | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' 
  | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z';

type ValidCaracter = Letter | Digit | '_';

export type CannotStartWithNumber<T extends string> =
  T extends `${Digit}${string}` ? "The string cannot start with a number!" : T;

export type OnlyAlphanumeric<T extends string> = 
T extends "" 
    ? "The string cannot be empty" 
    : T extends `${infer Head}${infer Tail}`
      ? Head extends ValidCaracter
        ? (Tail extends "" ? T : OnlyAlphanumeric<Tail>) 
        : "Only a-z, A-Z, 0-9, and _ are permitted."
      : T;