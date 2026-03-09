import { beforeAll, describe, expect, test } from "vitest";
import { createClient } from "./helpers/index.js";
import {
    assertErrID,
    assertMessage,
    assertValidationError,
    shouldFail,
} from "./helpers/assert.js";
import {
    Validate,
    ValidateExact,
    AnyUUID,
    AnyDate,
    AsString,
    InOrder,
    ByKey,
    Store,
} from "./helpers/validate.js";
import { ValidPassword } from "./fixtures/auth/testdata.js";

// Error IDs
// RequestValidationError              = fail.ID(0, "REQ", 4, false, ...)   → 0_REQ_0004_D
// RequestMissingQueryParam            = fail.ID(0, "REQ", 1, false, ...)   → 0_REQ_0001_D
// SchemaNotOwnedByPrincipal           = fail.ID(0, "SCHEMA", 0, false, ...) → 0_SCHEMA_0000_D
// SCHEMANoPublishedVersion            = fail.ID(0, "SCHEMA", 0, true, ...)  → 0_SCHEMA_0000_S
// SchemaFlowIDAlreadyExistsInType     = fail.ID(0, "SCHEMA", 1, true, ...)  → 0_SCHEMA_0001_S
// SchemaInvalidSchemaType             = fail.ID(0, "SCHEMA", 2, true, ...)  → 0_SCHEMA_0002_S
// SchemaFlowIDIsReserved              = fail.ID(0, "SCHEMA", 3, false, ...) → 0_SCHEMA_0003_D
// SchemaInvalidFlowID                 = fail.ID(0, "SCHEMA", 2, false, ...) → 0_SCHEMA_0002_D
// SchemaHasOnlyDraftVersion           = fail.ID(0, "SCHEMA", 3, true, ...)  → 0_SCHEMA_0003_S
// SchemaTryingToPublishPublished      = fail.ID(0, "SCHEMA", 5, true, ...)  → 0_SCHEMA_0005_S
// SchemaVersionDraftDoesntExist       = fail.ID(0, "SCHEMAVERSION", 1, true, ...) → 0_SCHEMAVERSION_0001_S
// SchemaVersionPublishWithNoFields    = fail.ID(0, "SCHEMAVERSION", 0, true, ...) → 0_SCHEMAVERSION_0000_S
// SchemaVersionTryingToPublishPublished = fail.ID(0, "SCHEMAVERSION", 2, true, ...) → 0_SCHEMAVERSION_0002_S
// SchemaVersionDraftOnNonPublished    = fail.ID(0, "SCHEMAVERSION", 7, true, ...) → 0_SCHEMAVERSION_0007_S
// SchemaVersionNoChanges              = fail.ID(0, "SCHEMAVERSION", 8, true, ...) → 0_SCHEMAVERSION_0008_S
// FIELDSamePositionForMultipleFields  = fail.ID(0, "FIELD", 6, false, ...) → 0_FIELD_0006_D
// FIELDSameKeyForMultipleFields       = fail.ID(0, "FIELD", 5, false, ...) → 0_FIELD_0005_D
// FIELDInvalidCharactersInKey         = fail.ID(0, "FIELD", 7, false, ...) → 0_FIELD_0007_D
// FIELDKeyAlreadyExists               = fail.ID(0, "FIELD", 8, false, ...) → 0_FIELD_0008_D
// AuthNotClient                       = fail.ID(0, "AUTH", 2, true, ...)   → 0_AUTH_0002_S
// SQLNotFound                         = fail.ID(0, "SQL", 0, false, ...)   → 0_SQL_0000_D
const ErrRequestValidation            = "0_REQ_0004_D";
const ErrRequestMissingQueryParam     = "0_REQ_0001_D";
const ErrSchemaNotOwned               = "0_SCHEMA_0000_D";
const ErrSchemaNoPublishedVersion     = "0_SCHEMA_0000_S";
const ErrSchemaFlowIDExists           = "0_SCHEMA_0001_S";
const ErrSchemaInvalidType            = "0_SCHEMA_0002_S";
const ErrSchemaFlowIDReserved         = "0_SCHEMA_0003_D";
const ErrSchemaInvalidFlowID          = "0_SCHEMA_0002_D";
const ErrSchemaHasOnlyDraft           = "0_SCHEMA_0003_S";
const ErrSchemaTryingToPublishPub     = "0_SCHEMA_0005_S";
const ErrSchemaVersionNoDraft         = "0_SCHEMAVERSION_0001_S";
const ErrSchemaVersionNoFields        = "0_SCHEMAVERSION_0000_S";
const ErrSchemaVersionAlreadyPub      = "0_SCHEMAVERSION_0002_S";
const ErrSchemaVersionDraftOnNonPub   = "0_SCHEMAVERSION_0007_S";
const ErrSchemaVersionNoChanges       = "0_SCHEMAVERSION_0008_S";
const ErrFieldSamePosition            = "0_FIELD_0006_D";
const ErrFieldSameKey                 = "0_FIELD_0005_D";
const ErrFieldInvalidChars            = "0_FIELD_0007_D";
const ErrFieldKeyExists               = "0_FIELD_0008_D";
const ErrAuthNotClient                = "0_AUTH_0002_S";
const ErrSQLNotFound                  = "0_SQL_0000_D";

// ============================================================================
// SCHEMAS TESTS
// ============================================================================

