// ============================================================================
// TEST DATA - Mirrors Go testdata.go
// ============================================================================

export interface ValidationSpec {
    name: string;
    email: string;
    pass: string;
    errors: string[];
}

export interface PasswordSpec {
    name: string;
    password: string;
    errors: string[];
}

export const ValidPassword = "Str0ngP4$$!";

export const ValidationTests: ValidationSpec[] = [
    {
        name: "NoEmail",
        email: "",
        pass: ValidPassword,
        errors: ["email is required"],
    },
    {
        name: "InvalidEmail",
        email: "not-an-email",
        pass: ValidPassword,
        errors: ["valid email address"],
    },
    {
        name: "NoPassword",
        email: "test@mail.com",
        pass: "",
        errors: ["password is required"],
    },
];

export const WeakPasswordTests: PasswordSpec[] = [
    { name: "OnlyLetters",     password: "abc",   errors: ["uppercase", "number", "symbol"] },
    { name: "LettersNumber",   password: "abc3",  errors: ["uppercase", "symbol"] },
    { name: "LettersSymbol",   password: "abc#",  errors: ["uppercase", "number"] },
    { name: "LettersUppercase",password: "Abc",   errors: ["number", "symbol"] },
    { name: "NoNumber",        password: "Abc#",  errors: ["number"] },
    { name: "NoSymbol",        password: "Abc3",  errors: ["symbol"] },
    { name: "TooShort",        password: "Abc#3", errors: ["at least 8 characters"] },
];