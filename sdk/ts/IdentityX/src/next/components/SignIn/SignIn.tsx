export interface SignInProps {
  /** What background color to use */
  backgroundColor?: string;
  /** Click Handler - Perform Login */
  onSubmit: () => void;
}

export function SignIn({
  backgroundColor="#dfcdcd"
}: SignInProps) {
  return (
    <div className="test" style={{ backgroundColor }}>
      <p>Faça seu Login</p>
    </div>
  );
}