describe("Schemas", () => {
    let user: any;
    let projectID: string;
    let schemaID: string;
    let schemaVersion1ID: string;
    let schemaVersion2ID: string;
    let schemaVersion3ID: string;

    beforeAll(async () => {
        await createClient().withCredentials("schemas@mail.com", ValidPassword).register();
        user = await createClient().withCredentials("schemas@mail.com", ValidPassword).login();
        const project = await user.post("/projects", {
            project_name: "schema testing",
            metadata: { env: "test" },
        });
        projectID = project.id;
    });

    test("PublishSchemaRandomID", async () => {
        const fakeID = "00000000-0000-7000-8000-000000000001";
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${fakeID}/publish`),
            401
        );
        assertErrID(body, ErrSchemaNotOwned);
        assertMessage(body, "cannot publish a schema you don't own");
    });

    test("Draft", async () => {
        const data = await user.post(`/projects/${projectID}/schemas`, {
            schema_type: "context",
            title: "scti-register-flow",
            flow_id: "scti-register",
        });

        const ref = { current: "" };
        Validate(data, {
            id:                 Store(ref, AnyUUID),
            project_id:         AsString(projectID, AnyUUID),
            title:              "scti-register-flow",
            flow_id:            "scti-register",
            type:               "context",
            status:             "draft",
            current_version_id: null,
            created_at:         AnyDate,
            updated_at:         AnyDate,
        });
        schemaID = ref.current as string;
    });

    test("DraftAnother", async () => {
        const data = await user.post(`/projects/${projectID}/schemas`, {
            schema_type: "context",
            title: "eenge",
            flow_id: "estudante",
        });
        Validate(data, {
            id:                 AnyUUID,
            project_id:         AsString(projectID, AnyUUID),
            title:              "eenge",
            flow_id:            "estudante",
            type:               "context",
            status:             "draft",
            current_version_id: null,
            created_at:         AnyDate,
            updated_at:         AnyDate,
        });
    });

    test("DraftSameFlowIDAndType", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas`, {
                schema_type: "context",
                title: "eenge",
                flow_id: "estudante",
            }),
            409
        );
        assertErrID(body, ErrSchemaFlowIDExists);
        assertMessage(body, "schema with this flow ID already exists in this type");
    });

    test("DraftReservedFlowID", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas`, {
                schema_type: "context",
                title: "Reserved",
                flow_id: "none",
            }),
            400
        );
        assertErrID(body, ErrSchemaFlowIDReserved);
        assertMessage(body, "flow id can't be the reserved keyword 'none'");
    });

    test("DraftFlowIDSameAsType", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas`, {
                schema_type: "context",
                title: "SameAsType",
                flow_id: "context",
            }),
            400
        );
        assertErrID(body, ErrSchemaInvalidFlowID);
        assertMessage(body, "flow id can't be the same as a schema type");
    });

    describe("DraftValidation", () => {
        test("InvalidType", async () => {
            const body = await shouldFail(
                user.post(`/projects/${projectID}/schemas`, {
                    schema_type: "invalid-type",
                    title: "test",
                    flow_id: "test",
                }),
                400
            );
            assertErrID(body, ErrRequestValidation);
            assertValidationError(body, "schema_type must be one of: core, context, sub-context");
        });

        test("FlowIDTooLong", async () => {
            const body = await shouldFail(
                user.post(`/projects/${projectID}/schemas`, {
                    schema_type: "context",
                    title: "test",
                    flow_id: "this-flow-id-is-way-too-long-and-should-fail-validation-because-it-exceeds-63-chars",
                }),
                400
            );
            assertErrID(body, ErrRequestValidation);
            assertValidationError(body, "flow_id must be at most 63 characters long");
        });
    });

    test("PublishSchemaNoVersion", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/publish`),
            400
        );
        assertErrID(body, ErrSchemaNoPublishedVersion);
        assertMessage(body, "cannot publish a schema with no versions");
    });

    test("PublishVersionNoDraft", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/versions/publish`),
            401
        );
        assertErrID(body, ErrSchemaVersionNoDraft);
        assertMessage(body, "cannot publish a schema with a version draft that doesn't exist");
    });

    test("DraftVersion", async () => {
        const data = await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/draft`);
        const ref = { current: "" };
        Validate(data, {
            id:             Store(ref, AnyUUID),
            schema_id:      AsString(schemaID, AnyUUID),
            version_number: 1,
        });
        schemaVersion1ID = ref.current as string;
    });

    test("CheckSchemaVersion", async () => {
        const data = await user.get(`/projects/${projectID}/schemas/${schemaID}`);
        Validate(data, {
            id:                 AsString(schemaID, AnyUUID),
            project_id:         AsString(projectID, AnyUUID),
            title:              "scti-register-flow",
            flow_id:            "scti-register",
            type:               "context",
            status:             "draft",
            current_version_id: AsString(schemaVersion1ID, AnyUUID),
        });
    });

    test("PublishSchemaDraftVersion", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/publish`),
            400
        );
        assertErrID(body, ErrSchemaHasOnlyDraft);
        assertMessage(body, "cannot publish a schema with only draft versions");
    });

    test("DraftVersionError", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/versions/draft`),
            400
        );
        assertErrID(body, ErrSchemaVersionDraftOnNonPub);
        assertMessage(body, "new versions can only be drafted from published versions");
    });

    test("PublishVersionFieldsError", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/versions/publish`),
            400
        );
        assertErrID(body, ErrSchemaVersionNoFields);
        assertMessage(body, "cannot publish a schema version with no fields");
    });

    test("CreateFieldsSamePosition", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/v1`, {
                fields: [
                    { key: "matricula", type: "string", owner: "user", title: "Numero da Matrícula", description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045", required: true, mutable: true, position: 0 },
                    { key: "curso",     type: "string", owner: "user", title: "Curso de Matrícula",  description: "O curso que você está matrículado na UENF",              placeholder: "Ciência da Computação",                    required: true, mutable: true, position: 0 },
                ],
            }),
            409
        );
        assertErrID(body, ErrFieldSamePosition);
        assertMessage(body, "two fields can't occupy the same position");
    });

    test("CreateFieldsSameKey", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/v1`, {
                fields: [
                    { key: "matricula", type: "string", owner: "user", title: "Numero da Matrícula", description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045", required: true, mutable: true, position: 0 },
                    { key: "matricula", type: "string", owner: "user", title: "Numero da Matrícula", description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045", required: true, mutable: true, position: 1 },
                ],
            }),
            409
        );
        assertErrID(body, ErrFieldSameKey);
        assertMessage(body, "two fields can't have the same key");
    });

    test("CreateFields", async () => {
        const data = await user.post(`/projects/${projectID}/schemas/${schemaID}/v1`, {
            fields: [
                { key: "matricula", type: "string", owner: "user", title: "Numero da Matrícula", description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045", required: true, mutable: true, position: 0 },
                { key: "curso",     type: "string", owner: "user", title: "Curso de Matrícula",  description: "O curso que você está matrículado na UENF",              placeholder: "Ciência da Computação",                    required: true, mutable: true, position: 1 },
            ],
        });
        Validate(data, [
            { object_id: AnyUUID, id: AnyUUID },
            { object_id: AnyUUID, id: AnyUUID },
        ]);
    });

    test("PublishVersionSuccess", async () => {
        await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/publish`);
    });

    test("PublishVersionAlreadyPublished", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/versions/publish`),
            401
        );
        assertErrID(body, ErrSchemaVersionAlreadyPub);
        assertMessage(body, "cannot publish a schema version that is already published");
    });

    test("PublishSchemaSuccess", async () => {
        await user.post(`/projects/${projectID}/schemas/${schemaID}/publish`);
    });

    test("PublishSchemaAlreadyPublished", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/publish`),
            401
        );
        assertErrID(body, ErrSchemaTryingToPublishPub);
        assertMessage(body, "cannot publish a schema that is already published");
    });

    test("DraftVersion2", async () => {
        const data = await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/draft`);
        const ref = { current: "" };
        Validate(data, {
            id:             Store(ref, AnyUUID),
            schema_id:      AsString(schemaID, AnyUUID),
            version_number: 2,
        });
        schemaVersion2ID = ref.current as string;
    });

    test("CheckSchemaVersionAfterV2Draft", async () => {
        const data = await user.get(`/projects/${projectID}/schemas/${schemaID}`);
        Validate(data, {
            id:                 AsString(schemaID, AnyUUID),
            project_id:         AsString(projectID, AnyUUID),
            title:              "scti-register-flow",
            flow_id:            "scti-register",
            type:               "context",
            status:             "published",
            current_version_id: AsString(schemaVersion1ID, AnyUUID),
        });
    });

    test("PublishVersion2NoChanges", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/versions/publish`),
            400
        );
        assertErrID(body, ErrSchemaVersionNoChanges);
        assertMessage(body, "cannot publish a version with no changes");
    });

    test("AddFieldToV2FailKeyCheck", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/v2`, {
                fields: [
                    { key: "período", type: "int", owner: "user", title: "Período Atual", description: "O período da sua matéria mais avançada da grade", required: true, mutable: true, position: 2 },
                ],
            }),
            400
        );
        assertErrID(body, ErrFieldInvalidChars);
        assertMessage(body, "field key must start with a lowercase letter and contain only lowercase letters, numbers, or underscores");
    });

    test("AddFieldToV2Success", async () => {
        const data = await user.post(`/projects/${projectID}/schemas/${schemaID}/v2`, {
            fields: [
                { key: "periodo", type: "int", owner: "user", title: "Período Atual", description: "O período da sua matéria mais avançada da grade", required: true, mutable: true, position: 2 },
            ],
        });
        Validate(data, [
            { object_id: AnyUUID, id: AnyUUID, key: "periodo", type: "int", owner: "user", title: "Período Atual", description: "O período da sua matéria mais avançada da grade", required: true, mutable: true, position: 2 },
        ]);
    });

    test("CreateFieldDuplicateInherited", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/v2`, {
                fields: [
                    { key: "matricula", type: "string", owner: "user", title: "Numero da Matrícula", description: "Duplicate", required: true, mutable: true, position: 3 },
                ],
            }),
            409
        );
        assertErrID(body, ErrFieldSameKey);
        assertMessage(body, "two fields can't have the same key");
    });

    test("CreateFieldDuplicateInDraft", async () => {
        const body = await shouldFail(
            user.post(`/projects/${projectID}/schemas/${schemaID}/v2`, {
                fields: [
                    { key: "periodo", type: "int", owner: "user", title: "Duplicate", description: "Duplicate", required: true, mutable: true, position: 4 },
                ],
            }),
            409
        );
        assertErrID(body, ErrFieldSameKey);
        assertMessage(body, "two fields can't have the same key");
    });

    test("PublishVersion2Success", async () => {
        await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/publish`);
    });

    test("GetSchemaVerbose", async () => {
        const schema = await user.get(`/projects/${projectID}/schemas/${schemaID}/verbose`);

        const matriculaV1Ref = { current: "" };
        const matriculaV2Ref = { current: "" };
        const cursoV1Ref     = { current: "" };
        const cursoV2Ref     = { current: "" };

        Validate(schema, {
            id:                 AsString(schemaID, AnyUUID),
            project_id:         AsString(projectID, AnyUUID),
            title:              "scti-register-flow",
            flow_id:            "scti-register",
            type:               "context",
            status:             "published",
            current_version_id: AsString(schemaVersion2ID, AnyUUID),
            created_at:         AnyDate,
            updated_at:         AnyDate,
            versions: InOrder([
                {
                    id:             AsString(schemaVersion2ID, AnyUUID),
                    schema_id:      AsString(schemaID, AnyUUID),
                    version_number: 2,
                    fields: ByKey("key", {
                        matricula: { object_id: AnyUUID, id: Store(matriculaV2Ref, AnyUUID), key: "matricula", type: "string", owner: "user", title: "Numero da Matrícula", description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045", required: true, mutable: true, position: 0 },
                        curso:     { object_id: AnyUUID, id: Store(cursoV2Ref, AnyUUID),     key: "curso",     type: "string", owner: "user", title: "Curso de Matrícula",  description: "O curso que você está matrículado na UENF",              placeholder: "Ciência da Computação",                    required: true, mutable: true, position: 1 },
                        periodo:   { object_id: AnyUUID, id: AnyUUID,                        key: "periodo",   type: "int",    owner: "user", title: "Período Atual",       description: "O período da sua matéria mais avançada da grade",                                                                  required: true, mutable: true, position: 2 },
                    }),
                },
                {
                    id:             AsString(schemaVersion1ID, AnyUUID),
                    schema_id:      AsString(schemaID, AnyUUID),
                    version_number: 1,
                    fields: ByKey("key", {
                        matricula: { object_id: AnyUUID, id: Store(matriculaV1Ref, AnyUUID), key: "matricula", type: "string", owner: "user", title: "Numero da Matrícula", description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045", required: true, mutable: true, position: 0 },
                        curso:     { object_id: AnyUUID, id: Store(cursoV1Ref, AnyUUID),     key: "curso",     type: "string", owner: "user", title: "Curso de Matrícula",  description: "O curso que você está matrículado na UENF",              placeholder: "Ciência da Computação",                    required: true, mutable: true, position: 1 },
                    }),
                },
            ]),
        });

        // Cross-version field ID stability
        expect(matriculaV1Ref.current).toBe(matriculaV2Ref.current);
        expect(cursoV1Ref.current).toBe(cursoV2Ref.current);
    });

    test("DraftVersion3WithOptionsAndRules", async () => {
        const data = await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/draft`);
        const ref = { current: "" };
        ValidateExact(data, {
            id:                  Store(ref, AnyUUID),
            schema_id:           AsString(schemaID, AnyUUID),
            version_number:      3,
            status:              "draft",
            based_on_version_id: AsString(schemaVersion2ID, AnyUUID),
            created_at:          AnyDate,
            updated_at:          AnyDate,
        });
        schemaVersion3ID = ref.current as string;
    });

    test("AddFieldsWithOptionsAndRules", async () => {
        const data = await user.post(`/projects/${projectID}/schemas/${schemaID}/v3`, {
            fields: [
                {
                    key: "user_type", type: "select", owner: "user", title: "Tipo de Usuário", required: true, mutable: true, position: 3,
                    options: [
                        { value: "student",   label: "Estudante", position: 0 },
                        { value: "professor", label: "Professor", position: 1 },
                        { value: "visitor",   label: "Visitante", position: 2 },
                    ],
                },
                { key: "needs_scholarship", type: "bool", owner: "user", title: "Necessita de Bolsa?", required: true, mutable: true, position: 4 },
                {
                    key: "income", type: "int", owner: "user", title: "Renda Familiar", description: "Renda mensal familiar em reais", required: false, mutable: true, position: 5,
                    visibility_rules: [{ depends_on_field_key: "needs_scholarship", operator: "equals", value: true }],
                    required_rules:   [{ depends_on_field_key: "needs_scholarship", operator: "equals", value: true }],
                },
                {
                    key: "scholarship_type", type: "radio", owner: "user", title: "Tipo de Bolsa", required: false, mutable: true, position: 6,
                    options: [
                        { value: "full",    label: "Integral", position: 0 },
                        { value: "partial", label: "Parcial",  position: 1 },
                    ],
                    visibility_rules: [
                        { depends_on_field_key: "user_type",          operator: "equals", value: "student" },
                        { depends_on_field_key: "needs_scholarship",  operator: "equals", value: true },
                    ],
                },
            ],
        });

        Validate(data, [
            { object_id: AnyUUID, id: AnyUUID, schema_id: AsString(schemaID, AnyUUID), schema_version_id: AsString(schemaVersion3ID, AnyUUID), key: "user_type",          type: "select", owner: "user", title: "Tipo de Usuário",     description: null, placeholder: null, required: true,  mutable: true, default_value: null, position: 3, created_at: AnyDate, updated_at: AnyDate },
            { object_id: AnyUUID, id: AnyUUID, schema_id: AsString(schemaID, AnyUUID), schema_version_id: AsString(schemaVersion3ID, AnyUUID), key: "needs_scholarship",  type: "bool",   owner: "user", title: "Necessita de Bolsa?", description: null, placeholder: null, required: true,  mutable: true, default_value: null, position: 4, created_at: AnyDate, updated_at: AnyDate },
            { object_id: AnyUUID, id: AnyUUID, schema_id: AsString(schemaID, AnyUUID), schema_version_id: AsString(schemaVersion3ID, AnyUUID), key: "income",             type: "int",    owner: "user", title: "Renda Familiar",      description: "Renda mensal familiar em reais", placeholder: null, required: false, mutable: true, default_value: null, position: 5, created_at: AnyDate, updated_at: AnyDate },
            { object_id: AnyUUID, id: AnyUUID, schema_id: AsString(schemaID, AnyUUID), schema_version_id: AsString(schemaVersion3ID, AnyUUID), key: "scholarship_type",   type: "radio",  owner: "user", title: "Tipo de Bolsa",       description: null, placeholder: null, required: false, mutable: true, default_value: null, position: 6, created_at: AnyDate, updated_at: AnyDate },
        ]);
    });

    test("PublishVersion3Success", async () => {
        await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/publish`);
    });

    test("GetLatestFormByID", async () => {
        const data = await user.get(`/projects/${projectID}/schemas/${schemaID}/latest`);
        ValidateExact(data, {
            id:             AsString(schemaVersion3ID, AnyUUID),
            schema_id:      AsString(schemaID, AnyUUID),
            title:          "scti-register-flow",
            flow_id:        "scti-register",
            schema_type:    "context",
            version_id:     AsString(schemaVersion3ID, AnyUUID),
            version_number: 3,
            status:         "published",
            created_at:     AnyDate,
            updated_at:     AnyDate,
            fields: [
                { id: AnyUUID, object_id: AnyUUID, key: "matricula",          type: "string", owner: "user", title: "Numero da Matrícula",     description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045",          required: true,  mutable: true, default_value: null, position: 0, created_at: AnyDate, updated_at: AnyDate, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "curso",              type: "string", owner: "user", title: "Curso de Matrícula",      description: "O curso que você está matrículado na UENF",              placeholder: "Ciência da Computação", required: true,  mutable: true, default_value: null, position: 1, created_at: AnyDate, updated_at: AnyDate, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "periodo",            type: "int",    owner: "user", title: "Período Atual",           description: "O período da sua matéria mais avançada da grade",        placeholder: null,                    required: true,  mutable: true, default_value: null, position: 2, created_at: AnyDate, updated_at: AnyDate, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "user_type",          type: "select", owner: "user", title: "Tipo de Usuário",         description: null, placeholder: null, required: true,  mutable: true, default_value: null, position: 3, created_at: AnyDate, updated_at: AnyDate,
                    options: [
                        { id: AnyUUID, value: "student",   label: "Estudante", position: 0 },
                        { id: AnyUUID, value: "professor", label: "Professor", position: 1 },
                        { id: AnyUUID, value: "visitor",   label: "Visitante", position: 2 },
                    ],
                    visibility_rules: [], required_rules: [],
                },
                { id: AnyUUID, object_id: AnyUUID, key: "needs_scholarship",  type: "bool",   owner: "user", title: "Necessita de Bolsa?",     description: null, placeholder: null, required: true,  mutable: true, default_value: null, position: 4, created_at: AnyDate, updated_at: AnyDate, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "income",             type: "int",    owner: "user", title: "Renda Familiar",           description: "Renda mensal familiar em reais", placeholder: null, required: false, mutable: true, default_value: null, position: 5, created_at: AnyDate, updated_at: AnyDate, options: [],
                    visibility_rules: [{ id: AnyUUID, depends_on_field_id: AnyUUID, operator: "equals", value: true }],
                    required_rules:   [{ id: AnyUUID, depends_on_field_id: AnyUUID, operator: "equals", value: true }],
                },
                { id: AnyUUID, object_id: AnyUUID, key: "scholarship_type",   type: "radio",  owner: "user", title: "Tipo de Bolsa",           description: null, placeholder: null, required: false, mutable: true, default_value: null, position: 6, created_at: AnyDate, updated_at: AnyDate,
                    options: [
                        { id: AnyUUID, value: "full",    label: "Integral", position: 0 },
                        { id: AnyUUID, value: "partial", label: "Parcial",  position: 1 },
                    ],
                    visibility_rules: [
                        { id: AnyUUID, depends_on_field_id: AnyUUID, operator: "equals", value: "student" },
                        { id: AnyUUID, depends_on_field_id: AnyUUID, operator: "equals", value: true },
                    ],
                    required_rules: [],
                },
            ],
        });
    });

    test("GetSpecificFormV2", async () => {
        const data = await user.get(`/projects/${projectID}/schemas/${schemaID}/v2`);
        Validate(data, {
            id:             AsString(schemaVersion2ID, AnyUUID),
            schema_id:      AsString(schemaID, AnyUUID),
            title:          "scti-register-flow",
            flow_id:        "scti-register",
            schema_type:    "context",
            version_id:     AsString(schemaVersion2ID, AnyUUID),
            version_number: 2,
            status:         "published",
            created_at:     AnyDate,
            updated_at:     AnyDate,
            fields: [
                { id: AnyUUID, object_id: AnyUUID, key: "matricula", type: "string", owner: "user", title: "Numero da Matrícula", description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045",          required: true, mutable: true, position: 0, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "curso",     type: "string", owner: "user", title: "Curso de Matrícula",  description: "O curso que você está matrículado na UENF",              placeholder: "Ciência da Computação", required: true, mutable: true, position: 1, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "periodo",   type: "int",    owner: "user", title: "Período Atual",       description: "O período da sua matéria mais avançada da grade",        placeholder: null,                    required: true, mutable: true, position: 2, options: [], visibility_rules: [], required_rules: [] },
            ],
        });
    });

    test("GetLatestFormByFlowLookup", async () => {
        const data = await user.http.get(
            `/projects/${projectID}/schemas/lookup/latest?flow_id=scti-register&schema_type=context`,
            { headers: { Cookie: `access_token=${user.auth!.accessToken}; refresh_token=${user.auth!.refreshToken}` } }
        );
        Validate(data.data.data, {
            id:             AsString(schemaVersion3ID, AnyUUID),
            schema_id:      AsString(schemaID, AnyUUID),
            title:          "scti-register-flow",
            flow_id:        "scti-register",
            schema_type:    "context",
            version_number: 3,
            status:         "published",
            created_at:     AnyDate,
            updated_at:     AnyDate,
            fields: [
                { id: AnyUUID, object_id: AnyUUID, key: "matricula",         type: "string", owner: "user", title: "Numero da Matrícula",     description: "Sua matrícula da UENF como aparece no sistema acadêmico", placeholder: "20223200045",          required: true,  mutable: true, default_value: null, position: 0, created_at: AnyDate, updated_at: AnyDate, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "curso",             type: "string", owner: "user", title: "Curso de Matrícula",      description: "O curso que você está matrículado na UENF",              placeholder: "Ciência da Computação", required: true,  mutable: true, default_value: null, position: 1, created_at: AnyDate, updated_at: AnyDate, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "periodo",           type: "int",    owner: "user", title: "Período Atual",           description: "O período da sua matéria mais avançada da grade",        placeholder: null,                    required: true,  mutable: true, default_value: null, position: 2, created_at: AnyDate, updated_at: AnyDate, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "user_type",         type: "select", owner: "user", title: "Tipo de Usuário",         description: null, placeholder: null, required: true,  mutable: true, default_value: null, position: 3, created_at: AnyDate, updated_at: AnyDate,
                    options: [
                        { id: AnyUUID, value: "student",   label: "Estudante", position: 0 },
                        { id: AnyUUID, value: "professor", label: "Professor", position: 1 },
                        { id: AnyUUID, value: "visitor",   label: "Visitante", position: 2 },
                    ],
                    visibility_rules: [], required_rules: [],
                },
                { id: AnyUUID, object_id: AnyUUID, key: "needs_scholarship", type: "bool",   owner: "user", title: "Necessita de Bolsa?",     description: null, placeholder: null, required: true,  mutable: true, default_value: null, position: 4, created_at: AnyDate, updated_at: AnyDate, options: [], visibility_rules: [], required_rules: [] },
                { id: AnyUUID, object_id: AnyUUID, key: "income",            type: "int",    owner: "user", title: "Renda Familiar",           description: "Renda mensal familiar em reais", placeholder: null, required: false, mutable: true, default_value: null, position: 5, created_at: AnyDate, updated_at: AnyDate, options: [],
                    visibility_rules: [{ id: AnyUUID, depends_on_field_id: AnyUUID, operator: "equals", value: true }],
                    required_rules:   [{ id: AnyUUID, depends_on_field_id: AnyUUID, operator: "equals", value: true }],
                },
                { id: AnyUUID, object_id: AnyUUID, key: "scholarship_type",  type: "radio",  owner: "user", title: "Tipo de Bolsa",           description: null, placeholder: null, required: false, mutable: true, default_value: null, position: 6, created_at: AnyDate, updated_at: AnyDate,
                    options: [
                        { id: AnyUUID, value: "full",    label: "Integral", position: 0 },
                        { id: AnyUUID, value: "partial", label: "Parcial",  position: 1 },
                    ],
                    visibility_rules: [
                        { id: AnyUUID, depends_on_field_id: AnyUUID, operator: "equals", value: "student" },
                        { id: AnyUUID, depends_on_field_id: AnyUUID, operator: "equals", value: true },
                    ],
                    required_rules: [],
                },
            ],
        });
    });

    test("GetFormByFlowLookupInvalidType", async () => {
        const body = await shouldFail(
            user.http.get(
                `/projects/${projectID}/schemas/lookup/latest?flow_id=scti-register&schema_type=core`,
                { headers: { Cookie: `access_token=${user.auth!.accessToken}; refresh_token=${user.auth!.refreshToken}` } }
            ),
            404
        );
        assertErrID(body, ErrSQLNotFound);
        assertMessage(body, "schema not found");
    });

    test("GetFormByFlowLookupMissingRequiredQuery", async () => {
        const body = await shouldFail(
            user.http.get(
                `/projects/${projectID}/schemas/lookup/latest?schema_type=context`,
                { headers: { Cookie: `access_token=${user.auth!.accessToken}; refresh_token=${user.auth!.refreshToken}` } }
            ),
            400
        );
        assertErrID(body, ErrRequestMissingQueryParam);
        assertMessage(body, "missing query parameter: flow_id");
    });

    test("ProjectUserAccessDenied", async () => {
        await createClient().http.post(`/projects/${projectID}/register`, {
            email: "proj-user-schema@mail.com",
            password: ValidPassword,
        });
        const projUser = await createClient()
            .withCredentials("proj-user-schema@mail.com", ValidPassword)
            .projectLogin(projectID);

        const expectForbidden = async (promise: Promise<any>) => {
            const body = await shouldFail(promise, 403);
            assertErrID(body, ErrAuthNotClient);
            assertMessage(body, "only clients can access this endpoint");
        };

        await expectForbidden(projUser.post(`/projects/${projectID}/schemas`, { schema_type: "context", title: "forbidden", flow_id: "forbidden" }));
        await expectForbidden(projUser.post(`/projects/${projectID}/schemas/${schemaID}/publish`));
        await expectForbidden(projUser.get(`/projects/${projectID}/schemas/${schemaID}`));
        await expectForbidden(projUser.get(`/projects/${projectID}/schemas/${schemaID}/verbose`));
        await expectForbidden(projUser.post(`/projects/${projectID}/schemas/${schemaID}/versions/draft`));
    });
});

