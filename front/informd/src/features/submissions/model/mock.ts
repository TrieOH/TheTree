import type { FieldI, FieldSelectConfigI } from "#/features/fields/model";
import type { FormI } from "#/features/forms/model";
import type { StepI } from "#/features/steps/model";
import type { SubmitRequestI } from ".";

const MOCK_FORM: FormI = {
  id: "form-001",
  owner_id: "user-001",
  created_by: "user-001",
  title: "Personal Identity Verification",
  status: "open",
  created_at: "2024-01-15T10:00:00Z",
  updated_at: "2024-01-15T10:00:00Z",
};

const MOCK_STEPS: StepI[] = [
  {
    id: "step-001",
    form_id: "form-001",
    title: "Personal Identity",
    description: "Collect full legal name, date of birth, and social security number for verification.",
    position_hint: 1,
  },
  {
    id: "step-002",
    form_id: "form-001",
    title: "Contact Information",
    description: "Provide your email address and phone number so we can reach you.",
    position_hint: 2,
  },
  {
    id: "step-003",
    form_id: "form-001",
    title: "Preferences",
    description: "Select your preferred contact methods and notification settings.",
    position_hint: 3,
  },
  {
    id: "step-004",
    form_id: "form-001",
    title: "Review & Submit",
    description: "Review your information before final submission.",
    position_hint: 4,
  },
];

const MOCK_FIELDS: Record<string, FieldI[]> = {
  "step-001": [
    {
      id: "field-001",
      step_id: "step-001",
      key: "full_legal_name",
      title: "Full Legal Name",
      description: "Include middle name if applicable.",
      position_hint: 1,
      required: true,
      type: "string",
      placeholder: { value: "Enter your full name as it appears on your ID" },
      created_at: "",
      updated_at: "",
    },
    {
      id: "field-002",
      step_id: "step-001",
      key: "date_of_birth",
      title: "Date of Birth",
      position_hint: 2,
      required: true,
      type: "date",
      placeholder: { value: "dd/mm/aaaa" },
      created_at: "",
      updated_at: "",
    },
    {
      id: "field-003",
      step_id: "step-001",
      key: "social_security",
      title: "Social Security Number",
      position_hint: 3,
      required: true,
      type: "string",
      placeholder: { value: "000-00-0000" },
      created_at: "",
      updated_at: "",
    },
  ],
  "step-002": [
    {
      id: "field-004",
      step_id: "step-002",
      key: "email",
      title: "Email Address",
      description: "We'll send confirmation to this address.",
      position_hint: 1,
      required: true,
      type: "email",
      placeholder: { value: "you@example.com" },
      created_at: "",
      updated_at: "",
    },
    {
      id: "field-005",
      step_id: "step-002",
      key: "phone",
      title: "Phone Number",
      description: "Include country code if outside US.",
      position_hint: 2,
      required: true,
      type: "phone",
      placeholder: { value: "+1 (000) 000-0000" },
      created_at: "",
      updated_at: "",
    },
    {
      id: "field-006",
      step_id: "step-002",
      key: "website",
      title: "Personal Website (Optional)",
      position_hint: 3,
      required: false,
      type: "url",
      placeholder: { value: "https://yourwebsite.com" },
      created_at: "",
      updated_at: "",
    },
  ],
  "step-003": [
    {
      id: "field-007",
      step_id: "step-003",
      key: "contact_method",
      title: "Preferred Contact Method",
      description: "How would you like us to reach you?",
      position_hint: 1,
      required: true,
      type: "select",
      config: {
        behaviour: "radio",
        value_type: "string",
        options: [
          { value: "email", label: "Email" },
          { value: "phone", label: "Phone Call" },
          { value: "sms", label: "SMS/Text" },
        ],
      },
      created_at: "",
      updated_at: "",
    },
    {
      id: "field-008",
      step_id: "step-003",
      key: "notifications",
      title: "Notification Preferences",
      description: "Select all that apply.",
      position_hint: 2,
      required: false,
      type: "select",
      config: {
        behaviour: "checkbox",
        value_type: "string",
        options: [
          { value: "marketing", label: "Marketing updates" },
          { value: "security", label: "Security alerts" },
          { value: "product", label: "Product news" },
        ],
      },
      created_at: "",
      updated_at: "",
    },
    {
      id: "field-009",
      step_id: "step-003",
      key: "age_verification",
      title: "I confirm I am 18 years or older",
      position_hint: 3,
      required: true,
      type: "bool",
      created_at: "",
      updated_at: "",
    },
  ],
  "step-004": [],
};

function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export async function getForm(): Promise<FormI> {
  await delay(400);
  return MOCK_FORM;
}

export async function getSteps(): Promise<StepI[]> {
  await delay(400);
  return MOCK_STEPS.sort((a, b) => a.position_hint - b.position_hint);
}

export async function getFields(stepId: string): Promise<FieldI[]> {
  await delay(400);
  return MOCK_FIELDS[stepId];
}

export async function getFieldSelectConfig(
  fieldId: string,
  _formId: string,
  _stepId: string,
  _namespaceId?: string
): Promise<FieldSelectConfigI> {
  await delay(200);
  const field = Object.values(MOCK_FIELDS)
    .flat()
    .find((f) => f.id === fieldId);
  return {
    field_id: fieldId,
    behaviour: field?.config?.behaviour ?? "dropdown-radio",
    value_type: field?.config?.value_type ?? "string",
    options: field?.config?.options ?? [],
  };
}

export async function submitForm(request: SubmitRequestI): Promise<{ success: boolean }> {
  await delay(800);
  console.log("Submit Request:", request);
  return { success: true };
}