import type { FieldDefinitionResultI } from "../../types/fields-types";

export const createMockField = (
  id: string, 
  key: string, 
  title: string, 
  type: FieldDefinitionResultI['type'] = 'string',
  overrides: Partial<FieldDefinitionResultI> = {}
): FieldDefinitionResultI => ({
  id,
  object_id: "obj_1",
  key,
  title,
  type,
  placeholder: `Digite seu ${title.toLowerCase()}`,
  description: "",
  position: 0,
  options: [],
  default_value: "",
  mutable: true,
  required: true,
  owner: "user",
  visibility_rules: [],
  required_rules: [],
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  ...overrides,
});

export const MOCK_FIELDS = {
  USER_TYPE: createMockField("f_utype", "user_type", "Tipo de Usuário", "select", {
    options: [
      { id: "opt1", label: "Pessoa Física", value: "personal", position: 1 },
      { id: "opt2", label: "Pessoa Jurídica", value: "business", position: 2 },
    ],
    default_value: "personal",
  }),

  COMPANY_NAME: (dependsOnId: string) => createMockField("f_company", "company_name", "Nome da Empresa", "string", {
    visibility_rules: [
      {
        id: "rule_vis_company",
        depends_on_field_id: dependsOnId,
        operator: "equals",
        value: "business",
      },
    ],
  }),

  COUNTRY: createMockField("f_country", "country", "País", "select", {
    options: [
      { id: "c1", label: "Brasil", value: "br", position: 1 },
      { id: "c2", label: "Estados Unidos", value: "us", position: 2 },
      { id: "c3", label: "Portugal", value: "pt", position: 3 },
    ],
    default_value: "us",
  }),

  CPF: (dependsOnId: string) => createMockField("f_cpf", "cpf", "CPF", "string", {
    required: false,
    placeholder: "000.000.000-00",
    required_rules: [
      {
        id: "rule_req_cpf",
        depends_on_field_id: dependsOnId,
        operator: "equals",
        value: "br",
      },
    ],
  }),

  AGE: createMockField("f_age", "age", "Idade", "int", { required: false }),

  NEWSLETTER: createMockField("f_news", "newsletter", "Desejo receber novidades", "bool", { required: false }),
  
  GENDER: createMockField("f_gender", "gender", "Gênero", "radio", {
    options: [
      { id: "g1", label: "Masculino", value: "m", position: 1 },
      { id: "g2", label: "Feminino", value: "f", position: 2 },
      { id: "g3", label: "Outro", value: "o", position: 3 },
    ],
  }),

  INTERESTS: createMockField("f_interests", "interests", "Interesses", "checkbox", {
    required: false,
    options: [
      { id: "i1", label: "Tecnologia", value: "tech", position: 1 },
      { id: "i2", label: "Esportes", value: "sports", position: 2 },
      { id: "i3", label: "Música", value: "music", position: 3 },
    ],
  }),

  BIO: (dependsOnId: string) => createMockField("f_bio", "bio", "Biografia Curta", "string", {
    description: "Aparece apenas se tiver mais de 18 anos",
    visibility_rules: [
      {
        id: "rule_vis_age",
        depends_on_field_id: dependsOnId,
        operator: "gte",
        value: "18",
      }
    ]
  }),

  ZIPCODE: (countryFieldId: string) => createMockField("f_zip", "zipcode", "CEP/Zip Code", "string", {
    required_rules: [
      {
        id: "rule_req_zip",
        depends_on_field_id: countryFieldId,
        operator: "in",
        value: "br,us",
      }
    ]
  })
};