// ============================================================================
// BATCH UPDATE FIELDS
// ============================================================================

describe("BatchUpdateFields", () => {
    let user: any;
    let projectID: string;
    let schemaID: string;
    let fieldAID: string;
    let fieldBID: string;
    let fieldCID: string;

    beforeAll(async () => {
        await createClient().withCredentials("batch-fields@mail.com", ValidPassword).register();
        user = await createClient().withCredentials("batch-fields@mail.com", ValidPassword).login();
        const project = await user.post("/projects", { project_name: "batch field testing", metadata: { env: "test" } });
        projectID = project.id;

        // SetupSchema
        const schema = await user.post(`/projects/${projectID}/schemas`, { schema_type: "context", title: "batch-test-flow", flow_id: "batch-test" });
        schemaID = schema.id;

        // SetupVersion
        await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/draft`);

        // SetupInitialFields
        await user.post(`/projects/${projectID}/schemas/${schemaID}/v1`, {
            fields: [
                { key: "field_a", type: "string", owner: "user", title: "Field A", required: true,  mutable: true, position: 0 },
                { key: "field_b", type: "string", owner: "user", title: "Field B", required: true,  mutable: true, position: 1 },
                { key: "field_c", type: "string", owner: "user", title: "Field C", required: false, mutable: true, position: 2 },
            ],
        });

        // GetInitialFields — capture object_ids
        const form = await user.get(`/projects/${projectID}/schemas/${schemaID}/v1`);
        const fields = form.fields;
        fieldAID = fields[0].object_id;
        fieldBID = fields[1].object_id;
        fieldCID = fields[2].object_id;
    });

    test("BatchReorderFields", async () => {
        const data = await user.put(`/projects/${projectID}/schemas/${schemaID}/v1/fields`, {
            fields: [
                { object_id: fieldCID, position: 0 },
                { object_id: fieldBID, position: 1 },
                { object_id: fieldAID, position: 2 },
            ],
        });
        Validate(data, [
            { object_id: AsString(fieldCID, AnyUUID), key: "field_c", title: "Field C", position: 0 },
            { object_id: AsString(fieldBID, AnyUUID), key: "field_b", title: "Field B", position: 1 },
            { object_id: AsString(fieldAID, AnyUUID), key: "field_a", title: "Field A", position: 2 },
        ]);
    });

    test("BatchReorderWithUpdates", async () => {
        const updatedDesc = "Updated description";
        const data = await user.put(`/projects/${projectID}/schemas/${schemaID}/v1/fields`, {
            fields: [
                { object_id: fieldAID, title: "Field A Updated", position: 0, required: false, mutable: false, description: updatedDesc },
                { object_id: fieldBID, title: "Field B Updated", position: 1 },
                { object_id: fieldCID, title: "Field C Updated", position: 2, required: true },
            ],
        });
        Validate(data, [
            { object_id: AsString(fieldAID, AnyUUID), key: "field_a", title: "Field A Updated", position: 0, required: false, mutable: false, description: updatedDesc },
            { object_id: AsString(fieldBID, AnyUUID), key: "field_b", title: "Field B Updated", position: 1 },
            { object_id: AsString(fieldCID, AnyUUID), key: "field_c", title: "Field C Updated", position: 2, required: true },
        ]);
    });

    test("BatchUpdateAndAddNewFields", async () => {
        const data = await user.put(`/projects/${projectID}/schemas/${schemaID}/v1/fields`, {
            fields: [
                { object_id: fieldAID, title: "Field A Final", position: 0 },
                { object_id: fieldBID, position: 1 },
                { object_id: fieldCID, position: 2 },
                { key: "new_alpha", type: "string", owner: "user", title: "New Alpha", required: true,  position: 3 },
                { key: "new_beta",  type: "int",    owner: "user", title: "New Beta",  required: false, position: 4 },
            ],
        });
        Validate(data, [
            { object_id: AsString(fieldAID, AnyUUID), key: "field_a",  title: "Field A Final", position: 0 },
            { object_id: AsString(fieldBID, AnyUUID), key: "field_b",  position: 1 },
            { object_id: AsString(fieldCID, AnyUUID), key: "field_c",  position: 2 },
            { object_id: AnyUUID,                     key: "new_alpha", type: "string", owner: "user", title: "New Alpha", required: true,  position: 3 },
            { object_id: AnyUUID,                     key: "new_beta",  type: "int",    owner: "user", title: "New Beta",  required: false, position: 4 },
        ]);
    });

    test("BatchUpdateWithInvalidPosition", async () => {
        const body = await shouldFail(
            user.put(`/projects/${projectID}/schemas/${schemaID}/v1/fields`, {
                fields: [
                    { object_id: fieldAID, position: 0 },
                    { object_id: fieldBID, position: 0 },
                ],
            }),
            409
        );
        assertErrID(body, ErrFieldSamePosition);
        assertMessage(body, "two fields can't occupy the same position");
    });

    test("BatchUpdateWithInvalidKey", async () => {
        const body = await shouldFail(
            user.put(`/projects/${projectID}/schemas/${schemaID}/v1/fields`, {
                fields: [
                    { object_id: fieldAID, key: "field_b", position: 0 },
                    { object_id: fieldBID, position: 1 },
                ],
            }),
            409
        );
        assertErrID(body, ErrFieldKeyExists);
        assertMessage(body, "field key 'field_b' already exists in this version");
    });

    test("BatchUpdateDuplicatePositionInNewFields", async () => {
        const body = await shouldFail(
            user.put(`/projects/${projectID}/schemas/${schemaID}/v1/fields`, {
                fields: [
                    { object_id: fieldAID, position: 0 },
                    { key: "new_field_1", type: "string", owner: "user", title: "New Field 1", position: 1 },
                    { key: "new_field_2", type: "string", owner: "user", title: "New Field 2", position: 1 },
                ],
            }),
            409
        );
        assertErrID(body, ErrFieldSamePosition);
        assertMessage(body, "two fields can't occupy the same position");
    });

    test("BatchUpdateDuplicateKeyInNewFields", async () => {
        const body = await shouldFail(
            user.put(`/projects/${projectID}/schemas/${schemaID}/v1/fields`, {
                fields: [
                    { object_id: fieldAID, position: 0 },
                    { key: "duplicate_key", type: "string", owner: "user", title: "Duplicate 1", position: 1 },
                    { key: "duplicate_key", type: "string", owner: "user", title: "Duplicate 2", position: 2 },
                ],
            }),
            409
        );
        assertErrID(body, ErrFieldSameKey);
        assertMessage(body, "two fields can't have the same key");
    });

    test("BatchUpdateWithNewFieldDuplicateKey", async () => {
        const body = await shouldFail(
            user.put(`/projects/${projectID}/schemas/${schemaID}/v1/fields`, {
                fields: [
                    { object_id: fieldAID, position: 0 },
                    { key: "field_b", type: "string", owner: "user", title: "Duplicate of existing", position: 1 },
                ],
            }),
            409
        );
        assertErrID(body, ErrFieldKeyExists);
        assertMessage(body, "field key 'field_b' already exists in this version");
    });

    test("BatchUpdateOnlyExistingFields", async () => {
        const desc = "Only updating field A";
        const data = await user.put(`/projects/${projectID}/schemas/${schemaID}/v1/fields`, {
            fields: [
                { object_id: fieldAID, title: "Field A Only", description: desc, position: 0 },
                { object_id: fieldBID, title: "Field B Only", position: 1 },
                { object_id: fieldCID, title: "Field C Only", position: 2 },
                { key: "field_x", type: "string", position: 5 },
                { key: "field_y", type: "string", position: 6 },
            ],
        });
        Validate(data, [
            { object_id: AsString(fieldAID, AnyUUID), key: "field_a",  title: "Field A Only", description: desc, position: 0 },
            { object_id: AsString(fieldBID, AnyUUID), key: "field_b",  title: "Field B Only", position: 1 },
            { object_id: AsString(fieldCID, AnyUUID), key: "field_c",  title: "Field C Only", position: 2 },
            { key: "new_alpha", type: "string", position: 3 },
            { key: "new_beta",  type: "int",    position: 4 },
            { key: "field_x",   type: "string", position: 5 },
            { key: "field_y",   type: "string", position: 6 },
        ]);
    });
});

// ============================================================================
// DELETE FIELD OPTIONS AND RULES
// ============================================================================

describe("DeleteFieldOptionsAndRules", () => {
    let user: any;
    let projectID: string;
    let schemaID: string;
    let fieldWithOptionsID: string;
    let optionID1: string;
    let fieldWithRulesID: string;
    let visibilityRuleID: string;
    let requiredRuleID: string;

    beforeAll(async () => {
        await createClient().withCredentials("delete-options-rules@mail.com", ValidPassword).register();
        user = await createClient().withCredentials("delete-options-rules@mail.com", ValidPassword).login();
        const project = await user.post("/projects", { project_name: "delete options rules test", metadata: { env: "test" } });
        projectID = project.id;

        // SetupSchema
        const schema = await user.post(`/projects/${projectID}/schemas`, { schema_type: "context", title: "delete-test-flow", flow_id: "delete-test" });
        schemaID = schema.id;

        // SetupVersion
        await user.post(`/projects/${projectID}/schemas/${schemaID}/versions/draft`);

        // CreateFieldWithOptions
        await user.post(`/projects/${projectID}/schemas/${schemaID}/v1`, {
            fields: [
                {
                    key: "user_type", type: "select", owner: "user", title: "User Type", required: true, mutable: true, position: 0,
                    options: [
                        { value: "admin", label: "Admin", position: 0 },
                        { value: "user",  label: "User",  position: 1 },
                    ],
                },
            ],
        });

        const form1 = await user.get(`/projects/${projectID}/schemas/${schemaID}/v1`);
        fieldWithOptionsID = form1.fields[0].object_id;
        optionID1          = form1.fields[0].options[0].id;

        // CreateFieldWithRules
        await user.post(`/projects/${projectID}/schemas/${schemaID}/v1`, {
            fields: [
                { key: "has_car",   type: "bool",   owner: "user", title: "Has Car",   required: false, mutable: true, position: 1 },
                {
                    key: "car_model", type: "string", owner: "user", title: "Car Model", required: false, mutable: true, position: 2,
                    visibility_rules: [{ depends_on_field_key: "has_car", operator: "equals", value: true }],
                    required_rules:   [{ depends_on_field_key: "has_car", operator: "equals", value: true }],
                },
            ],
        });

        const form2 = await user.get(`/projects/${projectID}/schemas/${schemaID}/v1`);
        fieldWithRulesID = form2.fields[2].object_id;
        visibilityRuleID = form2.fields[2].visibility_rules[0].id;
        requiredRuleID   = form2.fields[2].required_rules[0].id;
    });

    test("DeleteOption", async () => {
        await user.del(`/projects/${projectID}/schemas/${schemaID}/v1/fields/${fieldWithOptionsID}/options/${optionID1}`);
    });

    test("VerifyOptionDeleted", async () => {
        const form = await user.get(`/projects/${projectID}/schemas/${schemaID}/v1`);
        expect(form.fields[0].options).toHaveLength(1);
    });

    test("DeleteVisibilityRule", async () => {
        await user.del(`/projects/${projectID}/schemas/${schemaID}/v1/fields/${fieldWithRulesID}/visibility-rules/${visibilityRuleID}`);
    });

    test("VerifyVisibilityRuleDeleted", async () => {
        const form = await user.get(`/projects/${projectID}/schemas/${schemaID}/v1`);
        expect(form.fields[2].visibility_rules).toHaveLength(0);
    });

    test("DeleteRequiredRule", async () => {
        await user.del(`/projects/${projectID}/schemas/${schemaID}/v1/fields/${fieldWithRulesID}/required-rules/${requiredRuleID}`);
    });

    test("VerifyRequiredRuleDeleted", async () => {
        const form = await user.get(`/projects/${projectID}/schemas/${schemaID}/v1`);
        expect(form.fields[2].required_rules).toHaveLength(0);
    });
});