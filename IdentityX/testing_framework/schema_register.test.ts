import { beforeAll, describe, expect, test } from "vitest";
import { createClient } from "./helpers/index.js";
import {
    assertErrID,
    assertMessage,
    shouldFail,
} from "./helpers/assert.js";
import {
    Validate,
    AnyUUID,
    AnyDate,
    AnyNumber,
    AnyString,
    AsString,
    Store,
} from "./helpers/validate.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

// Error IDs
const ErrRequestMissingSchemaCustomFields = "0_REQ_0000_S";
const ErrFieldValidationOnSchemaRegister  = "0_FIELD_0000_D";
const ErrProjectUserRegisterNoVersion     = "0_PROJECTUSER_0005_S";
const ErrProjectUserRegisterOnSchemaDraft = "0_PROJECTUSER_0001_S";
const ErrAuthEmailAlreadyUsed             = "0_AUTH_0000_D";

function assertTrace(body: any, ...expected: string[]) {
    const trace: string[] = body?.trace ?? [];
    for (const exp of expected) {
        const found = trace.some((entry: string) => entry.includes(exp));
        expect(found, `missing trace entry: expected trace to contain "${exp}"\ntrace=${JSON.stringify(trace)}`).toBe(true);
    }
}

describe("SchemaRegister", () => {
    let user: any;
    let projectID: string;
    let schemaID: string;
    let schemaVersion1ID: string;

    // Unauthenticated client for project register calls
    let anonClient: any;

    beforeAll(async () => {
        anonClient = createClient();

        await createClient().withCredentials("schemas_register@mail.com", ValidPassword).register();
        user = await createClient().withCredentials("schemas_register@mail.com", ValidPassword).login();
        const project = await user.post("/projects", { project_name: "schema testing", metadata: { env: "test" } });
        projectID = project.id;

        // Draft schema
        const schema = await user.post(`/projects/${projectID}/schemas`, {
            schema_type: "context",
            title: "scti",
            flow_id: "estudante",
        });
        const schemaRef = { current: "" };
        Validate(schema, {
            id:                 Store(schemaRef, AnyUUID),
            project_id:         AsString(projectID, AnyUUID),
            title:              "scti",
            flow_id:            "estudante",
            type:               "context",
            status:             "draft",
            current_version_id: null,
            created_at:         AnyDate,
            updated_at:         AnyDate,
        });
        schemaID = schemaRef.current as string;

        // Draft version
        const version = await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/draft`);
        const versionRef = { current: "" };
        Validate(version, {
            id:             Store(versionRef, AnyUUID),
            schema_id:      AsString(schemaID, AnyUUID),
            version_number: 1,
        });
        schemaVersion1ID = versionRef.current as string;

        // Verify schema state
        const schemaState = await user.get(`/projects/${projectID}/schemas/${schemaID}`);
        Validate(schemaState, {
            id:                 AsString(schemaID, AnyUUID),
            project_id:         AsString(projectID, AnyUUID),
            title:              "scti",
            flow_id:            "estudante",
            type:               "context",
            status:             "draft",
            current_version_id: AsString(schemaVersion1ID, AnyUUID),
        });

        // Create fields
        const fields = await user.post(`/projects/${projectID}/schemas/${schemaID}/v1`, {
            fields: [
                { key: "matricula", type: "string", owner: "user", title: "Numero da Matrícula", description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045",          required: true,  mutable: true, position: 0 },
                { key: "curso",     type: "string", owner: "user", title: "Curso de Matrícula",  description: "O curso que você está matrículado na UENF",              placeholder: "Ciência da Computação", required: true,  mutable: true, position: 1 },
                { key: "periodo",   type: "int",    owner: "user", title: "Período Atual",       description: "O período da sua matéria mais avançada da grade",                                              required: true,  mutable: true, position: 2 },
                { key: "ativo",     type: "bool",   owner: "user", title: "Ativo",               description: "Se o aluno está ativo",                                                                        required: false, mutable: true, position: 3 },
            ],
        });
        Validate(fields, [
            { object_id: AnyUUID, id: AnyUUID },
            { object_id: AnyUUID, id: AnyUUID },
            { object_id: AnyUUID, id: AnyUUID },
            { object_id: AnyUUID, id: AnyUUID },
        ]);

        // Publish
        await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/publish`);
        await user.post(`/projects/${projectID}/schemas/${schemaID}/publish`);
    });

    // Helper: POST to project register with query params
    async function projectRegister(params: { flow_id: string; schema_type: string }, body: object) {
        return anonClient.http.post(
            `/projects/${projectID}/register?schema_type=${params.schema_type}&flow_id=${params.flow_id}`,
            body
        );
    }

    test("RegisterOnSchemaNoCustomFields", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "client@email.com",
                password: ValidPassword,
            }),
            400
        );
        assertErrID(body, ErrRequestMissingSchemaCustomFields);
        assertMessage(body, "schema custom fields are required on a schema register");
    });

    test("RegisterOnSchemaEmptyCustomFields", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "client@email.com",
                password: ValidPassword,
                custom_fields: {},
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body,
            "form missing required field: matricula",
            "form missing required field: curso",
            "form missing required field: periodo",
        );
    });

    test("RegisterOnSchemaNoCursoField", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "client@email.com",
                password: ValidPassword,
                custom_fields: { matricula: "20221100033" },
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body,
            "form missing required field: curso",
            "form missing required field: periodo",
        );
    });

    test("RegisterOnSchemaNoMatriculaField", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "client@email.com",
                password: ValidPassword,
                custom_fields: { curso: "Ciência da Computação" },
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body,
            "form missing required field: matricula",
            "form missing required field: periodo",
        );
    });

    test("RegisterOnSchemaUnknownField", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "client@email.com",
                password: ValidPassword,
                custom_fields: { valor: "4", bing: "bong" },
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body,
            "form missing required field: matricula",
            "form missing required field: curso",
            "form missing required field: periodo",
        );
    });

    test("RegisterOnSchemaWrongTypeStringOnInt", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "client@email.com",
                password: ValidPassword,
                custom_fields: { periodo: "abc" },
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body,
            "form missing required field: matricula",
            "form missing required field: curso",
            "invalid form value for periodo: type(int) value(abc)",
        );
    });

    test("RegisterOnSchemaWrongTypeIntOnString", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "client@email.com",
                password: ValidPassword,
                custom_fields: { matricula: 20221100033 },
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body,
            "invalid form value for matricula: type(string) value(2.0221100033e+10)",
            "form missing required field: curso",
            "form missing required field: periodo",
        );
    });

    test("RegisterOnSchemaTypeFloatOnInt", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "float_on_int@email.com",
                password: ValidPassword,
                custom_fields: { matricula: "20221100033", curso: "Ciência da Computação", periodo: 4.5 },
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body, "invalid form value for periodo: type(int) value(4.5)");
    });

    test("RegisterOnSchemaTypeStringOnBool", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "string_on_bool@email.com",
                password: ValidPassword,
                custom_fields: { matricula: "20221100033", curso: "Ciência da Computação", periodo: 4, ativo: "true" },
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body, "invalid form value for ativo: type(bool) value(true)");
    });

    test("RegisterOnSchemaTypeIntOnBool", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "int_on_bool@email.com",
                password: ValidPassword,
                custom_fields: { matricula: "20221100033", curso: "Ciência da Computação", periodo: 4, ativo: 1 },
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body, "invalid form value for ativo: type(bool) value(1)");
    });

    test("RegisterOnSchemaTypeBoolOnString", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "bool_on_string@email.com",
                password: ValidPassword,
                custom_fields: { matricula: true },
            }),
            400
        );
        assertErrID(body, ErrFieldValidationOnSchemaRegister);
        assertMessage(body, "error validating form for schema register");
        assertTrace(body,
            "invalid form value for matricula: type(string) value(true)",
            "form missing required field: curso",
            "form missing required field: periodo",
        );
    });

    test("RegisterOnSchemaTypeFloatZeroOnInt", async () => {
        // 4.0 is a valid integer representation in JSON — should succeed
        const res = await projectRegister({ schema_type: "context", flow_id: "estudante" }, {
            email: "float_zero@email.com",
            password: ValidPassword,
            custom_fields: { matricula: "20221100033", curso: "Ciência da Computação", periodo: 4.0, ativo: true },
        });
        expect(res.status).toBe(201);
    });

    test("RegisterOnSchemaSuccess", async () => {
        const res = await projectRegister({ schema_type: "context", flow_id: "estudante" }, {
            email: "client@email.com",
            password: ValidPassword,
            custom_fields: { matricula: "20221100033", curso: "Ciência da Computação", periodo: 4 },
        });
        expect(res.status).toBe(201);
    });

    test("RegisterOnSchemaDuplicateEmail", async () => {
        const body = await shouldFail(
            projectRegister({ schema_type: "context", flow_id: "estudante" }, {
                email: "client@email.com",
                password: ValidPassword,
                custom_fields: { matricula: "20221100033", curso: "Ciência da Computação", periodo: 4 },
            }),
            409
        );
        assertErrID(body, ErrAuthEmailAlreadyUsed);
        assertMessage(body, "email already in use");
    });

    test("SchemaUserSessionInfo", async () => {
        const projClient = await createClient()
            .withCredentials("client@email.com", ValidPassword)
            .projectLogin(projectID);

        const data = await projClient.get("/sessions/me");

        Validate(data, {
            refresh_expire_date: AnyNumber,
            access: {
                iss: AsString(projectID, AnyUUID),
                exp: AnyNumber,
                iat: AnyNumber,
                jti: AnyUUID,
                sub: {
                    id:         AnyUUID,
                    email:      "client@email.com",
                    project_id: projectID,
                    user_type:  "project",
                    session_id: AnyUUID,
                    user_agent: AnyString,
                    user_ip:    AnyString,
                    metadata: {
                        context: {
                            estudante: {
                                schema_id:         schemaID,
                                schema_version_id: schemaVersion1ID,
                                fields: {
                                    curso:     "Ciência da Computação",
                                    matricula: "20221100033",
                                    periodo:   4,
                                },
                            },
                        },
                    },
                },
            },
        });
    });

    // ============================================================
    // SchemaStateEdgeCases
    // ============================================================

    describe("SchemaStateEdgeCases", () => {
        let edgeUser: any;
        let edgeProjectID: string;
        let edgeSchemaID: string;
        let edgeClient: any;

        beforeAll(async () => {
            edgeClient = createClient();
            await createClient().withCredentials("schema_state@mail.com", ValidPassword).register();
            edgeUser = await createClient().withCredentials("schema_state@mail.com", ValidPassword).login();
            const project = await edgeUser.post("/projects", { project_name: "Schema State Project", metadata: { env: "test" } });
            edgeProjectID = project.id;

            const schema = await edgeUser.post(`/projects/${edgeProjectID}/schemas`, {
                schema_type: "context",
                title: "No Version Schema",
                flow_id: "noversion",
            });
            edgeSchemaID = schema.id;
        });

        async function edgeRegister(body: object) {
            return edgeClient.http.post(
                `/projects/${edgeProjectID}/register?schema_type=context&flow_id=noversion`,
                body
            );
        }

        test("RegisterFailsNoVersion", async () => {
            const body = await shouldFail(
                edgeRegister({ email: "user@noversion.com", password: ValidPassword, custom_fields: {} }),
                400
            );
            assertErrID(body, ErrProjectUserRegisterNoVersion);
            assertMessage(body, "can't register on a schema that has no published version");
        });

        test("RegisterFailsVersionNotPublished", async () => {
            // Draft version + add field (still not published)
            await edgeUser.post(`/projects/${edgeProjectID}/schemas/${edgeSchemaID}/versions/draft`);
            await edgeUser.post(`/projects/${edgeProjectID}/schemas/${edgeSchemaID}/v1`, {
                fields: [{ key: "test", type: "string", owner: "user", title: "Test", position: 0, required: true }],
            });

            const body = await shouldFail(
                edgeRegister({ email: "user@notpublished.com", password: ValidPassword, custom_fields: { test: "val" } }),
                400
            );
            assertErrID(body, ErrProjectUserRegisterOnSchemaDraft);
            assertMessage(body, "can't register to a draft schema");
        });

        test("RegisterFailsSchemaNotPublished", async () => {
            // Publish version but NOT the schema
            await edgeUser.post(`/projects/${edgeProjectID}/schemas/${edgeSchemaID}/versions/publish`);

            const body = await shouldFail(
                edgeRegister({ email: "user@schemadraft.com", password: ValidPassword, custom_fields: { test: "val" } }),
                400
            );
            assertErrID(body, ErrProjectUserRegisterOnSchemaDraft);
            assertMessage(body, "can't register to a draft schema");
        });
    });

    // ============================================================
    // UnimplementedFieldTypes
    // ============================================================

    test("UnimplementedFieldTypes", async () => {
        await createClient().withCredentials("unimplemented@mail.com", ValidPassword).register();
        const userU = await createClient().withCredentials("unimplemented@mail.com", ValidPassword).login();
        const project = await userU.post("/projects", { project_name: "Unimplemented Project", metadata: { env: "test" } });
        const pid = project.id;

        const schema = await userU.post(`/projects/${pid}/schemas`, { schema_type: "context", title: "Email Schema", flow_id: "emailsync" });
        const sid = schema.id;

        await userU.post(`/projects/${pid}/schemas/${sid}/versions/draft`);
        await userU.post(`/projects/${pid}/schemas/${sid}/v1`, {
            fields: [{ key: "contact", type: "email", owner: "user", title: "Contact Email", position: 0, required: true }],
        });
        await userU.post(`/projects/${pid}/schemas/${sid}/versions/publish`);
        await userU.post(`/projects/${pid}/schemas/${sid}/publish`);
    });

    // ============================================================
    // Amplified Tests
    // ============================================================

    describe("Amplified", () => {
        let ampUser: any;
        let ampProjectID: string;
        let ampSchemaID: string;
        let ampClient: any;

        beforeAll(async () => {
            ampClient = createClient();
            await createClient().withCredentials("amplified_schemas@mail.com", ValidPassword).register();
            ampUser = await createClient().withCredentials("amplified_schemas@mail.com", ValidPassword).login();
            const project = await ampUser.post("/projects", { project_name: "Amplified Schema Testing", metadata: { env: "test" } });
            ampProjectID = project.id;

            // Draft schema + version
            const schema = await ampUser.post(`/projects/${ampProjectID}/schemas`, { schema_type: "context", title: "Registration with Options", flow_id: "registration-v2" });
            ampSchemaID = schema.id;
            await ampUser.post(`/projects/${ampProjectID}/schemas/${ampSchemaID}/versions/draft`);

            // Create fields with options and rules
            await ampUser.post(`/projects/${ampProjectID}/schemas/${ampSchemaID}/v1`, {
                fields: [
                    { key: "name",            type: "string", owner: "user", title: "Nome Completo",   description: "Seu nome completo", required: true,  mutable: true, position: 0 },
                    { key: "email",           type: "email",  owner: "user", title: "Email",                                             required: true,  mutable: true, position: 1 },
                    { key: "user_type",       type: "select", owner: "user", title: "Tipo de Usuário",                                   required: true,  mutable: true, position: 2,
                        options: [
                            { value: "student", label: "Estudante",   position: 0 },
                            { value: "teacher", label: "Professor",   position: 1 },
                            { value: "staff",   label: "Funcionário", position: 2 },
                        ],
                    },
                    { key: "is_active",       type: "bool",   owner: "user", title: "Está Ativo?",                                       required: false, mutable: true, position: 3, default_value: true },
                    { key: "student_id",      type: "string", owner: "user", title: "RA do Aluno",                                       required: false, mutable: true, position: 4,
                        visibility_rules: [{ depends_on_field_key: "user_type", operator: "equals", value: "student" }],
                    },
                    { key: "activation_date", type: "string", owner: "user", title: "Data de Ativação",                                  required: false, mutable: true, position: 5,
                        required_rules: [{ depends_on_field_key: "is_active", operator: "equals", value: true }],
                    },
                    { key: "shift",           type: "radio",  owner: "user", title: "Turno",                                             required: false, mutable: true, position: 6,
                        options: [
                            { value: "morning",   label: "Matutino",   position: 0 },
                            { value: "afternoon", label: "Vespertino", position: 1 },
                            { value: "night",     label: "Noturno",    position: 2 },
                        ],
                    },
                ],
            });

            await ampUser.post(`/projects/${ampProjectID}/schemas/${ampSchemaID}/versions/publish`);
            await ampUser.post(`/projects/${ampProjectID}/schemas/${ampSchemaID}/publish`);
        });

        async function ampRegister(body: object) {
            return ampClient.http.post(
                `/projects/${ampProjectID}/register?schema_type=context&flow_id=registration-v2`,
                body
            );
        }

        test("InvalidOptionValue", async () => {
            const body = await shouldFail(
                ampRegister({ email: "invalid_option@test.com", password: ValidPassword, custom_fields: { name: "John Doe", email: "john@test.com", user_type: "invalid_type" } }),
                400
            );
            assertErrID(body, ErrFieldValidationOnSchemaRegister);
            assertTrace(body, "invalid form value for user_type: type(select) value(invalid_type)");
        });

        test("CaseSensitiveOption", async () => {
            const body = await shouldFail(
                ampRegister({ email: "case_sensitive@test.com", password: ValidPassword, custom_fields: { name: "John Doe", email: "john@test.com", user_type: "STUDENT" } }),
                400
            );
            assertErrID(body, ErrFieldValidationOnSchemaRegister);
            assertTrace(body, "invalid form value for user_type: type(select) value(STUDENT)");
        });

        test("ValidOptionSuccess", async () => {
            const res = await ampRegister({ email: "valid_option@test.com", password: ValidPassword, custom_fields: { name: "John Doe", email: "john@test.com", user_type: "student" } });
            expect(res.status).toBe(201);
        });

        test("VisibilityRuleHiddenFieldMissing", async () => {
            // user_type=teacher → student_id hidden → missing is OK
            const res = await ampRegister({ email: "hidden_field@test.com", password: ValidPassword, custom_fields: { name: "Jane Doe", email: "jane@test.com", user_type: "teacher" } });
            expect(res.status).toBe(201);
        });

        test("VisibilityRuleVisibleButOptional", async () => {
            // user_type=student → student_id visible but required=false
            const res = await ampRegister({ email: "visible_optional@test.com", password: ValidPassword, custom_fields: { name: "Student Optional", email: "optional@test.com", user_type: "student" } });
            expect(res.status).toBe(201);
        });

        test("RequiredRuleTriggered", async () => {
            const body = await shouldFail(
                ampRegister({ email: "req_rule_triggered@test.com", password: ValidPassword, custom_fields: { name: "Active No Date", email: "activenodate@test.com", user_type: "staff", is_active: true } }),
                400
            );
            assertErrID(body, ErrFieldValidationOnSchemaRegister);
            assertTrace(body, "form missing required field: activation_date");
        });

        test("RequiredRuleSatisfied", async () => {
            const res = await ampRegister({ email: "req_rule_satisfied@test.com", password: ValidPassword, custom_fields: { name: "Active With Date", email: "activewithdate@test.com", user_type: "staff", is_active: true, activation_date: "2024-01-15" } });
            expect(res.status).toBe(201);
        });

        test("RequiredRuleNotTriggered", async () => {
            const res = await ampRegister({ email: "req_rule_not_triggered@test.com", password: ValidPassword, custom_fields: { name: "Inactive User", email: "inactive@test.com", user_type: "staff", is_active: false } });
            expect(res.status).toBe(201);
        });

        test("RequiredRuleBothMissing", async () => {
            const res = await ampRegister({ email: "both_missing@test.com", password: ValidPassword, custom_fields: { name: "Partial User", email: "partial@test.com", user_type: "teacher" } });
            expect(res.status).toBe(201);
        });

        test("RadioInvalidValue", async () => {
            const body = await shouldFail(
                ampRegister({ email: "radio_invalid@test.com", password: ValidPassword, custom_fields: { name: "Radio User", email: "radio@test.com", user_type: "student", shift: "weekend" } }),
                400
            );
            assertErrID(body, ErrFieldValidationOnSchemaRegister);
            assertTrace(body, "invalid form value for shift: type(radio) value(weekend)");
        });

        test("RadioValidValue", async () => {
            const res = await ampRegister({ email: "radio_valid@test.com", password: ValidPassword, custom_fields: { name: "Radio User", email: "radio2@test.com", user_type: "teacher", shift: "night" } });
            expect(res.status).toBe(201);
        });

        test("EmailInvalidFormat", async () => {
            const body = await shouldFail(
                ampRegister({ email: "valid@example.com", password: ValidPassword, custom_fields: { name: "Invalid Email User", email: "not-an-email", user_type: "student" } }),
                400
            );
            assertErrID(body, ErrFieldValidationOnSchemaRegister);
            assertTrace(body, "invalid form value for email: type(email) value(not-an-email)");
        });

        test("EmailMissingAt", async () => {
            const body = await shouldFail(
                ampRegister({ email: "valid2@example.com", password: ValidPassword, custom_fields: { name: "No At User", email: "invalidemail.com", user_type: "teacher" } }),
                400
            );
            assertErrID(body, ErrFieldValidationOnSchemaRegister);
            assertTrace(body, "invalid form value for email: type(email) value(invalidemail.com)");
        });

        test("EmailValidSuccess", async () => {
            const res = await ampRegister({ email: "valid_email_test@example.com", password: ValidPassword, custom_fields: { name: "Valid Email User", email: "user@company.com", user_type: "student", student_id: "2023005001" } });
            expect(res.status).toBe(201);
        });

        test("ExtraFieldsIgnored", async () => {
            const res = await ampRegister({ email: "extra_ignored@test.com", password: ValidPassword, custom_fields: { name: "Extra Fields", email: "extra@test.com", user_type: "teacher", unknown_field: "should be ignored", another_extra: 123 } });
            expect(res.status).toBe(201);
        });

        test("NullOptionalFields", async () => {
            const res = await ampRegister({ email: "null_optional@test.com", password: ValidPassword, custom_fields: { name: "Null User", email: "null@test.com", user_type: "staff", is_active: null, activation_date: null, shift: null } });
            expect(res.status).toBe(201);
        });
    });
});