import type { Meta, StoryObj } from '@storybook/react-vite';
import SignUpWithProvider from './components/SignUpWithProvider';

const meta = {
  title: "Authentication/SignUp",
  component: SignUpWithProvider,
  parameters: { layout: 'centered' },
  tags: ['autodocs'],
  argTypes: {
    flow_id: { control: 'text' },
  },
  args: {
    flow_id: 'default',
  },
} satisfies Meta<typeof SignUpWithProvider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    flow_id: "test"
  }
};

export const WithLoginLink: Story = {
  args: {
    flow_id: 'test',
    loginRedirect: () => alert("Redirect to Login")
  },
};

export const WithCustomFields: Story = {
  args: {
    flow_id: 'onboarding',
    fields: [
      {
        id: "f1",
        key: "name",
        title: "Nome Completo",
        type: "string",
        placeholder: "Seu nome aqui",
        required: true,
        position: 1,
        object_id: "obj1",
        description: "",
        options: [],
        default_value: "",
        mutable: true,
        owner: "user",
        visibility_rules: [],
        required_rules: [],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      },
      {
        id: "f2",
        key: "age",
        title: "Idade",
        type: "int",
        placeholder: "Sua idade",
        required: false,
        position: 2,
        object_id: "obj1",
        description: "",
        options: [],
        default_value: "",
        mutable: true,
        owner: "user",
        visibility_rules: [],
        required_rules: [],
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      },
            {
              id: "f3",
              key: "gender",
              title: "Gênero",
              type: "radio",
              placeholder: "Selecione",
              required: true,
              position: 3,
              object_id: "obj1",
              description: "",
              options: [
                { id: "o1", label: "Masculino", value: "male", position: 1 },
                { id: "o2", label: "Feminino", value: "female", position: 2 },
                { id: "o3", label: "Outro", value: "other", position: 3 },
              ],
              default_value: "",
              mutable: true,
              owner: "user",
              visibility_rules: [],
              required_rules: [],
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
            {
              id: "f4",
              key: "interests",
              title: "Interesses",
              type: "checkbox",
              placeholder: "",
              required: false,
              position: 4,
              object_id: "obj1",
              description: "",
              options: [
                { id: "i1", label: "Tecnologia", value: "tech", position: 1 },
                { id: "i2", label: "Esportes", value: "sports", position: 2 },
                { id: "i3", label: "Música", value: "music", position: 3 },
              ],
              default_value: "",
              mutable: true,
              owner: "user",
              visibility_rules: [],
              required_rules: [],
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
            {
              id: "f5",
              key: "newsletter",
              title: "Desejo receber novidades",
              type: "bool",
              placeholder: "",
              required: false,
              position: 5,
              object_id: "obj1",
              description: "",
              options: [],
              default_value: "false",
              mutable: true,
              owner: "user",
              visibility_rules: [],
              required_rules: [],
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
            {
              id: "f6",
              key: "country",
              title: "País",
              type: "select",
              placeholder: "Selecione seu país",
              required: true,
              position: 6,
              object_id: "obj1",
              description: "",
              options: [
                { id: "c1", label: "Brasil", value: "br", position: 1 },
                { id: "c2", label: "Estados Unidos", value: "us", position: 2 },
                { id: "c3", label: "Portugal", value: "pt", position: 3 },
              ],
              default_value: "",
              mutable: true,
              owner: "user",
              visibility_rules: [],
              required_rules: [],
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            }
          ]
        }
      };
      