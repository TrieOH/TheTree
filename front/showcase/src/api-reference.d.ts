declare module "virtual:api-reference" {
  const reference: Record<
    string,
    {
      name: string;
      type: string;
      required: boolean;
      defaultValue?: string;
      description: string;
    }[]
  >;
  export default reference;
}
